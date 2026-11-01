package session

import (
	"testing"
	"time"
)

func TestResolvePrecedence(t *testing.T) {
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	hook := &Report{State: StateWorking, At: at}
	rex := &Report{State: StateIdle, At: at.Add(time.Hour)}
	code := 1
	exit := &Report{State: StateExited, Code: &code, At: at.Add(2 * time.Hour)}
	cases := []struct {
		name   string
		status Status
		state  string
		src    Source
	}{
		{"none", Status{}, StateStarting, SourceDefault},
		{"hook alone", Status{Hook: hook}, StateWorking, SourceHook},
		{"rex alone", Status{Rex: rex}, StateIdle, SourceRex},
		{"exit alone", Status{Exit: exit}, StateExited, SourceExit},
		{"hook over rex", Status{Hook: hook, Rex: rex}, StateWorking, SourceHook},
		{"exit over hook", Status{Hook: hook, Exit: exit}, StateExited, SourceExit},
		{"exit over rex", Status{Rex: rex, Exit: exit}, StateExited, SourceExit},
		{"exit over all", Status{Hook: hook, Rex: rex, Exit: exit}, StateExited, SourceExit},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, src := Resolve(Record{StartedAt: at, Status: c.status})
			if got.State != c.state || src != c.src {
				t.Errorf("Resolve = %s from %s, want %s from %s", got.State, src, c.state, c.src)
			}
		})
	}
}

func TestResolveKeepsDetails(t *testing.T) {
	code := 2
	r := Record{Status: Status{Rex: &Report{State: StateBlocked, Kind: KindPermission, Detail: "Bash"}}}
	got, _ := Resolve(r)
	if got.Kind != KindPermission || got.Detail != "Bash" {
		t.Errorf("got %+v", got)
	}
	r.Status.Exit = &Report{State: StateError, Code: &code}
	got, _ = Resolve(r)
	if got.Code == nil || *got.Code != 2 {
		t.Errorf("exit code lost: %+v", got)
	}
}

func TestResolveValidState(t *testing.T) {
	for _, s := range []string{StateStarting, StateWorking, StateIdle, StateBlocked, StateDone, StateError, StateExited} {
		if !ValidState(s) {
			t.Errorf("%s should be valid", s)
		}
	}
	if ValidState("napping") {
		t.Error("napping should be invalid")
	}
}

func TestSaidDoneLeavesResolveAlone(t *testing.T) {
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	r := Record{StartedAt: at, Status: Status{Hook: &Report{State: StateIdle, At: at}}}
	if r.Status.SaidDone() {
		t.Fatal("a session that never said done says done")
	}
	r.Status.Done = &Report{State: StateDone, At: at.Add(time.Hour)}
	if !r.Status.SaidDone() {
		t.Fatal("the done mark is not read")
	}
	// The mark does not outrank the live sources, nor the exit.
	if got, src := Resolve(r); got.State != StateIdle || src != SourceHook {
		t.Errorf("Resolve with a done mark = %s from %s, want idle from hook", got.State, src)
	}
	code := 0
	r.Status.Exit = &Report{State: StateExited, Code: &code, At: at.Add(2 * time.Hour)}
	if got, src := Resolve(r); got.State != StateExited || src != SourceExit {
		t.Errorf("Resolve after the exit = %s from %s, want exited from exit", got.State, src)
	}
	if !r.Status.SaidDone() {
		t.Error("the exit cleared the done mark")
	}
}

func TestStalled(t *testing.T) {
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	rec := func(s Status) Record { return Record{StartedAt: at, Status: s} }
	hook := func(state, detail string) Status {
		return Status{Hook: &Report{State: state, Detail: detail, At: at}}
	}
	code := 0
	cases := []struct {
		name string
		r    Record
		in   time.Duration
		want bool
	}{
		{"working at 19m", rec(hook(StateWorking, "")), 19 * time.Minute, false},
		{"working at 20m", rec(hook(StateWorking, "")), 20 * time.Minute, false},
		{"working at 21m", rec(hook(StateWorking, "")), 21 * time.Minute, true},
		{"blocked", rec(hook(StateBlocked, "")), 3 * time.Hour, false},
		{"idle", rec(hook(StateIdle, "")), 3 * time.Hour, false},
		{"usage paused", rec(hook(StateWorking, LimitPausedDetail+" until 2026-01-01T05:00:00Z")), 3 * time.Hour, false},
		{"exited", rec(Status{Hook: &Report{State: StateWorking, At: at}, Exit: &Report{State: StateExited, Code: &code, At: at}}), 3 * time.Hour, false},
		{"no report", rec(Status{}), 3 * time.Hour, false},
	}
	for _, c := range cases {
		if got := Stalled(c.r, at.Add(c.in)); got != c.want {
			t.Errorf("%s: Stalled = %v, want %v", c.name, got, c.want)
		}
	}
}
