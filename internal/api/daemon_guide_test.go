package api

import (
	"bytes"
	"errors"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/rootcfg"
)

func realExec(name string, args []string, dir string, env []string) ([]byte, error) {
	c := exec.Command(name, args...)
	c.Dir = dir
	out, err := c.CombinedOutput()
	return out, err
}

// oldGuide installs the guide into base for an older rhiza and drops the
// builder skill, as a guide from before the role topics lacks it.
func oldGuide(t *testing.T, ctx *Ctx, base string, project bool) string {
	t.Helper()
	c := *ctx
	c.Version = "0.0.1"
	if project {
		c.Root = base
	} else {
		c.Home = base
	}
	guideInstall(t, &c, `{"agent":"claude","project":`+map[bool]string{true: "true", false: "false"}[project]+`}`)
	dir := filepath.Join(base, ".claude", "skills", "rhiza")
	if err := os.RemoveAll(filepath.Join(dir, "skills", "rhiza-builder")); err != nil {
		t.Fatal(err)
	}
	return dir
}

func guideGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	out, err := realExec("git", append([]string{"-c", "user.email=t@e.x", "-c", "user.name=T"}, args...), dir, nil)
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return string(out)
}

func TestDaemonGuideRefresh(t *testing.T) {
	setup := func(t *testing.T) (*Ctx, string, rootcfg.Config) {
		ctx := testCtx(t.TempDir())
		ctx.Version = "1.2.3"
		ctx.Exec = realExec
		ctx.Home = t.TempDir()
		root := t.TempDir()
		proj := filepath.Join(root, "proj")
		if err := os.MkdirAll(proj, 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := callOp(t, testCtx(proj), "init", `{}`); err != nil {
			t.Fatal(err)
		}
		guideGit(t, proj, "init", "-q")
		reg := rootcfg.Config{Projects: map[string]rootcfg.Project{"proj": {Path: "proj"}}}
		return ctx, root, reg
	}
	noLog := func(string) {}

	t.Run("older guides are rewritten", func(t *testing.T) {
		ctx, root, reg := setup(t)
		proj := filepath.Join(root, "proj")
		home := oldGuide(t, ctx, ctx.Home, false)
		pdir := oldGuide(t, ctx, proj, true)
		var lines []string
		got := refreshGuides(ctx, root, reg, func(l string) { lines = append(lines, l) })
		if len(got) != 2 || !got[0].Refreshed || !got[1].Refreshed || len(lines) != 2 {
			t.Fatalf("got %+v, log %v", got, lines)
		}
		for _, d := range []string{home, pdir} {
			if !pluginCurrent(d, "1.2.3") {
				t.Errorf("%s is not current", d)
			}
		}
		if again := refreshGuides(ctx, root, reg, noLog); len(again) != 0 {
			t.Errorf("a current guide was rewritten: %+v", again)
		}
	})

	t.Run("no guide installs nothing, a foreign folder is left", func(t *testing.T) {
		ctx, root, reg := setup(t)
		foreign := filepath.Join(ctx.Home, ".claude", "skills", "rhiza", "SKILL.md")
		writePath(t, foreign, "mine")
		if got := refreshGuides(ctx, root, reg, noLog); len(got) != 0 {
			t.Fatalf("got %+v", got)
		}
		if b, _ := os.ReadFile(foreign); string(b) != "mine" {
			t.Errorf("foreign file changed: %q", b)
		}
		if _, err := os.Stat(filepath.Join(root, "proj", ".claude")); err == nil {
			t.Error("a guide was installed where none was")
		}
		os.RemoveAll(filepath.Join(ctx.Home, ".claude"))
		if got := refreshGuides(ctx, root, reg, noLog); len(got) != 0 {
			t.Fatalf("got %+v", got)
		}
		if _, err := os.Stat(filepath.Join(ctx.Home, ".claude")); err == nil {
			t.Error("a home with no guide got one")
		}
	})

	t.Run("a committed project guide is not rewritten", func(t *testing.T) {
		ctx, root, reg := setup(t)
		proj := filepath.Join(root, "proj")
		dir := oldGuide(t, ctx, proj, true)
		guideGit(t, proj, "add", "-Af", ".")
		guideGit(t, proj, "commit", "-qm", "guide")
		before := guideGit(t, proj, "status", "--porcelain")
		got := refreshGuides(ctx, root, reg, noLog)
		if len(got) != 1 || got[0].Refreshed || !got[0].Tracked {
			t.Fatalf("got %+v", got)
		}
		if !strings.Contains(got[0].Err, "rhiza guide install --agent claude --project") || got[0].Next() != "rhiza guide install --agent claude --project" {
			t.Errorf("error %q, next %q", got[0].Err, got[0].Next())
		}
		if pluginCurrent(dir, "1.2.3") {
			t.Error("the committed guide was rewritten")
		}
		if after := guideGit(t, proj, "status", "--porcelain"); after != before || strings.TrimSpace(after) != "" {
			t.Errorf("checkout not clean: %q", after)
		}
	})
}

// The daemon refreshes the older guides of the home and of a registered
// project before its first poll, and does it again when it starts as the
// new binary a re-exec runs: at the first poll of each start (ctx.Reap,
// called before a poll starts any session) both installs are already
// current for the running version.
func TestDaemonGuideRefreshAtStart(t *testing.T) {
	ctx, _ := daemonRepo(t)
	ctx.Home = t.TempDir()
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		if name != "git" {
			return nil, errors.New("unexpected " + name)
		}
		return realExec(name, args, dir, env)
	}
	var stderr bytes.Buffer
	ctx.Stderr = &stderr
	// A registered project, with its own git checkout, beside the home.
	proj := filepath.Join(ctx.Root, "proj")
	if err := os.MkdirAll(proj, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := callOp(t, testCtx(proj), "init", `{}`); err != nil {
		t.Fatal(err)
	}
	guideGit(t, proj, "init", "-q")
	writePath(t, filepath.Join(ctx.Root, "rhiza.toml"), "[projects.repo]\npath = \".\"\n\n[projects.proj]\npath = \"proj\"\n")
	ctx.Version = "1.2.3"
	home := oldGuide(t, ctx, ctx.Home, false)
	pdir := oldGuide(t, ctx, proj, true)

	start := daemon.Image{Size: 10, Mod: time.Unix(1700000000, 0)}
	old := daemonReexec
	t.Cleanup(func() { daemonReexec = old })

	// atFirstPoll runs the daemon and reports, for the first poll, whether
	// both guides were current for the running version when it began.
	atFirstPoll := func(f *reexecFake) (homeOK, projOK bool) {
		t.Helper()
		polled := false
		ctx.Reap = func(string) {
			if !polled {
				polled = true
				homeOK, projOK = pluginCurrent(home, ctx.Version), pluginCurrent(pdir, ctx.Version)
			}
		}
		daemonReexec = func() *daemon.Reexec { return f.reexec(start) }
		if err := runDaemonTicks(t, ctx, 1, nil); err != nil {
			t.Fatal(err)
		}
		if !polled {
			t.Fatal("the daemon never polled")
		}
		return homeOK, projOK
	}

	// A restart: the older guides are rewritten before the first poll, and
	// the rebuilt binary it finds then is re-executed.
	f := &reexecFake{img: daemon.Image{Size: 11, Mod: start.Mod}}
	if h, p := atFirstPoll(f); !h || !p {
		t.Fatalf("at the first poll after a restart: home current %v, project current %v", h, p)
	}
	if f.execs != 1 {
		t.Fatalf("execs = %d, want 1", f.execs)
	}
	for _, want := range []string{"refreshed the user guide", "refreshed the project proj guide"} {
		if !strings.Contains(stderr.String(), want) {
			t.Errorf("log lacks %q:\n%s", want, stderr.String())
		}
	}

	// The new image the re-exec ran starts as a newer version, over the
	// guides the older one installed: it rewrites both before its first
	// poll, and logs each.
	stderr.Reset()
	ctx.Version = "1.2.4"
	if h, p := atFirstPoll(&reexecFake{img: start}); !h || !p {
		t.Fatalf("at the first poll after the re-exec: home current %v, project current %v", h, p)
	}
	if !strings.Contains(stderr.String(), "refreshed the user guide in "+home+" to 1.2.4") ||
		!strings.Contains(stderr.String(), "refreshed the project proj guide in "+pdir+" to 1.2.4") {
		t.Errorf("log after the re-exec:\n%s", stderr.String())
	}

	// A start over current guides rewrites neither.
	stderr.Reset()
	if h, p := atFirstPoll(&reexecFake{img: start}); !h || !p {
		t.Fatalf("current guides at a later start: home %v, project %v", h, p)
	}
	if strings.Contains(stderr.String(), "guide") {
		t.Errorf("a current guide was refreshed:\n%s", stderr.String())
	}
}

// When the refresh fails for a scope, the daemon starts no step whose start
// command names a skill missing from the guide, opens one inbox item naming
// the scope and the install command however many polls pass, and starts the
// step on the poll after the guide is installed, with the item gone.
func TestDaemonGuideHold(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root writes into an unwritable folder")
	}
	f := usageFix(t)
	ctx := f.ctx
	ctx.Version = "1.2.3"
	ctx.Home = t.TempDir()
	home := oldGuide(t, ctx, ctx.Home, false)
	// The refresh cannot write the missing skill's folder.
	skills := filepath.Join(home, "skills")
	if err := os.Chmod(skills, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(skills, 0o755) })
	var lines []string
	scopes := refreshGuides(ctx, t.TempDir(), rootcfg.Config{}, func(l string) { lines = append(lines, l) })
	if len(scopes) != 1 || scopes[0].Refreshed || scopes[0].Err == "" {
		t.Fatalf("refresh: %+v, log %v", scopes, lines)
	}
	ctx.guideGate = newGuideGate(ctx, scopes)

	for poll := range 3 {
		var skips []daemon.Skip
		started, _, _, err := daemonPoll1(ctx, f.rs, f.pending, f.unj, 0, 4, 0, false, "", false, nil, nil, &skips)
		if err != nil {
			t.Fatal(err)
		}
		if len(started)+len(f.picked) != 0 {
			t.Fatalf("poll %d started %v (picked %v) over a stale guide", poll, started, f.picked)
		}
		var held bool
		for _, sk := range skips {
			if sk.Change == f.target && sk.Reason == daemon.SkipGuideStale {
				held = strings.Contains(strings.Join(sk.Detail, " "), "/rhiza-builder")
			}
		}
		if !held {
			t.Fatalf("poll %d skips %+v: want %s held for the guide", poll, skips, f.target)
		}
		items := guideHoldItems(ctx.guideGate.state(ctx.Root))
		if len(items) != 1 || items[0].Kind != InboxGuideStale {
			t.Fatalf("poll %d items %+v, want exactly one", poll, items)
		}
		d := items[0].Detail
		if !strings.Contains(d, "the user guide") || !strings.Contains(d, scopes[0].Err) || !strings.Contains(d, "rhiza guide install --agent claude") || strings.Contains(d, "--project") {
			t.Fatalf("item %q", d)
		}
	}

	// The owner installs the guide: the next poll starts the step.
	if err := os.Chmod(skills, 0o755); err != nil {
		t.Fatal(err)
	}
	guideInstall(t, ctx, `{"agent":"claude"}`)
	var skips []daemon.Skip
	started, _, _, err := daemonPoll1(ctx, f.rs, f.pending, f.unj, 0, 4, 0, false, "", false, nil, nil, &skips)
	if err != nil {
		t.Fatal(err)
	}
	if len(started) != 1 || started[0] != f.target {
		t.Fatalf("after the install started %v, want %s", started, f.target)
	}
	if items := guideHoldItems(ctx.guideGate.state(ctx.Root)); len(items) != 0 {
		t.Fatalf("items after the install: %+v", items)
	}
}
