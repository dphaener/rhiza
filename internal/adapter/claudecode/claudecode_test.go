package claudecode

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
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

func TestCommandPrompt(t *testing.T) {
	l := launch("manual")
	argv, err := Adapter{}.Command(l)
	if err != nil {
		t.Fatal(err)
	}
	if argv[len(argv)-2] != "--settings" {
		t.Errorf("no prompt: argv ends %q, want the settings last", argv[len(argv)-2:])
	}
	l.Prompt = "/shipit:go 01ABC"
	with, err := Adapter{}.Command(l)
	if err != nil {
		t.Fatal(err)
	}
	if len(with) != len(argv)+1 || with[len(with)-1] != l.Prompt {
		t.Errorf("prompt: argv = %q, want the prompt as the last, positional argument", with)
	}
}

// A launch naming a model passes --model before the prompt; one naming
// none passes no --model.
func TestCommandModel(t *testing.T) {
	l := launch("manual")
	argv, err := Adapter{}.Command(l)
	if err != nil {
		t.Fatal(err)
	}
	if slices.Contains(argv, "--model") {
		t.Errorf("no model: argv = %q, want no --model", argv)
	}
	l.Model, l.Prompt = "opus", "go"
	with, err := Adapter{}.Command(l)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("argv ends %q", with[len(with)-3:])
	if i := slices.Index(with, "--model"); i < 0 || with[i+1] != "opus" || i+2 != len(with)-1 || with[len(with)-1] != "go" {
		t.Errorf("model: argv = %q, want --model opus before the prompt", with)
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
		{"PostToolUse", `{"session_id":"s"}`, "working", "", true},
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

func TestHookStatusAskUserQuestion(t *testing.T) {
	cases := []struct {
		name, event, payload string
		state, kind, detail  string
	}{
		{"one question", "PreToolUse", `{"session_id":"s","tool_name":"AskUserQuestion","tool_input":{"questions":[{"question":"Which database?"}]}}`, "blocked", "question", "Which database?"},
		{"two questions", "PreToolUse", `{"session_id":"s","tool_name":"AskUserQuestion","tool_input":{"questions":[{"question":"Which database?"},{"question":" Which cache? "}]}}`, "blocked", "question", "Which database? / Which cache?"},
		{"multi-line", "PreToolUse", `{"session_id":"s","tool_name":"AskUserQuestion","tool_input":{"questions":[{"question":"Line one\nline two\t\u001bend\n"}]}}`, "blocked", "question", "Line one line two  end"},
		{"no questions", "PreToolUse", `{"session_id":"s","tool_name":"AskUserQuestion","tool_input":{}}`, "blocked", "question", ""},
		{"other tool", "PreToolUse", `{"session_id":"s","tool_name":"Bash","tool_input":{"command":"ls"}}`, "working", "", ""},
		{"post tool use", "PostToolUse", `{"session_id":"s","tool_name":"AskUserQuestion","tool_input":{"questions":[{"question":"Which database?"}]}}`, "working", "", ""},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r, ok, err := Adapter{}.HookStatus(c.event, []byte(c.payload))
			if err != nil || !ok || r.State != c.state || r.Kind != c.kind || r.Detail != c.detail || r.SessionID != "s" {
				t.Errorf("%+v ok=%v err=%v", r, ok, err)
			}
		})
	}
}

func TestHookStatusMalformedPreToolUse(t *testing.T) {
	for _, raw := range []string{`{"tool_input": 5}`, `{"tool_name":"AskUserQuestion","tool_input":"x"}`, `{"tool_name":`} {
		r, ok, err := Adapter{}.HookStatus("PreToolUse", []byte(raw))
		if err == nil || ok || r != (adapter.HookReport{}) {
			t.Errorf("%s: want error and no report, got %+v ok=%v err=%v", raw, r, ok, err)
		}
	}
}

func TestSettingsInstallsPostToolUse(t *testing.T) {
	s, err := Settings([]string{"rhiza", "session", "hook"})
	if err != nil || !strings.Contains(s, `"PostToolUse":[{"hooks":[{"type":"command","command":"rhiza session hook PostToolUse"}]}]`) {
		t.Fatalf("%s %v", s, err)
	}
}

