package daemon

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// DefaultMaxLoad is the load per core above which the daemon starts nothing
// new; DefaultMinMemory is the available memory under which it does not.
const (
	DefaultMaxLoad          = 3.0
	DefaultMinMemory uint64 = 2 << 30
)

// The reasons of a machine hold, also the skip reasons.
const (
	SkipLoadHigh  = "load-high"
	SkipMemoryLow = "memory-low"
)

// Reading is what the machine looked like at a poll. Load1 and Load5 are
// the 1- and 5-minute load averages per core, known when LoadOK; MemAvail
// is the available memory in bytes, known when MemOK. A figure that cannot
// be read is never held on.
type Reading struct {
	Load1, Load5 float64
	LoadOK       bool
	MemAvail     uint64
	MemOK        bool
}

// MachineHold is a hold of new runs for load or memory found by a poll.
type MachineHold struct {
	// Since is the first of the polls in a row that held for Reason.
	Since time.Time `json:"since"`
	// Reason is load-high or memory-low.
	Reason string `json:"reason"`
	// Detail names the reading and the limit.
	Detail string `json:"detail"`
}

// ParseLoadavg reads the 1- and 5-minute averages from /proc/loadavg text.
func ParseLoadavg(s string) (l1, l5 float64, err error) {
	f := strings.Fields(s)
	if len(f) < 2 {
		return 0, 0, fmt.Errorf("loadavg %q has no 1- and 5-minute figures", strings.TrimSpace(s))
	}
	if l1, err = strconv.ParseFloat(f[0], 64); err != nil {
		return 0, 0, fmt.Errorf("loadavg: %v", err)
	}
	if l5, err = strconv.ParseFloat(f[1], 64); err != nil {
		return 0, 0, fmt.Errorf("loadavg: %v", err)
	}
	return l1, l5, nil
}

// ParseMemAvailable reads MemAvailable, in bytes, from /proc/meminfo text.
func ParseMemAvailable(s string) (uint64, error) {
	for _, line := range strings.Split(s, "\n") {
		rest, ok := strings.CutPrefix(line, "MemAvailable:")
		if !ok {
			continue
		}
		f := strings.Fields(rest)
		if len(f) == 0 {
			break
		}
		n, err := strconv.ParseUint(f[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("meminfo: %v", err)
		}
		if len(f) > 1 && strings.EqualFold(f[1], "kB") {
			n <<= 10
		}
		return n, nil
	}
	return 0, fmt.Errorf("meminfo has no MemAvailable")
}

// ReadMachine reads the machine's load and available memory from /proc. A
// platform without them (macOS, others) gives a Reading with neither known.
func ReadMachine() Reading {
	var r Reading
	if b, err := os.ReadFile("/proc/loadavg"); err == nil {
		if l1, l5, err := ParseLoadavg(string(b)); err == nil {
			cores := float64(max(runtime.NumCPU(), 1))
			r.Load1, r.Load5, r.LoadOK = l1/cores, l5/cores, true
		}
	}
	if b, err := os.ReadFile("/proc/meminfo"); err == nil {
		if m, err := ParseMemAvailable(string(b)); err == nil {
			r.MemAvail, r.MemOK = m, true
		}
	}
	return r
}

// Evaluate decides the hold. Memory under minMem holds (0 turns that off).
// A load hold begins when the 1-minute load per core is over maxLoad and,
// once begun (prev is a load-high hold), ends only when the 5-minute load
// per core is back under it (0 turns the check off). It returns nil for no
// hold.
func Evaluate(r Reading, maxLoad float64, minMem uint64, prev *MachineHold, now time.Time) *MachineHold {
	since := func(reason string) time.Time {
		if prev != nil && prev.Reason == reason {
			return prev.Since
		}
		return now.UTC()
	}
	if minMem > 0 && r.MemOK && r.MemAvail < minMem {
		return &MachineHold{Since: since(SkipMemoryLow), Reason: SkipMemoryLow,
			Detail: fmt.Sprintf("%s available, under the minimum of %s", FormatSize(r.MemAvail), FormatSize(minMem))}
	}
	if maxLoad > 0 && r.LoadOK {
		if r.Load1 > maxLoad {
			return &MachineHold{Since: since(SkipLoadHigh), Reason: SkipLoadHigh,
				Detail: fmt.Sprintf("1-minute load %.2f per core, over the limit of %.2f", r.Load1, maxLoad)}
		}
		if prev != nil && prev.Reason == SkipLoadHigh && r.Load5 >= maxLoad {
			return &MachineHold{Since: prev.Since, Reason: SkipLoadHigh,
				Detail: fmt.Sprintf("5-minute load %.2f per core, not yet under the limit of %.2f", r.Load5, maxLoad)}
		}
	}
	return nil
}
