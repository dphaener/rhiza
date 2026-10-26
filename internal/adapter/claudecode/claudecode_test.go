package claudecode

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/adapter"
)

var hook = []string{"rhiza", "session", "hook"}

func launch(mode string) adapter.Launch {
	return adapter.Launch{SessionID: "sid-1", PermissionMode: mode, HookCommand: hook}
}

func TestRegisteredAsDefault(t *testing.T) {
	if a, ok := adapter.Lookup(Name); !ok || a.Name() != "claude-code" {
		t.Fatalf("Lookup = %v, %v", a, ok)
	}
	if adapter.Default().Name() != Name {
		t.Fatalf("Default = %v", adapter.Default().Name())
	}
}

func TestCommandForEachMode(t *testing.T) {
	a := Adapter{}
	for _, m := range a.PermissionModes() {
		argv, err := a.Command(launch(m))
		if err != nil {
			t.Fatalf("%s: %v", m, err)
		}
		i := slices.Index(argv, "--session-id")
		j := slices.Index(argv, "--permission-mode")
		if i < 0 || argv[i+1] != "sid-1" || j < 0 || argv[j+1] != m {
			t.Errorf("%s: argv = %q", m, argv)
		}
		t.Logf("argv: %q", argv[:len(argv)-1])
	}
}

func TestCommandManual(t *testing.T) {
	argv, err := Adapter{}.Command(launch("manual"))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("argv has --session-id sid-1 --permission-mode manual: %q", argv[:6])
	if got := strings.Join(argv[1:5], " "); got != "--session-id sid-1 --permission-mode manual" {
		t.Fatalf("argv = %q", argv)
	}
}

func TestCommandRefusesBadMode(t *testing.T) {
	for _, m := range []string{"", "yolo"} {
		argv, err := Adapter{}.Command(launch(m))
		if err == nil || argv != nil {
			t.Errorf("mode %q: want error and no argv, got %q, %v", m, argv, err)
		}
		t.Logf("mode %q: argv=%v err=%v", m, argv, err)
	}
}

func TestBypass(t *testing.T) {
	a := Adapter{}
	for _, m := range a.PermissionModes() {
		if a.Bypass(m) != (m == "bypassPermissions") {
			t.Errorf("Bypass(%q) = %v", m, a.Bypass(m))
		}
	}
}

func TestSettingsJSON(t *testing.T) {
	argv, err := Adapter{}.Command(launch("plan"))
	if err != nil {
		t.Fatal(err)
	}
	raw := argv[len(argv)-1]
	var got map[string]json.RawMessage
	if err := json.Unmarshal([]byte(raw), &got); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got["hooks"] == nil {
		t.Fatalf("settings keys = %v, want only hooks", got)
	}
	for _, bad := range []string{"permissions", "defaultMode", "permissionMode"} {
		if strings.Contains(raw, bad) || strings.Contains(raw, "plan") {
			t.Errorf("settings mention %q: %s", bad, raw)
		}
	}
	var hooks map[string][]hookGroup
	if err := json.Unmarshal(got["hooks"], &hooks); err != nil {
		t.Fatal(err)
	}
	for _, ev := range hookEvents {
		g := hooks[ev]
		if len(g) != 1 || len(g[0].Hooks) != 1 || g[0].Hooks[0].Type != "command" ||
			g[0].Hooks[0].Command != "rhiza session hook "+ev {
			t.Errorf("%s: %+v", ev, g)
		}
	}
	if len(hooks) != len(hookEvents) {
		t.Errorf("hooks = %d events", len(hooks))
	}
}

func TestSettingsQuotesHookCommand(t *testing.T) {
	s, err := Settings([]string{"/path with space/rhiza", "session", "hook"})
	if err != nil || !strings.Contains(s, `'/path with space/rhiza' session hook Stop`) {
		t.Fatalf("%s %v", s, err)
	}
	if _, err := Settings(nil); err == nil {
		t.Fatal("empty hook command should error")
	}
}

func TestHookStatus(t *testing.T) {
	cases := []struct {
		event, payload string
		state, kind    string
		ok             bool
	}{
		{"SessionStart", `{"session_id":"s"}`, "idle", "", true},
		{"UserPromptSubmit", `{"session_id":"s"}`, "working", "", true},
		{"PreToolUse", `{"session_id":"s"}`, "working", "", true},
		{"Notification", `{"session_id":"s","notification_type":"permission_prompt"}`, "blocked", "permission", true},
		{"Notification", `{"session_id":"s","notification_type":"idle_prompt"}`, "idle", "", true},
		{"Notification", `{"session_id":"s","notification_type":"auth_success"}`, "", "", false},
		{"Notification", `{"session_id":"s"}`, "", "", false},
		{"Stop", `{"session_id":"s"}`, "idle", "", true},
		{"SessionEnd", `{"session_id":"s"}`, "done", "", true},
	}
	for _, c := range cases {
		r, ok, err := Adapter{}.HookStatus(c.event, []byte(c.payload))
		if err != nil || ok != c.ok || r.State != c.state || r.Kind != c.kind || r.SessionID != "s" {
			t.Errorf("%s %s: %+v ok=%v err=%v", c.event, c.payload, r, ok, err)
		}
	}
	if _, _, err := (Adapter{}).HookStatus("Stop", []byte("{")); err == nil {
		t.Error("bad JSON should error")
	}
	if _, _, err := (Adapter{}).HookStatus("Bogus", nil); err == nil {
		t.Error("unknown event should error")
	}
}