// sendHookPayload returns a PostToolUse payload for a SendMessage call with the
// given tool_input, shaped like the agent's hook stdin. The content is made up.
func sendHookPayload(event, tool, input string) string {
	return `{"session_id":"s","transcript_path":"/tmp/x.jsonl","cwd":"/tmp/w","hook_event_name":"` + event +
		`","tool_name":"` + tool + `","tool_input":` + input +
		`,"tool_response":{"success":true},"tool_use_id":"toolu_01ABC"}`
}

var sendInputs = []struct {
	name, input string
	want        adapter.Message
	err         bool
}{
	{"plain", `{"to":"reviewer","summary":"Ready for review","message":"Change set 2 is ready."}`,
		adapter.Message{SessionID: "s", To: "reviewer", Text: "Change set 2 is ready.", Summary: "Ready for review"}, false},
	{"escapes", `{"to":"reviewer","summary":"Two lines","message":"Line one\nsaid \"done\""}`,
		adapter.Message{SessionID: "s", To: "reviewer", Text: "Line one\nsaid \"done\"", Summary: "Two lines"}, false},
	{"structured", `{"to":"lead","summary":"Status","message":{"type":"status","ok":true}}`,
		adapter.Message{SessionID: "s", To: "lead", Text: `{"type":"status","ok":true}`, Summary: "Status", Structured: true}, false},
	{"empty text", `{"to":"reviewer","summary":"Nothing","message":""}`,
		adapter.Message{SessionID: "s", To: "reviewer", Text: "", Summary: "Nothing"}, false},
	{"no summary", `{"to":"reviewer","message":"Hello"}`,
		adapter.Message{SessionID: "s", To: "reviewer", Text: "Hello"}, false},
	{"no to", `{"summary":"x","message":"Hello"}`, adapter.Message{}, true},
	{"null to", `{"to":null,"message":"Hello"}`, adapter.Message{}, true},
	{"empty to", `{"to":"","message":"Hello"}`, adapter.Message{}, true},
	{"no message", `{"to":"reviewer","summary":"x"}`, adapter.Message{}, true},
	{"null message", `{"to":"reviewer","message":null}`, adapter.Message{}, true},
}

func TestHookMessage(t *testing.T) {
	for _, c := range sendInputs {
		t.Run(c.name, func(t *testing.T) {
			m, ok, err := Adapter{}.HookMessage("PostToolUse", []byte(sendHookPayload("PostToolUse", sendTool, c.input)))
			if c.err {
				if err == nil || ok {
					t.Fatalf("want an error, got %+v ok=%v err=%v", m, ok, err)
				}
				return
			}
			if err != nil || !ok || m != c.want {
				t.Fatalf("got %+v ok=%v err=%v, want %+v", m, ok, err, c.want)
			}
		})
	}
}

