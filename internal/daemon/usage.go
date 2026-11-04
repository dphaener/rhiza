package daemon

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
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

// RecordLimitHit appends a hit to the limits file.
func RecordLimitHit(common string, h LimitHit) error {
	if err := os.MkdirAll(Dir(common), 0o755); err != nil {
		return err
	}
	b, err := json.Marshal(h)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(Dir(common), LimitsFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
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

// ReadLimitHits returns every recorded hit, oldest first as written.
// Lines that do not parse are skipped.
func ReadLimitHits(common string) []LimitHit {
	f, err := os.Open(filepath.Join(Dir(common), LimitsFile))
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
