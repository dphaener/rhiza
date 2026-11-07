package daemon

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// LimitsFile is the file, under Dir, that run drivers append usage-limit
// hits to, one JSON line each; the daemon reads it at each poll.
const LimitsFile = "usage.limits.jsonl"

// LimitHit is one usage-limit hit a run driver saw.
type LimitHit struct {
	// HitAt is when the driver saw the session stop; ResetAt is the reset
	// time the agent gave, zero when it gave none.
	HitAt   time.Time `json:"hit_at"`
	ResetAt time.Time `json:"reset_at,omitempty"`
	// Agent names the adapter that saw it.
	Agent string `json:"agent,omitempty"`
}

// LearnedLimit is the limit the daemon learned for one window kind: the
// window's use at the last hit.
type LearnedLimit struct {
	Kind   string    `json:"kind"`
	Tokens int64     `json:"tokens"`
	HitAt  time.Time `json:"hit_at"`
}

// RootStateDir is the root's state folder, <root>/.rhiza-root: what every
// project's drivers and the one root daemon share.
func RootStateDir(root string) string { return filepath.Join(root, ".rhiza-root") }

// RecordLimitHit appends a hit to the limits file under a repo's git common
// dir.
func RecordLimitHit(common string, h LimitHit) error { return RecordLimitHitIn(Dir(common), h) }

// ReadLimitHits returns every hit recorded under a repo's git common dir.
func ReadLimitHits(common string) []LimitHit { return ReadLimitHitsIn(Dir(common)) }

// RecordLimitHitIn appends a hit to the limits file in dir. The usage
// governor of the root daemon keeps it in RootStateDir, so the hits of every
// project's drivers land in one file.
func RecordLimitHitIn(dir string, h LimitHit) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(dir, LimitsFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	// One write, so drivers appending at once do not interleave lines.
	_, err = f.Write(append(b, '\n'))
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	return err
}

// ReadLimitHitsIn returns every hit recorded in dir, oldest first as
// written. Lines that do not parse are skipped.
func ReadLimitHitsIn(dir string) []LimitHit {
	f, err := os.Open(filepath.Join(dir, LimitsFile))
	if err != nil {
		return nil
	}
	defer f.Close()
	var out []LimitHit
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var h LimitHit
		if json.Unmarshal(sc.Bytes(), &h) == nil && !h.HitAt.IsZero() {
			out = append(out, h)
		}
	}
	return out
}

// ActiveHold is the machine-wide hold the governor keeps from the hits
// recorded in dir: the hit with the latest reset time that has been seen
// and not yet reset at now. It reports false when no hit holds, which is
// the case from a hit's reset time on, for every project at once.
func ActiveHold(dir string, now time.Time) (LimitHit, bool) {
	var best LimitHit
	found := false
	for _, h := range ReadLimitHitsIn(dir) {
		if h.ResetAt.IsZero() || h.HitAt.After(now) || !now.Before(h.ResetAt) {
			continue
		}
		if !found || h.ResetAt.After(best.ResetAt) {
			best, found = h, true
		}
	}
	return best, found
}

// LiveFile is the file, under a usage state folder, that status lines
// record the account's live rate-limit windows to: one JSON object holding
// the latest snapshot per window kind.
const LiveFile = "usage.live.json"

// LiveSnapshot is one window's live figure as an agent's status line
// reported it.
type LiveSnapshot struct {
	Kind           string    `json:"kind"`
	UsedPercentage float64   `json:"used_percentage"`
	ResetsAt       time.Time `json:"resets_at"`
	CapturedAt     time.Time `json:"captured_at"`
	Agent          string    `json:"agent,omitempty"`
}

// RecordLive saves snaps in dir's live file, each replacing the earlier
// snapshot of its kind and leaving the other kinds as they were.
func RecordLive(dir string, snaps []LiveSnapshot) error {
	if len(snaps) == 0 {
		return nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	byKind := map[string]LiveSnapshot{}
	for _, s := range ReadLive(dir) {
		byKind[s.Kind] = s
	}
	for _, s := range snaps {
		byKind[s.Kind] = s
	}
	kinds := make([]string, 0, len(byKind))
	for k := range byKind {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	all := make([]LiveSnapshot, 0, len(kinds))
	for _, k := range kinds {
		all = append(all, byKind[k])
	}
	b, err := json.Marshal(all)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, LiveFile+".*")
	if err != nil {
		return err
	}
	_, err = tmp.Write(append(b, '\n'))
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(tmp.Name(), filepath.Join(dir, LiveFile))
	}
	if err != nil {
		os.Remove(tmp.Name())
	}
	return err
}

// ReadLive returns the latest snapshot of each window kind recorded in dir,
// ordered by kind; none when the file is missing or does not parse.
func ReadLive(dir string) []LiveSnapshot {
	b, err := os.ReadFile(filepath.Join(dir, LiveFile))
	if err != nil {
		return nil
	}
	var out []LiveSnapshot
	if json.Unmarshal(b, &out) != nil {
		return nil
	}
	return out
}