func TestHookMessageEscapesDecoded(t *testing.T) {
	m, _, err := Adapter{}.HookMessage("PostToolUse", []byte(sendHookPayload("PostToolUse", "SendMessage",
		`{"to":"reviewer","message":"a\nb \"c\""}`)))
	if err != nil || !strings.Contains(m.Text, "\n") || !strings.Contains(m.Text, `"c"`) || strings.Contains(m.Text, `\`) {
		t.Fatalf("%q %v", m.Text, err)
	}
}

func TestHookMessageNoMessage(t *testing.T) {
	input := `{"to":"reviewer","summary":"x","message":"Hello"}`
	cases := []struct{ name, event, tool string }{
		{"other tool", "PostToolUse", "Bash"},
		{"SessionStart", "SessionStart", "SendMessage"},
		{"UserPromptSubmit", "UserPromptSubmit", "SendMessage"},
		{"PreToolUse", "PreToolUse", "SendMessage"},
		{"Notification", "Notification", "SendMessage"},
		{"Stop", "Stop", "SendMessage"},
		{"SessionEnd", "SessionEnd", "SendMessage"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			m, ok, err := Adapter{}.HookMessage(c.event, []byte(sendHookPayload(c.event, c.tool, input)))
			if err != nil || ok || m.To != "" || m.Text != "" {
				t.Fatalf("got %+v ok=%v err=%v", m, ok, err)
			}
		})
	}
	for _, ev := range hookEvents {
		if ev == "PostToolUse" {
			continue
		}
		if _, ok, err := (Adapter{}).HookMessage(ev, []byte(sendHookPayload(ev, sendTool, input))); ok || err != nil {
			t.Errorf("%s: ok=%v err=%v", ev, ok, err)
		}
	}
	if _, ok, err := (Adapter{}).HookMessage("PostToolUse", []byte(sendHookPayload("PostToolUse", "Bash", `{"to":5,"message":null}`))); ok || err != nil {
		t.Errorf("other tool's input: ok=%v err=%v", ok, err)
	}
	if _, ok, err := (Adapter{}).HookMessage("PostToolUse", []byte("not json")); err == nil || ok {
		t.Errorf("bad JSON: ok=%v err=%v", ok, err)
	}
}

func TestHookStatusUnchangedBySendMessage(t *testing.T) {
	want := adapter.HookReport{SessionID: "s", State: "working"}
	for _, c := range sendInputs {
		t.Run(c.name, func(t *testing.T) {
			r, ok, err := Adapter{}.HookStatus("PostToolUse", []byte(sendHookPayload("PostToolUse", sendTool, c.input)))
			if err != nil || !ok || r != want {
				t.Fatalf("got %+v ok=%v err=%v, want %+v", r, ok, err, want)
			}
		})
	}
}

func TestSettingsOnePostToolUseGroupWithoutMatcher(t *testing.T) {
	s, err := Settings(hook)
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Hooks map[string][]map[string]json.RawMessage `json:"hooks"`
	}
	if err := json.Unmarshal([]byte(s), &got); err != nil {
		t.Fatal(err)
	}
	groups := got.Hooks["PostToolUse"]
	if len(groups) != 1 {
		t.Fatalf("PostToolUse groups = %d, want 1: %s", len(groups), s)
	}
	if _, has := groups[0]["matcher"]; has {
		t.Fatalf("PostToolUse group has a matcher: %s", s)
	}
}

type opt struct {
	Label       string `json:"label"`
	Description string `json:"description"`
}

type q struct {
	Question    string `json:"question"`
	Header      string `json:"header"`
	Options     []opt  `json:"options"`
	MultiSelect bool   `json:"multiSelect"`
}

func opts(labels ...string) []opt {
	out := make([]opt, len(labels))
	for i, l := range labels {
		out[i] = opt{Label: l, Description: "about " + l}
	}
	return out
}

// answeredLine returns one transcript line for an answered question call,
// shaped like the agent writes it.
func answeredLine(t *testing.T, id string, qs []q, answers map[string]string, notes map[string]string) string {
	t.Helper()
	ann := map[string]map[string]string{}
	for k, v := range notes {
		ann[k] = map[string]string{"notes": v}
	}
	line := map[string]any{
		"type": "user",
		"message": map[string]any{
			"role":    "user",
			"content": []any{map[string]any{"type": "tool_result", "tool_use_id": id, "content": "answered"}},
		},
		"toolUseResult": map[string]any{"questions": qs, "answers": answers, "annotations": ann},
	}
	b, err := json.Marshal(line)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func TestAnswers(t *testing.T) {
	color := q{Question: "Which colour?", Header: "Colour", Options: opts("Red", "Blue (Recommended)", "Green")}
	tools := q{Question: "Which tools?", Header: "Tools", Options: opts("Go", "Make", "Lint, then test", "Lint"), MultiSelect: true}
	big := strings.Repeat("x", 2<<20)
	rejected := `{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"toolu_rej","content":"The user doesn't want to proceed","is_error":true}]},"toolUseResult":"Error: The user doesn't want to proceed with this tool use."}`

	cases := []struct {
		name  string
		lines func(t *testing.T) []string
		want  []adapter.Answer
	}{
		{
			name: "single picked",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_1", []q{color}, map[string]string{color.Question: "Red"}, nil)}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_1", Question: "Which colour?", Header: "Colour", Options: []string{"Red", "Blue (Recommended)", "Green"}, Picked: []string{"Red"}}},
		},
		{
			name: "picked recommended label",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_1", []q{color}, map[string]string{color.Question: "Blue (Recommended)"}, nil)}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_1", Question: "Which colour?", Header: "Colour", Options: []string{"Red", "Blue (Recommended)", "Green"}, Picked: []string{"Blue (Recommended)"}}},
		},
		{
			name: "typed other",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_1", []q{color}, map[string]string{color.Question: "Teal, but only on Fridays"}, nil)}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_1", Question: "Which colour?", Header: "Colour", Options: []string{"Red", "Blue (Recommended)", "Green"}, Typed: "Teal, but only on Fridays"}},
		},
		{
			name: "multiSelect all picked",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_2", []q{tools}, map[string]string{tools.Question: "Make, Go"}, nil)}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_2", Question: "Which tools?", Header: "Tools", Options: []string{"Go", "Make", "Lint, then test", "Lint"}, MultiSelect: true, Picked: []string{"Make", "Go"}}},
		},
		{
			name: "multiSelect label containing comma",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_2", []q{tools}, map[string]string{tools.Question: "Go, Lint, then test, Lint"}, nil)}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_2", Question: "Which tools?", Header: "Tools", Options: []string{"Go", "Make", "Lint, then test", "Lint"}, MultiSelect: true, Picked: []string{"Go", "Lint, then test", "Lint"}}},
		},
		{
			name: "multiSelect only the comma label",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_2", []q{tools}, map[string]string{tools.Question: "Lint, then test"}, nil)}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_2", Question: "Which tools?", Header: "Tools", Options: []string{"Go", "Make", "Lint, then test", "Lint"}, MultiSelect: true, Picked: []string{"Lint, then test"}}},
		},
		{
			name: "multiSelect mixed picked and typed",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_2", []q{tools}, map[string]string{tools.Question: "Go, Bazel, sometimes, Make, Gopls"}, nil)}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_2", Question: "Which tools?", Header: "Tools", Options: []string{"Go", "Make", "Lint, then test", "Lint"}, MultiSelect: true, Picked: []string{"Go", "Make"}, Typed: "Bazel, sometimes, Gopls"}},
		},
		{
			name: "multiSelect label text inside typed text is not picked",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_2", []q{tools}, map[string]string{tools.Question: "Golang tools, Make"}, nil)}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_2", Question: "Which tools?", Header: "Tools", Options: []string{"Go", "Make", "Lint, then test", "Lint"}, MultiSelect: true, Picked: []string{"Make"}, Typed: "Golang tools"}},
		},
		{
			name: "notes",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_1", []q{color}, map[string]string{color.Question: "Green"}, map[string]string{color.Question: "Green matches the logo."})}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_1", Question: "Which colour?", Header: "Colour", Options: []string{"Red", "Blue (Recommended)", "Green"}, Picked: []string{"Green"}, Notes: "Green matches the logo."}},
		},
		{
			name:  "rejected question",
			lines: func(t *testing.T) []string { return []string{rejected} },
			want:  nil,
		},
		{
			name: "unanswered question in a call is skipped",
			lines: func(t *testing.T) []string {
				return []string{answeredLine(t, "toolu_3", []q{color, tools}, map[string]string{tools.Question: "Go"}, nil)}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_3", Question: "Which tools?", Header: "Tools", Options: []string{"Go", "Make", "Lint, then test", "Lint"}, MultiSelect: true, Picked: []string{"Go"}}},
		},
		{
			name: "two tool calls in order",
			lines: func(t *testing.T) []string {
				return []string{
					answeredLine(t, "toolu_a", []q{tools, color}, map[string]string{color.Question: "Red", tools.Question: "Make"}, nil),
					`{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Thanks."}]}}`,
					answeredLine(t, "toolu_b", []q{color}, map[string]string{color.Question: "Green"}, nil),
				}
			},
			want: []adapter.Answer{
				{ToolUseID: "toolu_a", Question: "Which tools?", Header: "Tools", Options: []string{"Go", "Make", "Lint, then test", "Lint"}, MultiSelect: true, Picked: []string{"Make"}},
				{ToolUseID: "toolu_a", Question: "Which colour?", Header: "Colour", Options: []string{"Red", "Blue (Recommended)", "Green"}, Picked: []string{"Red"}},
				{ToolUseID: "toolu_b", Question: "Which colour?", Header: "Colour", Options: []string{"Red", "Blue (Recommended)", "Green"}, Picked: []string{"Green"}},
			},
		},
		{
			name: "malformed line among good ones",
			lines: func(t *testing.T) []string {
				return []string{
					answeredLine(t, "toolu_a", []q{color}, map[string]string{color.Question: "Red"}, nil),
					`{"type":"user","toolUseResult":{"answers":`,
					`not json at all`,
					`["an", "array"]`,
					`{"type":"user","toolUseResult":{"questions":"wrong","answers":{"x":"y"}}}`,
					``,
					answeredLine(t, "toolu_b", []q{color}, map[string]string{color.Question: "Green"}, nil),
				}
			},
			want: []adapter.Answer{
				{ToolUseID: "toolu_a", Question: "Which colour?", Header: "Colour", Options: []string{"Red", "Blue (Recommended)", "Green"}, Picked: []string{"Red"}},
				{ToolUseID: "toolu_b", Question: "Which colour?", Header: "Colour", Options: []string{"Red", "Blue (Recommended)", "Green"}, Picked: []string{"Green"}},
			},
		},
		{
			name: "line over 1 MB",
			lines: func(t *testing.T) []string {
				return []string{
					`{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"` + big + `"}]}}`,
					answeredLine(t, "toolu_big", []q{color}, map[string]string{color.Question: big}, nil),
				}
			},
			want: []adapter.Answer{{ToolUseID: "toolu_big", Question: "Which colour?", Header: "Colour", Options: []string{"Red", "Blue (Recommended)", "Green"}, Typed: big}},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// No trailing newline, so the last line is read at EOF.
			got, err := Adapter{}.Answers([]byte(strings.Join(c.lines(t), "\n")))
			if err != nil {
				t.Fatalf("Answers: %v", err)
			}
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("Answers =\n%s\nwant\n%s", brief(got), brief(c.want))
			}
			t.Logf("answers: %s", brief(got))
		})
	}
}

