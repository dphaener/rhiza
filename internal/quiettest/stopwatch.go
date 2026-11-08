package quiettest

import (
	"fmt"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Stopwatch times one run of a budgeted operation by the clock, and also
// by the time the operation spent runnable but waiting for a CPU, so a
// budget can judge the operation's own time rather than the machine's
// load. Waiting for a CPU is what another process on the machine (an
// agent's build, another worktree's suite, the runner's neighbours) adds to
// a read; it says nothing about the read. Disk, fsync and sleeps still
// count, as they would for a person waiting on the read.
//
// Start locks the calling goroutine to its OS thread until Stop, and only
// that thread's wait comes off: summed over every thread, a parallel test's
// or a runtime thread's wait would come off a read that never waited. Work
// the operation hands to other goroutines is timed by the clock alone, so
// own time can only be over the operation's, never under. Call Stop on the
// goroutine that called Start.
//
// It also reads the machine's CPU pressure over the run (see Timing).
type Stopwatch struct {
	start    time.Time
	wait     time.Duration
	pressure time.Duration
}

// Start starts a Stopwatch.
func Start() Stopwatch {
	runtime.LockOSThread()
	return Stopwatch{wait: threadWait(), pressure: cpuPressure(), start: time.Now()}
}

// Lap is one timed run: its wall-clock time, its own time, and Pressure,
// how long some task on the machine waited for a CPU during it.
type Lap struct {
	Wall, Own, Pressure time.Duration
}

// Stop reads the stopwatch and unlocks the goroutine from its thread: the
// wall-clock time since Start, and Own, that time less what the thread
// spent runnable but not running. Where the kernel reports no such wait
// (macOS, or a kernel without schedstats), Own is Wall.
func (s Stopwatch) Stop() Lap {
	wall := time.Since(s.start)
	wait := threadWait() - s.wait
	pressure := min(max(cpuPressure()-s.pressure, 0), wall)
	runtime.UnlockOSThread()
	if wait < 0 || wait > wall {
		wait = 0
	}
	return Lap{Wall: wall, Own: wall - wait, Pressure: pressure}
}

// cpuPressure is how long, in all, some task on the machine has waited for
// a CPU: the total of the "some" line of /proc/pressure/cpu (Linux PSI, in
// microseconds). It is 0 where that file is missing or unreadable.
func cpuPressure() time.Duration {
	b, err := os.ReadFile("/proc/pressure/cpu")
	if err != nil {
		return 0
	}
	for line := range strings.Lines(string(b)) {
		f := strings.Fields(line)
		if len(f) == 0 || f[0] != "some" {
			continue
		}
		for _, kv := range f[1:] {
			if v, ok := strings.CutPrefix(kv, "total="); ok {
				if n, err := strconv.ParseInt(v, 10, 64); err == nil {
					return time.Duration(n) * time.Microsecond
				}
			}
		}
	}
	return 0
}

// threadWait is how long the calling thread has spent runnable but not
// running: the second field of /proc/thread-self/schedstat (Linux, in
// nanoseconds). It is 0 where that file is missing or unreadable.
func threadWait() time.Duration {
	b, err := os.ReadFile("/proc/thread-self/schedstat")
	if err != nil {
		return 0
	}
	f := strings.Fields(string(b))
	if len(f) < 2 {
		return 0
	}
	n, err := strconv.ParseInt(f[1], 10, 64)
	if err != nil {
		return 0
	}
	return time.Duration(n)
}

// Timing sums up the laps of one budgeted operation: the wall-clock
// minimum, median and maximum, Own, the median of the laps' own times,
// which is what a budget judges, and Pressure, the share of the laps' time
// in which some task on the machine waited for a CPU.
type Timing struct {
	Min, Median, Max, Own time.Duration
	Pressure              float64
}

// SaturatedPressure is the CPU pressure at or over which a machine is too
// busy to judge a budget. Own time takes off the waits for a CPU, but not
// a CPU slowed by its neighbours: on the four-core VM, two busy loops per
// core of the test still doubled the reads' own time. A quiet machine,
// such as CI's budgets job, runs at a few percent.
const SaturatedPressure = 0.2

// Saturated says whether t was measured on a machine too busy to judge a
// budget by: an over-budget Timing that is Saturated is inconclusive, not
// a failure. Never under GitHub Actions: CI's budgets job runs alone and
// is the judge, so a budget there always passes or fails. Never where the
// kernel reports no pressure (macOS).
func (t Timing) Saturated() bool {
	return t.Pressure >= SaturatedPressure && !onCI
}

// onCI is whether the tests run under GitHub Actions.
var onCI = os.Getenv("GITHUB_ACTIONS") == "true"

// OnCI reports whether the tests run under GitHub Actions, where no
// Timing is Saturated.
func OnCI() bool { return onCI }

// Over reports whether t's own time is at or over budget on a machine
// quiet enough to judge it; when it is over only on a Saturated machine, it
// reports false and why, for the test's log.
func (t Timing) Over(budget time.Duration) (over bool, why string) {
	if t.Own < budget {
		return false, ""
	}
	if t.Saturated() {
		return false, fmt.Sprintf("inconclusive: own time %v is over %v, but the machine's CPU pressure was %.0f%%, too busy to judge (CI's budgets job judges it)", t.Own, budget, 100*t.Pressure)
	}
	return true, ""
}

// Summarize returns laps' Timing; the zero Timing for no laps.
func Summarize(laps []Lap) Timing {
	if len(laps) == 0 {
		return Timing{}
	}
	wall := make([]time.Duration, len(laps))
	own := make([]time.Duration, len(laps))
	var total, pressure time.Duration
	for i, l := range laps {
		wall[i], own[i] = l.Wall, l.Own
		total += l.Wall
		pressure += l.Pressure
	}
	slices.Sort(wall)
	slices.Sort(own)
	t := Timing{Min: wall[0], Median: median(wall), Max: wall[len(wall)-1], Own: median(own)}
	if total > 0 {
		t.Pressure = float64(pressure) / float64(total)
	}
	return t
}

// median is the middle of sorted ds, or the mean of the middle two.
func median(ds []time.Duration) time.Duration {
	n := len(ds)
	if n%2 == 1 {
		return ds[n/2]
	}
	return (ds[n/2-1] + ds[n/2]) / 2
}

// String reads "median 3ms, own 2ms (min 1ms, max 9ms), cpu pressure 4%".
func (t Timing) String() string {
	return fmt.Sprintf("median %v, own %v (min %v, max %v), cpu pressure %.0f%%", t.Median, t.Own, t.Min, t.Max, 100*t.Pressure)
}
