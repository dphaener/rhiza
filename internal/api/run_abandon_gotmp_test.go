package api

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
)

// A run whose process died left <runs>/<id>/gotmp behind: abandoning it
// removes the run's folder, and leaves its journal alone.
//
// Serial: it sets its own XDG_CACHE_HOME and inspects the run temp
// folders under it, which every run in the binary shares.
func TestAbandonRemovesRunGoTmp(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	f := newSessionFixture(t, "gotmp")
	rs := run.NewStore(run.Dir(f.common))
	rs.Getenv = f.ctx.Getenv
	if err := rs.Create("run-1", run.Entry{Type: run.TypeStarted, At: time.Now()}); err != nil {
		t.Fatal(err)
	}
	tmp := rs.TmpDir("run-1")
	if err := os.MkdirAll(filepath.Join(tmp, "go-build123"), 0o755); err != nil {
		t.Fatal(err)
	}
	other := run.TmpDir("run-2", f.ctx.Getenv)
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

// A session of a run gets TMPDIR=<run tmp>: what it writes there is gone
// once the run is abandoned.
//
// Serial: see TestAbandonRemovesRunGoTmp.
func TestAbandonRemovesSessionTmpdirFiles(t *testing.T) {
	t.Setenv("XDG_CACHE_HOME", t.TempDir())
	f, s, fs := runFixture(t, `{"permission_mode":"manual","run":"RT-abandon"}`)
	tmp := startTmpdirRun(t, f, s.Change, "RT-abandon", "name = \"wf\"\n"+tmpdirSteps)
	var wrote string
	f.ctx.Start = func(p Program) (Process, error) {
		if p.Terminal {
			dirs := envValues(p.Env, "TMPDIR")
			if len(dirs) != 1 {
				t.Fatalf("agent TMPDIR %v, want one", dirs)
			}
			wrote = filepath.Join(dirs[0], "scratch", "out.txt")
			if err := os.MkdirAll(filepath.Dir(wrote), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(wrote, []byte("x"), 0o644); err != nil {
				t.Fatal(err)
			}
		}
		return fs.start(p)
	}
	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(wrote, tmp+string(filepath.Separator)) {
		t.Fatalf("the session wrote %q, not under the run temp folder %s", wrote, tmp)
	}
	ex := &fakeExec{}
	ctx := *f.ctx
	ctx.Exec = ex.run
	if _, err := callOp(t, &ctx, "run.abandon", `{"run":"RT-abandon"}`); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(tmp); !os.IsNotExist(err) {
		t.Errorf("run folder %s still there after abandon (stat err %v)", tmp, err)
	}
}