func TestAnswersEmpty(t *testing.T) {
	for _, in := range []string{"", "\n\n", "{}\n"} {
		got, err := Adapter{}.Answers([]byte(in))
		if err != nil || got != nil {
			t.Errorf("Answers(%q) = %v, %v", in, got, err)
		}
	}
}

// brief renders answers for logs, cutting long strings.
func brief(as []adapter.Answer) string {
	cut := func(s string) string {
		if len(s) > 60 {
			return fmt.Sprintf("%s...(%d bytes)", s[:20], len(s))
		}
		return s
	}
	var b strings.Builder
	for _, a := range as {
		fmt.Fprintf(&b, "{id=%s q=%q picked=%q typed=%q notes=%q multi=%v} ", a.ToolUseID, cut(a.Question), a.Picked, cut(a.Typed), a.Notes, a.MultiSelect)
	}
	if b.Len() == 0 {
		return "none"
	}
	return b.String()
}

func TestTranscriptGlobs(t *testing.T) {
	env := func(vars map[string]string) func(string) string {
		return func(k string) string { return vars[k] }
	}
	home := map[string]string{"HOME": "/home/o"}
	both := map[string]string{"HOME": "/home/o", "CLAUDE_CONFIG_DIR": "/cfg/claude"}
	cases := []struct {
		name string
		sid  string
		env  map[string]string
		want []string
	}{
		{"home only", "sid-1", home, []string{"/home/o/.claude/projects/*/sid-1.jsonl"}},
		{"config dir first", "sid-1", both, []string{"/cfg/claude/projects/*/sid-1.jsonl", "/home/o/.claude/projects/*/sid-1.jsonl"}},
		{"config dir pattern syntax escaped", "sid-1", map[string]string{"CLAUDE_CONFIG_DIR": "/c[1]/*"}, []string{`/c\[1]/\*/projects/*/sid-1.jsonl`}},
		{"no environment", "sid-1", map[string]string{}, nil},
		{"hostile empty", "", both, nil},
		{"hostile slash", "../../etc/passwd", both, nil},
		{"hostile dotdot", "..", both, nil},
		{"hostile star", "*", both, nil},
		{"hostile question", "sid?", both, nil},
		{"hostile bracket", "sid[1]", both, nil},
		{"hostile backslash", `sid\*`, both, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Adapter{}.TranscriptGlobs(c.sid, env(c.env))
			if !reflect.DeepEqual(got, c.want) {
				t.Fatalf("TranscriptGlobs(%q) = %q, want %q", c.sid, got, c.want)
			}
			t.Logf("globs: %q", got)
		})
	}
}

