package api

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// Two change sets with one title: the later one's ship target and abandon
// cleanup name its suffixed branch, not the earlier one's.
func TestUniqueSlugShipTargetAndAbandon(t *testing.T) {
	f := newSessionFixture(t, "twin")
	first := newChange(t, f.ctx, `{"intent":"Twin work","title":"Twin work"}`).ChangeSet
	second := newChange(t, f.ctx, `{"intent":"Twin work","title":"Twin work"}`).ChangeSet
	if first.String() >= second.String() {
		t.Fatalf("IDs out of order: %s, %s", first, second)
	}
	suffix := strings.ToLower(second.String()[:12])

	for _, tc := range []struct {
		cs   string
		want string
	}{{first.String(), "twin-work"}, {second.String(), "twin-work-" + suffix}} {
		ctx := *f.ctx
		ctx.ChangeRef = tc.cs
		tgt, err := findShipTarget(&ctx)
		if err != nil {
			t.Fatal(err)
		}
		if tgt.slug != tc.want || tgt.branch != "cs/"+tc.want || filepath.Base(tgt.worktree) != "repo."+tc.want {
			t.Errorf("ship target of %s = %q %q %q, want slug %q", tc.cs, tgt.slug, tgt.branch, tgt.worktree, tc.want)
		}
	}

	ex := &fakeExec{}
	ctx := *f.ctx
	ctx.Exec = ex.run
	st := run.State{Change: second.String(), Title: "Twin work", Worktree: filepath.Join(filepath.Dir(f.worktree), "repo.gone")}
	rs := run.NewStore(run.Dir(f.common))
	abandonCleanup(&ctx, f.common, rs, "run1", st)
	var showed bool
	for _, c := range ex.calls {
		if strings.Contains(strings.Join(c, " "), "refs/heads/"+session.Branch("twin-work-"+suffix)) {
			showed = true
		}
		if strings.HasSuffix(strings.Join(c, " "), "refs/heads/cs/twin-work") {
			t.Errorf("cleanup named the earlier change set's branch: %v", c)
		}
	}
	if !showed {
		t.Errorf("cleanup never named cs/twin-work-%s: %v", suffix, ex.calls)
	}
}
