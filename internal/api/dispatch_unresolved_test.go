package api

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

func workflowSkips(skips []daemon.Skip) []daemon.Skip {
	var out []daemon.Skip
	for _, sk := range skips {
		if sk.Reason == daemon.SkipWorkflow {
			out = append(out, sk)
		}
	}
	return out
}

// A ready change set whose [ready] rule names a workflow that does not
// resolve starts no run (never the default), is skipped naming the workflow
// and its source, and waits in the inbox; once the rule resolves it starts.
func TestDaemonSkipsReadyWithUnresolvedWorkflow(t *testing.T) {
	f := newDaemonFoldFix(t)
	file := filepath.Join(f.ctx.Root, ".rhiza", "dispatch.toml")
	if err := os.WriteFile(file, []byte("[ready]\nfeature = \"ghost\"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	f.poll(t)
	if len(f.picked) != 0 {
		t.Fatalf("started %v under an unresolved workflow", f.picked)
	}
	ws := workflowSkips(f.skips)
	if len(ws) != 1 || ws[0].Change != f.target || ws[0].Detail[0] != "ghost" || !strings.Contains(ws[0].Detail[1], "[ready] rule for feature") {
		t.Fatalf("workflow skips = %+v, want one for the target naming ghost and the [ready] rule", ws)
	}
	items := workflowItems(t.TempDir(), f.skips)
	if len(items) != 1 || items[0].Kind != InboxWorkflow || items[0].Change != f.target ||
		!strings.Contains(items[0].Detail, "ghost") || !strings.Contains(items[0].Detail, "[ready] rule") {
		t.Fatalf("inbox items = %+v", items)
	}
	if strings.Contains(items[0].Detail, f.target) {
		t.Errorf("item text names the change set ID: %s", items[0].Detail)
	}
	if err := os.Remove(file); err != nil {
		t.Fatal(err)
	}
	f.poll(t)
	if !slices.Contains(f.picked, f.target) || len(workflowSkips(f.skips)) != 0 {
		t.Fatalf("after the rule went: picked %v, skips %+v", f.picked, f.skips)
	}
}

// A queued request whose recorded workflow does not resolve starts no plan
// run and is skipped, while another request in the same poll still starts.
func TestDaemonSkipsRequestWithUnresolvedWorkflow(t *testing.T) {
	f := newSessionFixture(t, "unres")
	now := time.Date(2026, 11, 1, 12, 0, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return now }
	f.ctx.repoSession = true
	f.ctx.Self = "/bin/rhiza"
	var detached [][]string
	f.ctx.Detach = func(p Program, _ string) error {
		detached = append(detached, p.Args)
		return nil
	}
	for _, text := range []string{"Bad one", "Good one"} {
		if _, err := callOp(t, f.ctx, "request", `{"text":"`+text+`"}`); err != nil {
			t.Fatal(err)
		}
	}
	// The queue entry is edited directly: request --workflow refuses a name
	// that does not resolve when queued.
	if err := daemon.UpdateRequests(f.common, func(reqs []daemon.Request) ([]daemon.Request, error) {
		reqs[0].Workflow = "ghost"
		return reqs, nil
	}); err != nil {
		t.Fatal(err)
	}
	var skips []daemon.Skip
	rs := run.NewStore(run.Dir(f.common))
	started, err := daemonPlanRequests(f.ctx, rs, nil, 0, 5, "auto", "", &skips)
	if err != nil {
		t.Fatal(err)
	}
	if len(started) != 1 || len(detached) != 1 || !slices.Contains(detached[0], "Good one") {
		t.Fatalf("started %v, detached %v: want only the good request", started, detached)
	}
	if len(skips) != 1 || skips[0].Reason != daemon.SkipWorkflow || skips[0].Request == "" || skips[0].Title != "Bad one" ||
		skips[0].Detail[0] != "ghost" || !strings.Contains(skips[0].Detail[1], "request") {
		t.Fatalf("skips = %+v", skips)
	}
	items := workflowItems(f.common, skips)
	if len(items) != 1 || items[0].Kind != InboxWorkflow || !strings.Contains(items[0].Detail, "Bad one") || !strings.Contains(items[0].Detail, "ghost") {
		t.Fatalf("items = %+v", items)
	}
	if strings.Contains(items[0].Detail, skips[0].Request) {
		t.Errorf("item text names the request ID: %s", items[0].Detail)
	}
}