func TestTranscriptGlobsMatchOnlyTheSessionsOwn(t *testing.T) {
	dir := t.TempDir()
	own := filepath.Join(dir, "projects", "-repo", "sid-1.jsonl")
	for _, p := range []string{
		own,
		filepath.Join(dir, "projects", "-repo", "sid-1", "subagents", "agent-x.jsonl"),
		filepath.Join(dir, "projects", "-repo", "sid-10.jsonl"),
	} {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	globs := Adapter{}.TranscriptGlobs("sid-1", func(k string) string {
		if k == "CLAUDE_CONFIG_DIR" {
			return dir
		}
		return ""
	})
	var got []string
	for _, g := range globs {
		m, err := filepath.Glob(g)
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, m...)
	}
	if !reflect.DeepEqual(got, []string{own}) {
		t.Fatalf("matches = %q, want %q", got, []string{own})
	}
}

func TestAskReply(t *testing.T) {
	ask := `{"session_id":"s","tool_name":"AskUserQuestion","tool_input":{"questions":[{"question":"Which database?"}]}}`
	reply, ok, err := Adapter{}.AskReply("PreToolUse", []byte(ask), `no "owner"`)
	if err != nil || !ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	want := `{"hookSpecificOutput":{"hookEventName":"PreToolUse","permissionDecision":"deny","permissionDecisionReason":"no \"owner\""}}`
	if string(reply) != want {
		t.Errorf("reply = %s\nwant    %s", reply, want)
	}
	for _, c := range []struct{ event, payload string }{
		{"PreToolUse", `{"session_id":"s","tool_name":"Bash"}`},
		{"PostToolUse", ask},
		{"Stop", `{"session_id":"s"}`},
		{"Notification", `{"notification_type":"permission_prompt"}`},
	} {
		reply, ok, err := Adapter{}.AskReply(c.event, []byte(c.payload), "r")
		if err != nil || ok || reply != nil {
			t.Errorf("%s %s: reply=%s ok=%v err=%v", c.event, c.payload, reply, ok, err)
		}
	}
	if _, ok, err := (Adapter{}).AskReply("PreToolUse", []byte("{"), "r"); err == nil || ok {
		t.Errorf("bad JSON: ok=%v err=%v", ok, err)
	}
}

