package api

import (
	"os"
	"slices"
	"testing"

	"github.com/dphaener/rhiza/internal/run"
)

// A run started in Rex and resumed from another block journals a process
// entry each time: this host and pid, the block from REX_BLOCK and the
// window `rex block inspect` names for it. Outside Rex the entry has no
// window or block.
func TestRunJournalsProcess(t *testing.T) {
	f := newSessionFixture(t, "x")
	newChange(t, f.ctx, `{"intent":"X","title":"X"}`)
	writeWorkflow(t, f.ctx.Root, "w", `
[step.gate]
kind = "wait"
until = "ready"
timeout = "1m"
`)
	sh := &shStart{codes: map[string]int{"ready": 1}}
	f.ctx.Start = sh.start
	windows := map[string]string{"block:start": "window:ws", "block:resume": "window:wr"}
	var inspected [][]string
	f.ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		inspected = append(inspected, append([]string{name}, args...))
		block := args[len(args)-1]
		return []byte(`{"block":{"id":"` + block + `","window_id":"` + windows[block] + `"}}`), nil
	}
	host, _ := os.Hostname()
	store := run.NewStore(run.Dir(f.common))
	processes := func(id string) []run.Process {
		t.Helper()
		j, err := store.Read(id)
		if err != nil {
			t.Fatal(err)
		}
		var ps []run.Process
		for _, e := range j.Entries {
			if e.Type == run.TypeProcess {
				ps = append(ps, *e.Process)
			}
		}
		return ps
	}

	f.ctx.RexSession, f.ctx.RexBlock = "session:s", "block:start"
	f.ctx.repoSession = true // keep the caller's session: this checks the process record
	resp, err := callOp(t, f.ctx, "run.start", `{"foreground":true,"workflow":"w","no_wait":true}`)
	if err != nil {
		t.Fatal(err)
	}
	id := resp.(RunResponse).Run.ID
	started := run.Process{Host: host, Pid: os.Getpid(), RexSession: "session:s", RexWindow: "window:ws", RexBlock: "block:start"}
	if ps := processes(id); !slices.Equal(ps, []run.Process{started}) {
		t.Fatalf("started run's processes = %+v, want %+v", ps, started)
	}
	if want := []string{"rex", "block", "inspect", "-s", "session:s", "--json", "block:start"}; len(inspected) != 1 || !slices.Equal(inspected[0], want) {
		t.Errorf("rex calls = %q, want %q", inspected, want)
	}

	f.ctx.RexBlock = "block:resume"
	if _, err := callOp(t, f.ctx, "run.resume", `{"foreground":true,"run":"`+id+`","no_wait":true}`); err != nil {
		t.Fatal(err)
	}
	resumed := run.Process{Host: host, Pid: os.Getpid(), RexSession: "session:s", RexWindow: "window:wr", RexBlock: "block:resume"}
	if ps := processes(id); !slices.Equal(ps, []run.Process{started, resumed}) {
		t.Fatalf("resumed run's processes = %+v, want %+v then %+v", ps, started, resumed)
	}
	j, err := store.Read(id)
	if err != nil {
		t.Fatal(err)
	}
	st := j.State()
	if st.Process == nil || *st.Process != resumed || st.Liveness(host, nil) != run.LiveRunning {
		t.Errorf("state process = %+v, liveness %s", st.Process, st.Liveness(host, nil))
	}

	// Outside Rex: host and pid alone, and Rex is never asked.
	f.ctx.RexSession, f.ctx.RexBlock = "", ""
	inspected = nil
	if _, err := callOp(t, f.ctx, "run.resume", `{"foreground":true,"run":"`+id+`","no_wait":true}`); err != nil {
		t.Fatal(err)
	}
	bare := run.Process{Host: host, Pid: os.Getpid()}
	if ps := processes(id); len(ps) != 3 || ps[2] != bare {
		t.Errorf("processes outside Rex = %+v, want a third %+v", ps, bare)
	}
	if len(inspected) != 0 {
		t.Errorf("outside Rex, rex ran %q", inspected)
	}
}
