package api

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// A run whose process died left <runs>/<id>/gotmp behind: abandoning it
// removes the run's folder, and leaves its journal alone.
func TestAbandonRemovesRunGoTmp(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	f := newSessionFixture(t, "gotmp")
	rs := run.NewStore(run.Dir(f.common))
	if err := rs.Create("run-1", run.Entry{Type: run.TypeStarted, At: time.Now()}); err != nil {
		t.Fatal(err)
	}
	tmp := rs.TmpDir("run-1")
	if err := os.MkdirAll(filepath.Join(tmp, "go-build123"), 0o755); err != nil {
		t.Fatal(err)
	}
	other := run.TmpDir("run-2")
	if err := os.MkdirAll(other, 0o755); err != nil {
		t.Fatal(err)
	}
	ex := &fakeExec{}
	ctx := *f.ctx
	ctx.Exec = ex.run
	if _, err := callOp(t, &ctx, "run.abandon", `{"run":"run-1"}`); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Errorf("run folder %s still there after abandon (stat err %v)", tmp, err)
	}
	if _, err := os.Stat(other); err != nil {
		t.Errorf("another run's folder gone after abandon: %v", err)
	}
	if _, err := os.Stat(filepath.Dir(tmp)); err != nil {
		t.Errorf("shared run-tmp parent gone after abandon: %v", err)
	}
	j, err := rs.Read("run-1")
	if err != nil {
		t.Fatal(err)
	}
	if st := j.State(); !st.Ended || !st.Abandoned() {
		t.Errorf("run state %+v: want ended and abandoned", st)
	}
}
