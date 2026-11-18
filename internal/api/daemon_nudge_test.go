package api

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// TestDaemonNudgeStalled checks one autonomous and one owner-attended run
// under the same daemon: the autonomous run's session, idle at its prompt
// without saying done, stalls and gets one
// message per stall, keyed on its report time; the attended run's gets none;
// the inbox's stalled list is untouched by the nudge.
func TestDaemonNudgeStalled(t *testing.T) {
	t.Parallel()
	sessions := session.NewStore(t.TempDir())
	t0 := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
	report := func(id string, at time.Time) {
		// The autonomous session ended its turn: its hook reports idle,
		// which stalls it because its record is autonomous.
		state := session.StateWorking
		if id == "auto" {
			state = session.StateIdle
		}
		rec := session.Record{ID: id, Autonomous: id == "auto"}
		rec.Status.Hook = &session.Report{State: state, At: at}
		if err := sessions.Write(rec); err != nil {
			t.Fatal(err)
		}
	}
	report("auto", t0)
	report("owner", t0)
	states := []run.State{
		{ID: "A", Autonomous: true, Current: "build", Session: map[string]string{"build": "auto"}},
		{ID: "O", Current: "build", Session: map[string]string{"build": "owner"}},
	}
	var sent []string
	send := func(sid, text string) error {
		// The message names background shells and `rhiza session done`.
		for _, want := range []string{"background shell", "`rhiza session done`"} {
			if !strings.Contains(text, want) {
				t.Errorf("nudge to %s lacks %q: %q", sid, want, text)
			}
		}
		sent = append(sent, sid)
		return nil
	}
	var nudged []string
	poll := func(now time.Time) {
		t.Helper()
		var err error
		if nudged, err = daemonNudgeStalled(states, sessions, now, nudged, send); err != nil {
			t.Fatal(err)
		}
	}
	inbox := func(now time.Time) int {
		n := 0
		for _, st := range states {
			n += len(runStalledAll(st, true, sessions, now))
		}
		return n
	}

	poll(t0.Add(19 * time.Minute))
	if len(sent) != 0 {
		t.Fatalf("before the stall: sent %v", sent)
	}
	poll(t0.Add(21 * time.Minute))
	if !slices.Equal(sent, []string{"auto"}) {
		t.Fatalf("first stall: sent %v", sent)
	}
	poll(t0.Add(25 * time.Minute))
	if len(sent) != 1 {
		t.Fatalf("same stall again: sent %v", sent)
	}
	if n := inbox(t0.Add(25 * time.Minute)); n != 2 {
		t.Errorf("inbox lists %d stalled sessions, want 2", n)
	}
	// A new report, then another 21 silent minutes: a second stall.
	t1 := t0.Add(30 * time.Minute)
	report("auto", t1)
	poll(t1.Add(time.Minute))
	poll(t1.Add(21 * time.Minute))
	if !slices.Equal(sent, []string{"auto", "auto"}) {
		t.Fatalf("second stall: sent %v", sent)
	}
	// A failed send is retried by the next poll.
	report("auto", t1.Add(40*time.Minute))
	fail := true
	send = func(sid, text string) error {
		if fail {
			return errTest
		}
		sent = append(sent, sid)
		return nil
	}
	poll2 := t1.Add(62 * time.Minute)
	if _, err := daemonNudgeStalled(states, sessions, poll2, nudged, send); err == nil {
		t.Fatal("want the send error")
	}
	fail = false
	poll(poll2)
	if len(sent) != 3 {
		t.Fatalf("retry: sent %v", sent)
	}
}

// TestDaemonNudgeStalledSpawnedChild checks the daemon's nudge of a
// conduct run whose conductor spawned a fixer: the fixer, idle without
// saying done, is nudged and the conductor waiting on it is not; a fixer
// that said done or whose spawn ended is not nudged.
func TestDaemonNudgeStalledSpawnedChild(t *testing.T) {
	t.Parallel()
	t0 := time.Date(2026, 11, 16, 0, 0, 0, 0, time.UTC)
	now := t0.Add(21 * time.Minute)
	sessions := session.NewStore(t.TempDir())
	write := func(id, state string, done bool) {
		rec := session.Record{ID: id, Autonomous: true}
		rec.Status.Hook = &session.Report{State: state, At: t0}
		if done {
			rec.Status.Done = &session.Report{State: session.StateDone, At: t0}
		}
		if err := sessions.Write(rec); err != nil {
			t.Fatal(err)
		}
	}
	write("conductor", session.StateWorking, false)
	write("fixer", session.StateIdle, false)
	write("finished", session.StateIdle, true)
	write("ended", session.StateIdle, false)
	st := func(id, child string, ended bool) run.State {
		s := run.State{ID: id, Autonomous: true, Current: "conduct", Session: map[string]string{"conduct": "conductor"},
			Spawns: []run.Spawn{{Step: "conduct", Parent: "conductor", Session: child, Role: "fixer", Ended: ended}}}
		return s
	}
	cases := []struct {
		name string
		st   run.State
		want []string
	}{
		{"idle child", st("A", "fixer", false), []string{"fixer"}},
		{"child said done", st("B", "finished", false), nil},
		{"spawn ended", st("C", "ended", true), []string{"conductor"}},
	}
	for _, c := range cases {
		var sent []string
		send := func(sid, _ string) error { sent = append(sent, sid); return nil }
		if _, err := daemonNudgeStalled([]run.State{c.st}, sessions, now, nil, send); err != nil {
			t.Fatal(err)
		}
		if !slices.Equal(sent, c.want) {
			t.Errorf("%s: sent %v, want %v", c.name, sent, c.want)
		}
	}
}

type testErr string

func (e testErr) Error() string { return string(e) }

const errTest = testErr("send failed")
