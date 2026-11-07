package web_test

import (
	"errors"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

// rootCtx is a context in dir that reaches git only, through Exec, so the
// root's per-project contexts each run git in their own checkout. Anything
// else it is asked to run fails.
func rootCtx(t *testing.T, dir string) *api.Ctx {
	t.Helper()
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	ctx := &api.Ctx{
		Root: dir,
		Wd:   dir,
		Now:  clock,
		IDs:  id.NewGenerator(clock, rand.New(rand.NewSource(3))),
		Exec: func(name string, args []string, d string, env []string) ([]byte, error) {
			if name != "git" {
				return nil, errors.New("not run in tests: " + name)
			}
			c := exec.Command("git", args...)
			c.Dir = d
			return c.Output()
		},
		IndexPath: filepath.Join(t.TempDir(), "index.db"),
		GitConfig: func(k string) (string, error) {
			if k == "user.email" {
				return "test@example.com", nil
			}
			return "", errors.New("unset")
		},
	}
	ctx.Git = func(args ...string) (string, error) {
		out, err := ctx.Exec("git", args, dir, nil)
		return string(out), err
	}
	return ctx
}

// The all-projects home reads a real root: two initialised projects show
// their live runs, inbox count and health, a bare one (no .rhiza/) is
// marked not initialised, and each links to its overview. Every page links
// to the home.
func TestProjectsHomeReadsTheRoot(t *testing.T) {
	cache := settledDir(t)
	t.Setenv("XDG_CACHE_HOME", cache)
	t.Setenv("HOME", cache)
	root := settledDir(t)
	ctxs := map[string]*api.Ctx{}
	for _, name := range []string{"alpha", "beta", "bare"} {
		dir := filepath.Join(root, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command("git", "-C", dir, "init", "-q").CombinedOutput(); err != nil {
			t.Skipf("git init: %v %s", err, out)
		}
		if name == "bare" {
			continue
		}
		c := rootCtx(t, dir)
		callJSON(t, c, "init", `{}`)
		callJSON(t, c, "change.new", `{"intent":"Work in `+name+`"}`)
		ctxs[name] = c
	}
	if err := os.WriteFile(filepath.Join(root, "rhiza.toml"), []byte("[projects.alpha]\npath = \"alpha\"\n[projects.beta]\npath = \"beta\"\n[projects.bare]\npath = \"bare\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// alpha: one run going, one stopped at sign-off (waiting on the owner,
	// and in the inbox) and one failed (in the inbox, not live).
	at := time.Date(2026, 10, 24, 11, 0, 0, 0, time.UTC)
	rs := run.NewStore(run.Dir(filepath.Join(root, "alpha", ".git")))
	mk := func(rid string, es ...run.Entry) {
		t.Helper()
		first := run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "build", Change: "C-" + rid, Title: "Run " + rid}}
		if err := rs.Create(rid, first); err != nil {
			t.Fatal(err)
		}
		for _, e := range es {
			if err := rs.Append(rid, e); err != nil {
				t.Fatal(err)
			}
		}
	}
	enter := func(step string) run.Entry {
		return run.Entry{Type: run.TypeEntered, At: at.Add(time.Second), Step: step, Visit: 1}
	}
	mk("01LIVE", enter("leaves"))
	// A run is going while a process holds its lock.
	unlock, ok, err := rs.TryLock("01LIVE")
	if err != nil || !ok {
		t.Fatalf("lock 01LIVE: %v %v", ok, err)
	}
	t.Cleanup(unlock)
	mk("01SIGN", enter("signoff"), run.Entry{Type: run.TypeStopped, At: at.Add(time.Minute), Step: "signoff", Detail: run.SignoffWait})
	mk("01FAIL", enter("open-pr"),
		run.Entry{Type: run.TypeEnded, At: at.Add(2 * time.Minute), Step: "open-pr", Visit: 1, Outcome: "fail", Detail: "exit 3"},
		run.Entry{Type: run.TypeFinished, At: at.Add(2 * time.Minute), Step: "open-pr", Detail: `step "open-pr" ended with outcome "fail"`, End: run.Failed})

	home := ctxs["alpha"]
	h := api.WebRootHandler(home, root, "alpha", api.WebProjectHandler(home))
	body := serve(h, "GET", "/projects").Body.String()
	// Health lands in the background: wait for both projects' reports.
	for end := time.Now().Add(30 * time.Second); strings.Contains(body, `rz-pill--health rz-pill--"`) && time.Now().Before(end); {
		time.Sleep(50 * time.Millisecond)
		body = serve(h, "GET", "/projects").Body.String()
	}
	rows := map[string]string{}
	for _, row := range strings.Split(body, `<li class="rz-row">`)[1:] {
		for _, name := range []string{"alpha", "beta", "bare"} {
			if strings.Contains(row, `href="/p/`+name+`/"`) {
				rows[name] = row
			}
		}
	}
	if len(rows) != 3 {
		t.Fatalf("want a row for each of alpha, beta, bare; got %v:\n%s", rows, body)
	}
	for name, want := range map[string][]string{
		"alpha": {">alpha<small>1 running, 1 waiting on you · 2 in inbox</small>", "rz-pill--health rz-pill--"},
		"beta":  {">beta<small>none live · 0 in inbox</small>", "rz-pill--health rz-pill--"},
		"bare":  {">bare<small>not initialised</small>", `rz-pill--draft">not initialised`},
	} {
		for _, w := range want {
			if !strings.Contains(rows[name], w) {
				t.Errorf("%s row lacks %q:\n%s", name, w, rows[name])
			}
		}
	}
	for _, name := range []string{"alpha", "beta"} {
		if strings.Contains(rows[name], ">unknown<") || strings.Contains(rows[name], "not initialised") {
			t.Errorf("%s row has no health or is marked bare:\n%s", name, rows[name])
		}
	}
	if strings.Contains(rows["bare"], "in inbox") || strings.Contains(rows["bare"], "rz-pill--health") {
		t.Errorf("bare project shown with counts or health:\n%s", rows["bare"])
	}
	// Each health pill is the worst part of that project's own report.
	o, _ := api.Lookup("health")
	for _, name := range []string{"alpha", "beta"} {
		out, err := o.Run(ctxs[name], nil)
		if err != nil {
			t.Fatal(err)
		}
		if w := `rz-pill--health rz-pill--` + out.(api.HealthResponse).State + `"`; !strings.Contains(rows[name], w) {
			t.Errorf("%s row lacks %q:\n%s", name, w, rows[name])
		}
	}

	// The home is reachable from every page's nav.
	for _, path := range []string{"/", "/p/alpha/runs", "/p/alpha/roadmap", "/inbox", "/p/beta/", "/p/beta/runs"} {
		if b := serve(h, "GET", path).Body.String(); !strings.Contains(b, `<nav class="rz-nav" aria-label="Pages"><a href="/projects"`) {
			t.Errorf("%s does not link the projects home from its nav", path)
		}
	}
	if !strings.Contains(body, `<a href="/projects" aria-current="page">Projects</a>`) {
		t.Errorf("/projects does not mark its nav link current")
	}
}

// settledDir is a temporary directory removed once nothing writes to it
// any more. Every page served starts health refreshes in the background,
// which write each project's index under the cache and read its checkout;
// the test returns with some still running, so a plain t.TempDir's one
// RemoveAll can race them and fail with "directory not empty", or leave
// a directory a late write made again.
func settledDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "rhiza-web-projects-")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		// Removed, it must stay gone for a second: a late refresh can
		// write it again.
		var err error
		gone := time.Time{}
		for end := time.Now().Add(time.Minute); time.Now().Before(end); time.Sleep(100 * time.Millisecond) {
			if _, serr := os.Lstat(dir); serr == nil {
				gone = time.Time{}
				if err = os.RemoveAll(dir); err != nil {
					continue
				}
			}
			if gone.IsZero() {
				gone = time.Now()
			} else if time.Since(gone) >= time.Second {
				return
			}
		}
		t.Errorf("remove %s: %v", dir, err)
	})
	return dir
}

func TestProjectsHomeEmpty(t *testing.T) {
	rec := serve(web.RootWithProjects(web.Handler(nil), nil, nil), "GET", "/projects")
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "No projects") {
		t.Errorf("%d %s", rec.Code, rec.Body)
	}
}
