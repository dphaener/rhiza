package session

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// raceRounds is how many reports each writer makes in a race test.
const raceRounds = 200

// writeReports makes rounds reports as one source (hook or rex), each with
// its round number as the detail, through its own Store. Before each write
// it checks that its previous report is still in the record: a lost update
// would have replaced it with the other writer's stale copy.
func writeReports(dir, id string, src Source, rounds int) error {
	st := NewStore(dir)
	for i := 1; i <= rounds; i++ {
		var lost error
		if _, err := st.Update(id, func(r *Record) {
			mine := &r.Status.Hook
			if src == SourceRex {
				mine = &r.Status.Rex
			}
			if i > 1 && (*mine == nil || (*mine).Detail != strconv.Itoa(i-1)) {
				lost = fmt.Errorf("%s report %d was lost: found %+v", src, i-1, *mine)
			}
			*mine = &Report{State: StateWorking, Detail: strconv.Itoa(i), At: time.Now().UTC()}
		}); err != nil {
			return err
		}
		if lost != nil {
			return lost
		}
	}
	return nil
}

// checkBothKept fails unless each source's last report is in the record.
func checkBothKept(t *testing.T, dir, id string) {
	t.Helper()
	got, err := NewStore(dir).Read(id)
	if err != nil {
		t.Fatal(err)
	}
	want := strconv.Itoa(raceRounds)
	if got.Status.Hook == nil || got.Status.Hook.Detail != want || got.Status.Rex == nil || got.Status.Rex.Detail != want {
		t.Errorf("hook %+v, rex %+v; want both at round %s", got.Status.Hook, got.Status.Rex, want)
	}
}

// TestStoreUpdateConcurrentWriters runs a hook writer and a Rex writer at
// once, each through its own Store, as goroutines. Run it with -race.
func TestStoreUpdateConcurrentWriters(t *testing.T) {
	dir := Dir(t.TempDir())
	if err := NewStore(dir).Write(rec("abcd2222", time.Now())); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	errs := make([]error, 2)
	for i, src := range []Source{SourceHook, SourceRex} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs[i] = writeReports(dir, "abcd2222", src, raceRounds)
		}()
	}
	wg.Wait()
	for _, err := range errs {
		if err != nil {
			t.Error(err)
		}
	}
	checkBothKept(t, dir, "abcd2222")
}

// TestStoreUpdateConcurrentProcesses runs the two writers as separate
// processes, as a hook process and the wrapper's Rex copier are.
func TestStoreUpdateConcurrentProcesses(t *testing.T) {
	dir := Dir(t.TempDir())
	if err := NewStore(dir).Write(rec("abcd3333", time.Now())); err != nil {
		t.Fatal(err)
	}
	var cmds []*exec.Cmd
	var outs []*strings.Builder
	for _, src := range []Source{SourceHook, SourceRex} {
		cmd := exec.Command(os.Args[0], "-test.run=^TestStoreWriterProcess$", "-test.count=1")
		cmd.Env = append(os.Environ(), "RHIZA_TEST_WRITER_DIR="+dir, "RHIZA_TEST_WRITER_SOURCE="+string(src))
		out := &strings.Builder{}
		cmd.Stdout, cmd.Stderr = out, out
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		cmds, outs = append(cmds, cmd), append(outs, out)
	}
	for i, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			t.Errorf("writer process: %v\n%s", err, outs[i])
		}
	}
	checkBothKept(t, dir, "abcd3333")
}

// TestStoreWriterProcess is the writer TestStoreUpdateConcurrentProcesses
// starts. Run on its own, it has nothing to write and passes.
func TestStoreWriterProcess(t *testing.T) {
	dir := os.Getenv("RHIZA_TEST_WRITER_DIR")
	if dir == "" {
		return
	}
	if err := writeReports(dir, "abcd3333", Source(os.Getenv("RHIZA_TEST_WRITER_SOURCE")), raceRounds); err != nil {
		t.Fatal(err)
	}
}
