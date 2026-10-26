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