// tokenTranscript is made-up content in the shape of a Claude Code
// transcript: msg_B is split over three lines that repeat its usage, and the
// helper's msg_H is logged twice on sidechain lines.
const tokenTranscript = `{"type":"user","isSidechain":false,"message":{"role":"user","content":"hi"}}
{"type":"assistant","isSidechain":false,"message":{"id":"msg_A","usage":{"input_tokens":3,"output_tokens":10,"cache_creation_input_tokens":100,"cache_read_input_tokens":1000}}}
{"type":"assistant","isSidechain":false,"message":{"id":"msg_B","usage":{"input_tokens":5,"output_tokens":20,"cache_creation_input_tokens":200,"cache_read_input_tokens":2000}}}
{"type":"assistant","isSidechain":false,"message":{"id":"msg_B","usage":{"input_tokens":5,"output_tokens":20,"cache_creation_input_tokens":200,"cache_read_input_tokens":2000}}}
{not json
{"type":"assistant","isSidechain":false,"message":{"id":"msg_B","usage":{"input_tokens":5,"output_tokens":20,"cache_creation_input_tokens":200,"cache_read_input_tokens":2000}}}
{"type":"assistant","isSidechain":true,"message":{"id":"msg_H","usage":{"input_tokens":11,"output_tokens":40,"cache_creation_input_tokens":400,"cache_read_input_tokens":4000}}}
{"type":"assistant","isSidechain":true,"message":{"id":"msg_H","usage":{"input_tokens":11,"output_tokens":40,"cache_creation_input_tokens":400,"cache_read_input_tokens":4000}}}
{"type":"user","isSidechain":false,"message":{"role":"user","content":"more","usage":{"input_tokens":999}}}
{"type":"assistant","isSidechain":false,"message":{"id":"msg_C","usage":{"input_tokens":7,"output_tokens":30,"cache_creation_input_tokens":300,"cache_read_input_tokens":3000}}}
{"type":"assistant","isSidechain":false,"message":{"usage":{"input_tokens":1,"output_tokens":2,"cache_creation_input_tokens":4,"cache_read_input_tokens":8}}}
{"type":"assistant","isSidechain":false,"message":{"id":"msg_D"}}
`

