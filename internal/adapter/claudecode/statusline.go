package claudecode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
