package api

import (
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
		got := refreshGuides(ctx, root, reg, "proj", func(l string) { lines = append(lines, l) })
		if len(got) != 2 || !got[0].Refreshed || !got[1].Refreshed || len(lines) != 2 {
			t.Fatalf("got %+v, log %v", got, lines)
		}
		for _, d := range []string{home, pdir} {
			if !pluginCurrent(d, "1.2.3") {
				t.Errorf("%s is not current", d)
			}
		}
		if again := refreshGuides(ctx, root, reg, "proj", noLog); len(again) != 0 {
			t.Errorf("a current guide was rewritten: %+v", again)
		}
	})

	t.Run("no guide installs nothing, a foreign folder is left", func(t *testing.T) {
		ctx, root, reg := setup(t)
		foreign := filepath.Join(ctx.Home, ".claude", "skills", "rhiza", "SKILL.md")
		writePath(t, foreign, "mine")
		if got := refreshGuides(ctx, root, reg, "proj", noLog); len(got) != 0 {
			t.Fatalf("got %+v", got)
		}
		if b, _ := os.ReadFile(foreign); string(b) != "mine" {
			t.Errorf("foreign file changed: %q", b)
		}
		if _, err := os.Stat(filepath.Join(root, "proj", ".claude")); err == nil {
			t.Error("a guide was installed where none was")
		}
		os.RemoveAll(filepath.Join(ctx.Home, ".claude"))
		if got := refreshGuides(ctx, root, reg, "proj", noLog); len(got) != 0 {
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
		got := refreshGuides(ctx, root, reg, "proj", noLog)
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

// The daemon refreshes an older user guide before its first poll.
func TestDaemonGuideRefreshAtStart(t *testing.T) {
	ctx, _ := daemonRepo(t)
	ctx.Version = "1.2.3"
	ctx.Exec = realExec
	ctx.Home = t.TempDir()
	ctx.IndexPath = filepath.Join(t.TempDir(), "index.db")
	dir := oldGuide(t, ctx, ctx.Home, false)
	if _, err := callOp(t, ctx, "daemon", `{"once":true,"permission_mode":"manual"}`); err != nil {
		t.Fatal(err)
	}
	if !pluginCurrent(dir, "1.2.3") {
		t.Error("the daemon did not refresh the guide")
	}
}
