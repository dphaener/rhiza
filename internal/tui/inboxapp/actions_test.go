package inboxapp

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

const actionItems = `{"items":[
{"id":"q1","kind":"question","run":"RUNQ00000001","change":"C1","step":"build","session":"SESS1","question":"Which db?","recommendation":"sqlite"},
{"id":"s2","kind":"signoff","run":"RUNS00000002","change":"C2","step":"signoff"},
{"id":"f3","kind":"failed","run":"RUNF00000003","change":"C3","step":"verify"},
{"id":"g4","kind":"gave-up","run":"RUNG00000004","change":"C4","step":"review"}
]}`

type call struct {
	op   string
	args []string
}

type rig struct {
	m     *Model
	calls []call
	execs [][]string
	// changes is the findings file's content at the time of the call.
	changes string
}

func newRig(t *testing.T) *rig {
	t.Helper()
	r := &rig{}
	m := NewModel(func(op string, args []string) ([]byte, error) {
		if op == Name {
			return []byte(actionItems), nil
		}
		r.calls = append(r.calls, call{op, args})
		if op == Name+".request-changes" {
			b, err := os.ReadFile(args[1])
			if err != nil {
				t.Fatal(err)
			}
			r.changes = string(b)
		}
		return []byte(`{}`), nil
	})
	m.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		r.execs = append(r.execs, append([]string{name}, args...))
		return nil, nil
	}
	m.Self, m.Root = "/bin/rhiza", "/repo"
	r.m = m
	return r
}

func (r *rig) press(keys ...string) {
	for _, k := range keys {
		switch k {
		case "enter":
			r.m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
		case "esc":
			r.m.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
		default:
			r.m.Update(tea.KeyPressMsg{Code: []rune(k)[0], Text: k})
		}
	}
}

func (r *rig) typeText(s string) {
	for _, c := range s {
		r.press(string(c))
	}
}

func (r *rig) at(n int) {
	for range n {
		r.press("j")
	}
}

func TestOpenOnEnterAndO(t *testing.T) {
	r := newRig(t)
	r.press("o", "enter")
	if len(r.calls) != 2 || r.calls[0].op != "inbox.open" || r.calls[0].args[0] != "q1" || r.calls[1].op != "inbox.open" {
		t.Fatalf("calls: %v", r.calls)
	}
}

func TestAnswerSendsTypedText(t *testing.T) {
	r := newRig(t)
	r.press("a")
	r.typeText("use pg")
	r.press("enter")
	want := call{"session.message", []string{"SESS1", "use pg"}}
	if len(r.calls) != 1 || r.calls[0].op != want.op || strings.Join(r.calls[0].args, "|") != "SESS1|use pg" {
		t.Fatalf("calls: %v", r.calls)
	}
}

func TestTakeRecommendation(t *testing.T) {
	r := newRig(t)
	r.press("t")
	if len(r.calls) != 1 || r.calls[0].op != "session.message" || strings.Join(r.calls[0].args, "|") != "SESS1|sqlite" {
		t.Fatalf("calls: %v", r.calls)
	}
}

func TestApproveAndRequestChanges(t *testing.T) {
	r := newRig(t)
	r.at(1)
	r.press("y")
	if len(r.calls) != 1 || r.calls[0].op != "inbox.approve" || r.calls[0].args[0] != "s2" {
		t.Fatalf("calls: %v", r.calls)
	}
	r.calls = nil
	r.press("c")
	r.typeText("add tests")
	r.press("enter")
	if len(r.calls) != 1 || r.calls[0].op != "inbox.request-changes" || r.calls[0].args[0] != "s2" {
		t.Fatalf("calls: %v", r.calls)
	}
	var f struct {
		Findings []map[string]string `json:"findings"`
	}
	if err := json.Unmarshal([]byte(r.changes), &f); err != nil || len(f.Findings) != 1 || f.Findings[0]["reason"] != "add tests" {
		t.Fatalf("findings file %q: %v", r.changes, err)
	}
	if _, err := os.Stat(r.calls[0].args[1]); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("findings file left behind: %v", err)
	}
}

func TestRetryOpensTab(t *testing.T) {
	r := newRig(t)
	r.at(2)
	r.press("R")
	if len(r.execs) != 1 {
		t.Fatalf("execs: %v", r.execs)
	}
	got := strings.Join(r.execs[0], " ")
	if want := "rex window new RUNF00000 -c /repo --json -- /bin/rhiza run resume RUNF00000003 --retry"; got != want {
		t.Fatalf("exec %q, want %q", got, want)
	}
	r.at(1)
	r.press("R")
	if got := strings.Join(r.execs[1], " "); !strings.HasSuffix(got, "run resume RUNG00000004") {
		t.Fatalf("exec %q", got)
	}
}

func TestAbandonNeedsY(t *testing.T) {
	r := newRig(t)
	r.at(2)
	r.press("x", "n")
	if len(r.calls) != 0 {
		t.Fatalf("n abandoned: %v", r.calls)
	}
	r.press("x", "y")
	if len(r.calls) != 1 || r.calls[0].op != "run.abandon" || r.calls[0].args[0] != "RUNF00000003" {
		t.Fatalf("calls: %v", r.calls)
	}
}

func TestActionsOnlyOnTheirKinds(t *testing.T) {
	r := newRig(t)
	r.press("y", "c", "R", "x") // on a question
	r.at(2)
	r.press("a", "t", "y", "c") // on a failed run
	r.press("enter")
	if len(r.execs) != 0 || len(r.calls) != 1 || r.calls[0].op != "inbox.open" {
		t.Fatalf("calls %v execs %v", r.calls, r.execs)
	}
}

func TestEscCancelsTyping(t *testing.T) {
	r := newRig(t)
	r.press("a")
	r.typeText("x")
	r.press("esc")
	r.press("enter")
	if len(r.calls) != 1 || r.calls[0].op != "inbox.open" {
		t.Fatalf("calls: %v", r.calls)
	}
}
