package api

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/adapter/claudecode"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// trustedIn reads the config file at path and returns its projects' trust
// flags and the whole top level.
func trustedIn(t *testing.T, path string) (map[string]bool, map[string]json.RawMessage) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		t.Fatal(err)
	}
	var projects map[string]struct {
		Trusted bool `json:"hasTrustDialogAccepted"`
	}
	if err := json.Unmarshal(top["projects"], &projects); err != nil {
		t.Fatal(err)
	}
	out := map[string]bool{}
	for k, v := range projects {
		out[k] = v.Trusted
	}
	return out, top
}

func TestSessionRunTrustsFolders(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	home := t.TempDir()
	f.ctx.Home = home
	f.ctx.Getenv = func(string) string { return "" }
	cfg := filepath.Join(home, ".claude.json")
	existing := `{"numStartups": 7, "theme": "dark", "projects": {"/other": {"allowedTools": ["x"], "hasTrustDialogAccepted": false}, "` +
		f.worktree + `": {"lastCost": 1.5}}}`
	if err := os.WriteFile(cfg, []byte(existing), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	trusted, top := trustedIn(t, cfg)
	root := filepath.Dir(f.common)
	if !trusted[root] || !trusted[f.worktree] || trusted["/other"] {
		t.Errorf("trusted %v: want %s and %s only", trusted, root, f.worktree)
	}
	if string(top["numStartups"]) != "7" || string(top["theme"]) != `"dark"` {
		t.Errorf("other keys changed: %v", top)
	}
	var projects map[string]map[string]json.RawMessage
	if err := json.Unmarshal(top["projects"], &projects); err != nil {
		t.Fatal(err)
	}
	var tools bytes.Buffer
	if err := json.Compact(&tools, projects["/other"]["allowedTools"]); err != nil {
		t.Fatal(err)
	}
	if tools.String() != `["x"]` || string(projects[f.worktree]["lastCost"]) != "1.5" {
		t.Errorf("project keys changed: %v", projects)
	}
	if fi, err := os.Stat(cfg); err != nil || fi.Mode().Perm() != 0o600 {
		t.Errorf("config mode %v, %v", fi, err)
	}
}

func TestSessionRunTrustHonoursConfigDirAndCreatesFile(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	home, conf := t.TempDir(), filepath.Join(t.TempDir(), "conf")
	if err := os.MkdirAll(conf, 0o755); err != nil {
		t.Fatal(err)
	}
	f.ctx.Home = home
	f.ctx.Getenv = func(k string) string {
		if k == "CLAUDE_CONFIG_DIR" {
			return conf
		}
		return ""
	}
	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	trusted, _ := trustedIn(t, filepath.Join(conf, ".claude.json"))
	if !trusted[filepath.Dir(f.common)] || !trusted[f.worktree] {
		t.Errorf("trusted %v", trusted)
	}
	if _, err := os.Stat(filepath.Join(home, ".claude.json")); err == nil {
		t.Error("HOME's config was written although CLAUDE_CONFIG_DIR is set")
	}
}

func TestSessionRunTrustDryRunWritesNothing(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	home := t.TempDir()
	f.ctx.Home = home
	f.ctx.Getenv = func(string) string { return "" }
	f.ctx.DryRun = true
	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	if ents, _ := os.ReadDir(home); len(ents) != 0 {
		t.Errorf("a dry run wrote %v", ents)
	}
}

func TestSessionRunTrustWarnsOnUnreadableConfig(t *testing.T) {
	f, s, _ := runFixture(t, `{"permission_mode":"manual"}`)
	home := t.TempDir()
	f.ctx.Home = home
	f.ctx.Getenv = func(string) string { return "" }
	if err := os.WriteFile(filepath.Join(home, ".claude.json"), []byte("{nope"), 0o600); err != nil {
		t.Fatal(err)
	}
	resp, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	if w := strings.Join(resp.(SessionRunResponse).Warnings, "\n"); !strings.Contains(w, "trust the session's folders") {
		t.Errorf("warnings %q", w)
	}
	if data, _ := os.ReadFile(filepath.Join(home, ".claude.json")); string(data) != "{nope" {
		t.Errorf("a broken config was rewritten: %q", data)
	}
}

// trustLimitAdapter is the limit fake that trusts folders as Claude Code does.
type trustLimitAdapter struct{ limitAdapter }

const trustLimitName = "trust-limit-fake"

var trustRegister sync.Once

func (trustLimitAdapter) Name() string { return trustLimitName }
func (trustLimitAdapter) TrustFile(getenv func(string) string) string {
	return claudecode.Adapter{}.TrustFile(getenv)
}
func (trustLimitAdapter) Trust(cur []byte, dirs []string) ([]byte, error) {
	return claudecode.Adapter{}.Trust(cur, dirs)
}

func TestSessionRunResumeTrustsAgain(t *testing.T) {
	f, s, path := limitFixture(t)
	f.ctx.Getenv = func(string) string { return "" }
	// Swap in the trusting variant under the same name.
	trustRegister.Do(func() { adapter.Register(trustLimitAdapter{}) })
	if err := os.WriteFile(path, []byte("limited"), 0o644); err != nil {
		t.Fatal(err)
	}
	st := session.NewStore(session.Dir(f.common))
	if _, err := st.Update(s.ID, func(r *session.Record) { r.Agent, r.Run, r.Step = trustLimitName, "run-1", "build" }); err != nil {
		t.Fatal(err)
	}
	if err := run.NewStore(run.Dir(f.common)).Create("run-1", run.Entry{Type: run.TypeStarted, At: time.Now()}); err != nil {
		t.Fatal(err)
	}
	cfg := filepath.Join(f.ctx.Home, ".claude.json")
	// The first sleep, at the limit, drops the trust as an owner's cleanup
	// would; the resumed start must put it back.
	useClock(t, trustResetClock{path: path, cfg: cfg})
	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	trusted, _ := trustedIn(t, cfg)
	if !trusted[f.worktree] || !trusted[filepath.Dir(f.common)] {
		t.Errorf("after the resumed start, trusted %v", trusted)
	}
	if n := len(limitLaunches); n != 2 || !limitLaunches[1].Resume {
		t.Errorf("launches %+v: want a resumed second one", limitLaunches)
	}
}

// trustResetClock lifts the limit and removes the trust file when the run sleeps.
type trustResetClock struct{ path, cfg string }

func (c trustResetClock) Now() time.Time { return time.Unix(0, 0) }
func (c trustResetClock) Sleep(time.Duration) {
	os.WriteFile(c.path, []byte("fine"), 0o644)
	os.Remove(c.cfg)
}