func TestTokenUse(t *testing.T) {
	got, err := Adapter{}.TokenUse([]byte(tokenTranscript))
	if err != nil {
		t.Fatal(err)
	}
	want := adapter.TokenUse{
		Session: adapter.TokenCounts{Input: 3 + 5 + 7 + 1, Output: 10 + 20 + 30 + 2,
			CacheCreation: 100 + 200 + 300 + 4, CacheRead: 1000 + 2000 + 3000 + 8},
		Helpers: adapter.TokenCounts{Input: 11, Output: 40, CacheCreation: 400, CacheRead: 4000},
	}
	if got != want {
		t.Fatalf("TokenUse = %+v, want %+v", got, want)
	}
}

func TestTokenUseEmpty(t *testing.T) {
	for _, in := range []string{"", "\n\n"} {
		got, err := Adapter{}.TokenUse([]byte(in))
		if err != nil || got != (adapter.TokenUse{}) {
			t.Fatalf("TokenUse(%q) = %+v, %v; want zeros, nil", in, got, err)
		}
	}
}

func TestTokenUseBroken(t *testing.T) {
	got, err := Adapter{}.TokenUse([]byte("{not json\n[1,2]\n\"usage\"\n{\"type\":\"assistant\",\"message\":{\"usage\":\"x\"}}"))
	if err != nil || got != (adapter.TokenUse{}) {
		t.Fatalf("TokenUse = %+v, %v; want zeros, nil", got, err)
	}
}

func TestTokenUseLongLine(t *testing.T) {
	pad := strings.Repeat("x", 70*1024)
	in := `{"type":"assistant","isSidechain":false,"pad":"` + pad + `","message":{"id":"msg_L","usage":{"input_tokens":2,"output_tokens":3,"cache_creation_input_tokens":5,"cache_read_input_tokens":7}}}` + "\n"
	got, err := Adapter{}.TokenUse([]byte(in))
	if err != nil {
		t.Fatal(err)
	}
	want := adapter.TokenUse{Session: adapter.TokenCounts{Input: 2, Output: 3, CacheCreation: 5, CacheRead: 7}}
	if got != want {
		t.Fatalf("TokenUse = %+v, want %+v", got, want)
	}
}

func TestCommandResume(t *testing.T) {
	argv, err := Adapter{}.Command(adapter.Launch{
		SessionID: "s-1", PermissionMode: "manual", HookCommand: []string{"rhiza", "session", "hook"},
		Resume: true, Prompt: "go on",
	})
	if err != nil {
		t.Fatal(err)
	}
	got := strings.Join(argv[:5], " ")
	if want := "claude --resume s-1 --permission-mode manual"; got != want {
		t.Errorf("argv = %q, want prefix %q", got, want)
	}
	if strings.Contains(strings.Join(argv, " "), "--session-id") {
		t.Errorf("a resumed session must not start a new ID: %v", argv)
	}
	if argv[len(argv)-1] != "go on" {
		t.Errorf("last arg = %q, want the prompt", argv[len(argv)-1])
	}
}

func TestUsageLimit(t *testing.T) {
	limited := `{"type":"user","message":{"content":"work"}}
{"type":"assistant","isApiErrorMessage":true,"message":{"model":"<synthetic>","content":[{"type":"text","text":"Claude AI usage limit reached|1700000000"}]}}
`
	l, ok := Adapter{}.UsageLimit([]byte(limited))
	if !ok || l.ResetAt.Unix() != 1700000000 {
		t.Errorf("UsageLimit = %v, %v; want reset 1700000000", l, ok)
	}
	l, ok = Adapter{}.UsageLimit([]byte(`{"type":"assistant","message":{"model":"<synthetic>","content":"You've hit your limit"}}`))
	if !ok || !l.ResetAt.IsZero() {
		t.Errorf("no reset given: got %v, %v", l, ok)
	}
	resumed := limited + `{"type":"assistant","message":{"content":[{"type":"text","text":"done"}]}}` + "\n"
	if _, ok := (Adapter{}).UsageLimit([]byte(resumed)); ok {
		t.Error("a session that went on past the limit is not limited")
	}
	if _, ok := (Adapter{}).UsageLimit(nil); ok {
		t.Error("an empty transcript is not limited")
	}
}

