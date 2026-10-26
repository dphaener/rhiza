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
