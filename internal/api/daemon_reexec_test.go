package api

import (
	"errors"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
)

// reexecFake is a Reexec whose executable's Image is img and whose Exec
// counts its calls and runs onExec. Exec returns err, as a failed exec
// would; a fake's nil return stands for the new image having run.
type reexecFake struct {
	img    daemon.Image
	stats  int
	execs  int
	err    error
	onExec func()
}

func (f *reexecFake) reexec(start daemon.Image) *daemon.Reexec {
	r := &daemon.Reexec{
		Path: "/bin/rhiza", Args: []string{"rhiza", "daemon"},
		Stat: func(string) (daemon.Image, error) { f.stats++; return f.img, nil },
		Exec: func(string, []string, []string) error {
			f.execs++
			if f.onExec != nil {
				f.onExec()
			}
			return f.err
		},
	}
	r.Start(start)
	return r
}

// newImageTakesLock checks, from inside an Exec, that the daemon let go of
// its lock, so that the new image takes it.
func newImageTakesLock(t *testing.T, common string) {
	l, err := daemon.TryLock(common, time.Now())
	if err != nil {
		t.Errorf("the new image cannot take the daemon lock: %v", err)
		return
	}
	l.Unlock()
}

// runDaemonTicks runs the daemon op on ctx with the ticker and Interrupts
// replaced, sends it n ticks (n+1 polls: one at start, one per tick), runs
// between after each tick, then interrupts it and returns its error.
// An unbuffered tick is taken only once a poll and its re-exec check are
// done, so each send returns after the poll before it.
func runDaemonTicks(t *testing.T, ctx *Ctx, n int, between func(tick int)) error {
	t.Helper()
	ticks := make(chan time.Time)
	old := daemonTicker
	daemonTicker = func() (<-chan time.Time, func()) { return ticks, func() {} }
	t.Cleanup(func() { daemonTicker = old })
	stop := make(chan struct{})
	ctx.Interrupts = func() (<-chan struct{}, func()) { return stop, func() {} }
	done := make(chan error, 1)
	go func() {
		_, err := callOp(t, ctx, "daemon", `{"permission_mode":"manual"}`)
		done <- err
	}()
	for i := 1; i <= n; i++ {
		select {
		case ticks <- time.Now():
		case err := <-done:
			t.Fatalf("daemon stopped before tick %d: %v", i, err)
		}
		if between != nil {
			between(i)
		}
	}
	close(stop)
	select {
	case err := <-done:
		return err
	case <-time.After(10 * time.Second):
		t.Fatal("daemon did not stop on interrupt")
		return nil
	}
}

// daemonRepo is a repository with an empty model, so a poll starts
// nothing, and the git common dir it reports.
func daemonRepo(t *testing.T) (*Ctx, string) {
	t.Helper()
	ctx := testCtx(t.TempDir())
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	common := t.TempDir()
	ctx.Git = func(args ...string) (string, error) {
		if len(args) == 2 && args[0] == "rev-parse" && args[1] == "--git-common-dir" {
			return common + "\n", nil
		}
		return "", errors.New("unexpected git " + args[0])
	}
	return ctx, common
}