// The phrases in the model's own turn, a tool call's input or its prose,
// are not the limit: only the synthetic API-error message is.
func TestUsageLimitIgnoresModelTurns(t *testing.T) {
	for name, line := range map[string]string{
		"tool_use input": `{"type":"assistant","message":{"model":"claude-opus-4","content":[{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"grep -rn \"hit your limit\" internal/"}}]}}`,
		"model prose":    `{"type":"assistant","message":{"model":"claude-opus-4","content":[{"type":"text","text":"UsageLimit matches \"usage limit reached|1700000000\"."}]}}`,
		"error tool_use": `{"type":"assistant","isApiErrorMessage":true,"message":{"model":"<synthetic>","content":[{"type":"tool_use","id":"t1","name":"Bash","input":{"command":"echo usage limit reached"}}]}}`,
	} {
		if l, ok := (Adapter{}).UsageLimit([]byte(line + "\n")); ok {
			t.Errorf("%s: UsageLimit = %v, true; want not limited", name, l)
		}
	}
}

func TestPaneUsageLimit(t *testing.T) {
	a := Adapter{}
	if !a.PaneUsageLimit("work\nClaude usage limit reached|1700000000\n> ") {
		t.Fatal("limit prompt not seen")
	}
	if !a.PaneUsageLimit("You've hit your limit · resets 3pm\n") {
		t.Fatal("hit-your-limit not seen")
	}
	old := "hit your limit\n" + strings.Repeat("later work\n", 30)
	if a.PaneUsageLimit(old) || a.PaneUsageLimit("all fine") {
		t.Fatal("false positive")
	}
}

func TestTrust(t *testing.T) {
	a := Adapter{}
	out, err := a.Trust(nil, []string{"/r", "/r/wt", ""})
	if err != nil {
		t.Fatal(err)
	}
	var got struct {
		Projects map[string]map[string]any `json:"projects"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.Projects) != 2 || got.Projects["/r"]["hasTrustDialogAccepted"] != true || got.Projects["/r/wt"]["hasTrustDialogAccepted"] != true {
		t.Errorf("created %s", out)
	}
	again, err := a.Trust(out, []string{"/r", "/r/wt"})
	if err != nil || string(again) != string(out) {
		t.Errorf("not idempotent: %v\n%s\n%s", err, out, again)
	}
	if _, err := a.Trust([]byte("{nope"), []string{"/r"}); err == nil {
		t.Error("broken JSON was accepted")
	}
	in := `{"big": 12345678901234567890, "projects": {"/r": {"x": [1, 2], "hasTrustDialogAccepted": false}}}`
	out, err = a.Trust([]byte(in), []string{"/r"})
	if err != nil || !strings.Contains(string(out), "12345678901234567890") || !strings.Contains(string(out), `"x": [`) {
		t.Errorf("other keys lost: %v %s", err, out)
	}
}

func TestTrustFile(t *testing.T) {
	a := Adapter{}
	env := map[string]string{"HOME": "/h"}
	get := func(k string) string { return env[k] }
	if got := a.TrustFile(get); got != "/h/.claude.json" {
		t.Errorf("got %q", got)
	}
	env["CLAUDE_CONFIG_DIR"] = "/c"
	if got := a.TrustFile(get); got != "/c/.claude.json" {
		t.Errorf("got %q", got)
	}
	if a.TrustFile(nil) != "" || a.TrustFile(func(string) string { return "" }) != "" {
		t.Error("a place was invented")
	}
}

// Command only builds argv: it writes no file, trust included.
func TestCommandWritesNoFile(t *testing.T) {
	home, conf := t.TempDir(), t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLAUDE_CONFIG_DIR", conf)
	if _, err := (Adapter{}).Command(launch("manual")); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{home, conf} {
		if ents, _ := os.ReadDir(d); len(ents) != 0 {
			t.Errorf("Command wrote %v in %s", ents, d)
		}
	}
}
