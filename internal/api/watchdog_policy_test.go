package api

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/workflow"
)

// watchdogCtx is a context under a fresh root whose project git config is
// the map cfg.
func watchdogCtx(t *testing.T, cfg map[string]string) (*Ctx, string) {
	t.Helper()
	root := t.TempDir()
	if _, err := rootcfg.Init(root); err != nil {
		t.Fatal(err)
	}
	return &Ctx{Wd: root, Git: func(args ...string) (string, error) {
		if len(args) == 3 && args[0] == "config" && args[1] == "--get" {
			if v, ok := cfg[args[2]]; ok {
				return v + "\n", nil
			}
		}
		return "", errors.New("unset")
	}}, root
}

func TestWatchdogPolicy(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		ctx, _ := watchdogCtx(t, nil)
		p, err := WatchdogPolicyOf(ctx)
		if err != nil {
			t.Fatal(err)
		}
		want := WatchdogPolicy{true, 20 * time.Minute, 40 * time.Minute, 30 * time.Minute, 90 * time.Minute, 2, 1, 3, true, true}
		if p != want {
			t.Fatalf("got %+v, want %+v", p, want)
		}
	})
	t.Run("machine value", func(t *testing.T) {
		ctx, root := watchdogCtx(t, nil)
		if err := rootcfg.SetDaemon(root, "rhiza.daemon.watchdog.retryBudget", "5"); err != nil {
			t.Fatal(err)
		}
		if err := rootcfg.SetDaemon(root, "rhiza.daemon.watchdog.waitLimit", "2h"); err != nil {
			t.Fatal(err)
		}
		p, err := WatchdogPolicyOf(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if p.RetryBudget != 5 || p.WaitLimit != 2*time.Hour || p.CommandLimit != DefaultCommandLimit {
			t.Fatalf("got %+v", p)
		}
	})
	t.Run("project over machine", func(t *testing.T) {
		ctx, root := watchdogCtx(t, map[string]string{
			"rhiza.project.watchdog.retryBudget": "0",
			"rhiza.project.watchdog.enabled":     "false",
		})
		if err := rootcfg.SetDaemon(root, "rhiza.daemon.watchdog.retryBudget", "5"); err != nil {
			t.Fatal(err)
		}
		if err := rootcfg.SetDaemon(root, "rhiza.daemon.watchdog.ciReruns", "4"); err != nil {
			t.Fatal(err)
		}
		p, err := WatchdogPolicyOf(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if p.RetryBudget != 0 || p.Enabled || p.CIReruns != 4 {
			t.Fatalf("got %+v", p)
		}
	})
	t.Run("no project pollBehind", func(t *testing.T) {
		if _, ok := findSetting("rhiza.project.watchdog.pollBehind"); ok {
			t.Fatal("pollBehind is machine only")
		}
		if d, ok := findSetting("rhiza.daemon.watchdog.pollBehind"); !ok || d.project {
			t.Fatal("pollBehind machine key missing")
		}
	})
	t.Run("step stall wins", func(t *testing.T) {
		def, err := workflow.Parse("w", []byte("name = \"w\"\nstart = \"a\"\n[step.a]\nkind = \"command\"\nrun = \"true\"\nstall = \"7m\"\n[step.b]\nkind = \"command\"\nrun = \"true\"\n"))
		if err != nil {
			t.Fatal(err)
		}
		p := WatchdogPolicy{CommandLimit: 30 * time.Minute, AgentRetry: 40 * time.Minute}
		a := def.Steps["a"]
		if got := p.StallLimit(a.Kind, a.Stall); got != 7*time.Minute {
			t.Fatalf("step a limit %v", got)
		}
		b := def.Steps["b"]
		if got := p.StallLimit(b.Kind, b.Stall); got != 30*time.Minute {
			t.Fatalf("step b limit %v", got)
		}
	})
	t.Run("invalid values name the key", func(t *testing.T) {
		for key, bad := range map[string]string{
			"rhiza.daemon.watchdog.agentNudge":    "soon",
			"rhiza.daemon.watchdog.retryBudget":   "-1",
			"rhiza.daemon.watchdog.pollBehind":    "0",
			"rhiza.daemon.watchdog.enabled":       "maybe",
			"rhiza.project.watchdog.commandLimit": "0s",
		} {
			ctx, root := watchdogCtx(t, nil)
			if strings.HasPrefix(key, "rhiza.project.") {
				ctx, _ = watchdogCtx(t, map[string]string{key: bad})
			} else if err := os.WriteFile(filepath.Join(root, rootcfg.File), []byte("[daemon]\n\""+key+"\" = \""+bad+"\"\n"), 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := WatchdogPolicyOf(ctx)
			if err == nil || !strings.Contains(err.Error(), key) {
				t.Errorf("%s=%s: error %v", key, bad, err)
			}
		}
	})
}
