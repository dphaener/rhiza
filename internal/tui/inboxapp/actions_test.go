package inboxapp

import (
	"encoding/json"
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
	// onCall sees each call but the inbox's own.
	onCall func(call)
	// items is what the inbox operation returns.
	items string
}

func newRig(t *testing.T) *rig {
	t.Helper()
	r := &rig{items: actionItems}
	m := NewModel(func(op string, args []string) ([]byte, error) {
		if op == Name {
			return []byte(r.items), nil
		}
		c := call{op, args}
		r.calls = append(r.calls, c)
		if r.onCall != nil {
			r.onCall(c)
		}
		return []byte(`{}`), nil
	})
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

func TestApproveAndRequestChangesCallTheOperations(t *testing.T) {
	r := newRig(t)
	r.at(1)
	r.press("y")
	if len(r.calls) != 1 || r.calls[0].op != "run.approve" || strings.Join(r.calls[0].args, "|") != "RUNS00000002" {
		t.Fatalf("calls %v", r.calls)
	}
	// run.request-changes reads the findings from the file it is given.
	var f struct {
		Findings []map[string]string `json:"findings"`
	}
	var file string
	r.onCall = func(c call) {
		file = c.args[0]
		b, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(b, &f); err != nil {
			t.Fatal(err)
		}
	}
	r.press("c")
	r.typeText("add tests")
	r.press("enter")
	if len(r.calls) != 2 || r.calls[1].op != "run.request-changes" || r.calls[1].args[1] != "RUNS00000002" {
		t.Fatalf("calls %v", r.calls)
	}
	if len(f.Findings) != 1 || f.Findings[0]["reason"] != "add tests" || f.Findings[0]["check"] != "quality" {
		t.Fatalf("findings %v", f)
	}
	if _, err := os.Stat(file); err == nil {
		t.Errorf("findings file %s left behind", file)
	}
}

func TestModeActsOnItsItem(t *testing.T) {
	r := newRig(t)
	r.at(2)
	r.press("x")
	// f3 stops waiting: the cursor moves to g4, and the abandon is cancelled.
	r.items = strings.Replace(actionItems, `{"id":"f3","kind":"failed","run":"RUNF00000003","change":"C3","step":"verify"},`, "", 1)
	r.m.Update(TickMsg{})
	r.press("y")
	if len(r.calls) != 0 {
		t.Fatalf("calls %v", r.calls)
	}
	if !strings.Contains(r.m.View().Content, "no longer waits") {
		t.Fatalf("view:\n%s", r.m.View().Content)
	}
	// The cursor moves while the answer is typed: it still goes to q1.
	r.items = actionItems
	r.m.cursor = 0
	r.press("a")
	r.m.cursor = 1
	r.typeText("pg")
	r.press("enter")
	if len(r.calls) != 1 || strings.Join(r.calls[0].args, "|") != "SESS1|pg" {
		t.Fatalf("calls: %v", r.calls)
	}
}

func TestRetryCallsResume(t *testing.T) {
	r := newRig(t)
	r.at(2)
	r.press("R")
	if len(r.calls) != 1 || r.calls[0].op != "run.resume" || strings.Join(r.calls[0].args, " ") != "RUNF00000003 --retry" {
		t.Fatalf("calls: %v", r.calls)
	}
	r.at(1)
	r.press("R")
	if len(r.calls) != 2 || r.calls[1].op != "run.resume" || strings.Join(r.calls[1].args, " ") != "RUNG00000004" {
		t.Fatalf("calls: %v", r.calls)
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
	if len(r.calls) != 1 || r.calls[0].op != "inbox.open" {
		t.Fatalf("calls %v", r.calls)
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

const blockerItems = `{"items":[
{"id":"b1","kind":"blocker","change":"C9","change_title":"Ship it","question":"Which db?","handle":"which-db","blocks":"C9"},
{"id":"q1","kind":"question","run":"RUNQ00000001","change":"C1","step":"build","session":"SESS1","question":"Which db?"}
]}`

func TestAnswerBlockerCallsQuestionAnswer(t *testing.T) {
	r := newRig(t)
	r.items = blockerItems
	r.m.Reload()
	r.at(1)
	r.press("a")
	r.typeText("use sqlite")
	r.press("enter")
	if len(r.calls) != 1 || r.calls[0].op != "question.answer" ||
		strings.Join(r.calls[0].args, "|") != "which-db|C9|use sqlite" {
		t.Fatalf("calls: %v", r.calls)
	}
}

func TestAnswerRunQuestionStillMessagesSession(t *testing.T) {
	r := newRig(t)
	r.items = blockerItems
	r.m.Reload()
	r.press("a")
	r.typeText("sqlite")
	r.press("enter")
	if len(r.calls) != 1 || r.calls[0].op != "session.message" ||
		strings.Join(r.calls[0].args, "|") != "SESS1|sqlite" {
		t.Fatalf("calls: %v", r.calls)
	}
}

func TestEscOnBlockerAnswerSendsNothing(t *testing.T) {
	r := newRig(t)
	r.items = blockerItems
	r.m.Reload()
	r.at(1)
	r.press("a")
	r.typeText("use sqlite")
	r.press("esc")
	if len(r.calls) != 0 {
		t.Fatalf("calls: %v", r.calls)
	}
}

const nudgeItems = `{"items":[
{"id":"n1","kind":"stalled","run":"RUNN00000001","change":"C1","step":"build","session":"SESSN1"},
{"id":"f2","kind":"failed","run":"RUNF00000002","change":"C2","step":"verify"}
]}`

func newNudgeRig(t *testing.T, pane string) *rig {
	t.Helper()
	r := &rig{}
	m := NewModel(func(op string, args []string) ([]byte, error) {
		switch op {
		case Name:
			return []byte(nudgeItems), nil
		case "session.pane":
			r.calls = append(r.calls, call{op, args})
			return []byte(pane), nil
		}
		c := call{op, args}
		r.calls = append(r.calls, c)
		if r.onCall != nil {
			r.onCall(c)
		}
		return []byte(`{}`), nil
	})
	r.m = m
	return r
}

func opsOf(cs []call) string {
	var ops []string
	for _, c := range cs {
		ops = append(ops, c.op+":"+strings.Join(c.args[1:], ","))
	}
	return strings.Join(ops, " ")
}

func TestNudgeSendsCarryOn(t *testing.T) {
	r := newNudgeRig(t, `{"session":"SESSN1","text":"working","usage_limit":false}`)
	r.at(1) // failed runs list before stalled sessions
	r.press("n")
	if got, want := opsOf(r.calls), "session.pane: session.message:"+carryOn; got != want {
		t.Fatalf("calls %q, want %q", got, want)
	}
	if r.calls[1].args[0] != "SESSN1" {
		t.Fatalf("calls: %v", r.calls)
	}
}

func TestNudgeEscapesUsageLimitPromptFirst(t *testing.T) {
	r := newNudgeRig(t, `{"session":"SESSN1","text":"You've hit your limit","usage_limit":true}`)
	r.at(1)
	r.press("n")
	if got, want := opsOf(r.calls), "session.pane: session.key:escape session.message:"+carryOn; got != want {
		t.Fatalf("calls %q, want %q", got, want)
	}
}

func TestNudgeOnlyOnStalled(t *testing.T) {
	r := newNudgeRig(t, `{}`)
	r.press("n") // on the failed run
	if len(r.calls) != 0 {
		t.Fatalf("calls: %v", r.calls)
	}
}