// The daemon re-executes itself when its executable changed since it
// started: once, at the end of its first poll, with its lock let go for the
// new image and taken back when the exec returns. An unchanged executable
// never re-executes.
func TestDaemonReexecsRebuiltBinary(t *testing.T) {
	start := daemon.Image{Size: 10, Mod: time.Unix(1700000000, 0)}
	old := daemonReexec
	t.Cleanup(func() { daemonReexec = old })

	t.Run("changed", func(t *testing.T) {
		ctx, common := daemonRepo(t)
		f := &reexecFake{img: daemon.Image{Size: 11, Mod: start.Mod}}
		f.onExec = func() {
			// After the poll: the poll has written the daemon's state.
			if s, ok := daemon.ReadState(common); !ok || s.LastPoll.IsZero() {
				t.Errorf("exec before the first poll finished: state %+v, %v", s, ok)
			}
			newImageTakesLock(t, common)
		}
		daemonReexec = func() *daemon.Reexec { return f.reexec(start) }
		err := runDaemonTicks(t, ctx, 3, func(tick int) {
			if _, held := daemon.Running(common); !held {
				t.Errorf("tick %d: daemon lock not held again after the exec", tick)
			}
		})
		if err != nil {
			t.Fatal(err)
		}
		if f.execs != 1 {
			t.Fatalf("execs over 4 polls = %d, want 1", f.execs)
		}
	})

	t.Run("unchanged", func(t *testing.T) {
		ctx, _ := daemonRepo(t)
		f := &reexecFake{img: start}
		daemonReexec = func() *daemon.Reexec { return f.reexec(start) }
		if err := runDaemonTicks(t, ctx, 3, nil); err != nil {
			t.Fatal(err)
		}
		if f.execs != 0 || f.stats != 4 {
			t.Fatalf("execs %d, stats %d over 4 polls, want 0 and 4", f.execs, f.stats)
		}
	})
}

// A start pending at the end of a poll delays the exec to the end of the
// first later poll with none pending. An exec that returns, failed or not,
// is not tried again, and the lock is held again.
func TestDaemonLoopReexecWaitsForPending(t *testing.T) {
	for _, execErr := range []error{nil, errors.New("exec format error")} {
		common := t.TempDir()
		lock, err := daemon.TryLock(common, time.Now())
		if err != nil {
			t.Fatal(err)
		}
		start := daemon.Image{Size: 10, Mod: time.Unix(1700000000, 0)}
		f := &reexecFake{img: daemon.Image{Size: 12, Mod: start.Mod}, err: execErr}
		polls, execAt := 0, 0
		f.onExec = func() {
			execAt = polls
			newImageTakesLock(t, common)
		}
		// Polls 1 and 2 leave a start pending, polls 3 to 5 do not.
		pendingAt := map[int]int{1: 1, 2: 2}
		ctx := testCtx(t.TempDir())
		ticks, stop := make(chan time.Time), make(chan struct{})
		type result struct {
			lock *daemon.Lock
			err  error
		}
		done := make(chan result, 1)
		go func() {
			l, err := daemonLoop(ctx, common, lock,
				func() error { polls++; return nil },
				func() int { return pendingAt[polls] },
				f.reexec(start), ticks, stop)
			done <- result{l, err}
		}()
		for range 4 {
			ticks <- time.Now()
		}
		close(stop)
		res := <-done
		if res.err != nil || res.lock == nil {
			t.Fatalf("exec error %v: loop = %v, %v", execErr, res.lock, res.err)
		}
		if polls != 5 || f.execs != 1 || execAt != 3 {
			t.Fatalf("exec error %v: polls %d, execs %d at poll %d; want 5 polls, 1 exec at poll 3", execErr, polls, f.execs, execAt)
		}
		if _, held := daemon.Running(common); !held {
			t.Fatalf("exec error %v: lock not held after the loop", execErr)
		}
		res.lock.Unlock()
	}
}

// A daemon that takes over, as a re-exec's new image does, keeps the
// unjournaled holds of the daemon before it: their start, and the runs
// their change set had then.
func TestDaemonKeepsUnjournaledHolds(t *testing.T) {
	at := time.Unix(1700000000, 0).UTC()
	prev := []daemon.Event{{Change: "B", At: at, Runs: 2}, {Change: "A", At: at.Add(time.Minute)}}
	got := daemonRestoreUnjournaled(prev)
	if len(got) != 2 || got["B"] != (daemonPending{at: at, runs: 2}) || got["A"] != (daemonPending{at: at.Add(time.Minute)}) {
		t.Fatalf("restored %v", got)
	}
	// Restored, the hold runs its 15 minutes from the first start.
	yes := func(string) bool { return true }
	daemonHeld(map[string]daemonPending{}, got, nil, yes, at.Add(daemonUnjournaledHold-time.Second))
	if _, ok := got["B"]; !ok {
		t.Fatalf("restored hold of B released before 15 minutes")
	}

}
