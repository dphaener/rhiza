package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/session"
)

// limitAdapter is a fake agent whose transcript is $HOME/limit/<session>.txt
// and which is usage-limited while that file says "limited".
type limitAdapter struct {
	interviewAdapter
	launches *[]adapter.Launch
}

const limitAgent = "limit-fake"

func (limitAdapter) Name() string { return limitAgent }
func (a limitAdapter) Command(l adapter.Launch) ([]string, error) {
	*a.launches = append(*a.launches, l)
	return []string{limitAgent}, nil
}
func (limitAdapter) TranscriptGlobs(sid string, getenv func(string) string) []string {
	return []string{filepath.Join(getenv("HOME"), "limit", sid+".txt")}
}
func (limitAdapter) UsageLimit(t []byte) (adapter.UsageLimit, bool) {
	return adapter.UsageLimit{}, strings.Contains(string(t), "limited")
}

// removeClock removes a file when the run sleeps.
type removeClock struct{ path string }

func (c removeClock) Now() time.Time      { return time.Unix(0, 0) }
func (c removeClock) Sleep(time.Duration) { os.Remove(c.path) }

func TestSessionRunFreshSessionWhenTranscriptGone(t *testing.T) {
	f, s, _ := interviewFixture(t)
	var launches []adapter.Launch
	adapter.Register(limitAdapter{launches: &launches})
	if _, err := session.NewStore(session.Dir(f.common)).Update(s.ID, func(r *session.Record) {
		r.Agent, r.Prompt = limitAgent, "do the step"
	}); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(f.ctx.Home, "limit", s.ID+".txt")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("limited"), 0o644); err != nil {
		t.Fatal(err)
	}
	useClock(t, removeClock{path})

	resp, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(SessionRunResponse)
	if len(launches) != 2 || launches[1].Resume || launches[1].Prompt != "do the step" || launches[1].SessionID != s.ID {
		t.Errorf("launches %+v: want a second, fresh launch of the same session with its first prompt", launches)
	}
	if !strings.Contains(strings.Join(r.Warnings, "\n"), "started a fresh session") {
		t.Errorf("warnings %v do not record the fresh session", r.Warnings)
	}
}
