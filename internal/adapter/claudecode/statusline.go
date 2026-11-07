package claudecode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
)

// LiveUsage implements adapter.StatusLineReader. Claude Code's status line
// input carries rate_limits with a five_hour and a seven_day window, each
// with used_percentage and resets_at (epoch seconds). The five-hour window
// is the session kind and the seven-day one the weekly kind; a window
// missing either figure is skipped.
func (Adapter) LiveUsage(payload []byte) ([]adapter.LiveWindow, error) {
	var in struct {
		RateLimits map[string]struct {
			UsedPercentage *float64 `json:"used_percentage"`
			ResetsAt       *float64 `json:"resets_at"`
		} `json:"rate_limits"`
	}
	if err := json.Unmarshal(payload, &in); err != nil {
		return nil, fmt.Errorf("status line payload: %w", err)
	}
	var out []adapter.LiveWindow
	for _, w := range []struct{ name, kind string }{{"five_hour", "session"}, {"seven_day", "weekly"}} {
		r, ok := in.RateLimits[w.name]
		if !ok || r.UsedPercentage == nil || r.ResetsAt == nil || *r.ResetsAt <= 0 {
			continue
		}
		out = append(out, adapter.LiveWindow{
			Kind:           w.kind,
			UsedPercentage: *r.UsedPercentage,
			ResetsAt:       time.Unix(int64(*r.ResetsAt), 0).UTC(),
		})
	}
	return out, nil
}

// UserStatusLine implements adapter.StatusLineReader: the command of the
// statusLine in <config>/settings.json, "" when there is none.
func (Adapter) UserStatusLine(getenv func(string) string) string {
	dir := configDir(getenv)
	if dir == "" {
		return ""
	}
	b, err := os.ReadFile(filepath.Join(dir, "settings.json"))
	if err != nil {
		return ""
	}
	var s struct {
		StatusLine struct {
			Type    string `json:"type"`
			Command string `json:"command"`
		} `json:"statusLine"`
	}
	if json.Unmarshal(b, &s) != nil || (s.StatusLine.Type != "" && s.StatusLine.Type != "command") {
		return ""
	}
	return s.StatusLine.Command
}

// liveFile is the file in the usage state folder that the status line
// recording keeps (daemon.LiveFile): a JSON array of snapshots, one per kind.
const liveFile = "usage.live.json"

// LiveSnapshots implements adapter.LiveUsageReader by reading the live-usage
// file in dir. Entries without a kind or a reset time are skipped, and a
// later entry of a kind replaces an earlier one.
func (Adapter) LiveSnapshots(dir string) []adapter.LiveSnapshot {
	b, err := os.ReadFile(filepath.Join(dir, liveFile))
	if err != nil {
		return nil
	}
	var in []struct {
		Kind           string    `json:"kind"`
		UsedPercentage float64   `json:"used_percentage"`
		ResetsAt       time.Time `json:"resets_at"`
		CapturedAt     time.Time `json:"captured_at"`
	}
	if json.Unmarshal(b, &in) != nil {
		return nil
	}
	byKind := map[string]adapter.LiveSnapshot{}
	for _, s := range in {
		if s.Kind == "" || s.ResetsAt.IsZero() {
			continue
		}
		byKind[s.Kind] = adapter.LiveSnapshot{Kind: s.Kind, UsedPercentage: s.UsedPercentage, ResetsAt: s.ResetsAt, CapturedAt: s.CapturedAt}
	}
	kinds := make([]string, 0, len(byKind))
	for k := range byKind {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	out := make([]adapter.LiveSnapshot, 0, len(kinds))
	for _, k := range kinds {
		out = append(out, byKind[k])
	}
	return out
}
