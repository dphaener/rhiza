package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
)

// roadmapRepo is a temporary repository, the working directory for the
// test, with a proposed change set and an accepted one blocked by a
// question. rhiza runs a command in it and returns its stdout.
type roadmapRepo struct {
	t   *testing.T
	env map[string]string
}

func (r roadmapRepo) getenv(k string) string { return r.env[k] }

func (r roadmapRepo) rhiza(args ...string) string {
	r.t.Helper()
	var out, errb bytes.Buffer
	if code := run(args, nil, &out, &errb, r.getenv); code != 0 {
		r.t.Fatalf("rhiza %s: exit %d: %s%s", strings.Join(args, " "), code, out.String(), errb.String())
	}
	return out.String()
}

// changeID makes a change set of intent and returns its ID.
func (r roadmapRepo) changeID(intent string) string {
	r.t.Helper()
	var c struct {
		ChangeSet string `json:"change_set"`
		ID        string `json:"id"`
	}
	if err := json.Unmarshal([]byte(r.rhiza("change", "new", intent, "--json")), &c); err != nil {
		r.t.Fatal(err)
	}
	return cmpOr(c.ChangeSet, c.ID)
}

func cmpOr(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func newRoadmapRepo(t *testing.T) (r roadmapRepo, proposed, blocked string) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	if err := os.MkdirAll(home, 0o755); err != nil {
		t.Fatal(err)
	}
	gitconfig := filepath.Join(home, ".gitconfig")
	if err := os.WriteFile(gitconfig, []byte(testGitConfig), 0o644); err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(dir, "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Chdir(repo)
	r = roadmapRepo{t: t, env: map[string]string{
		"HOME": home, "GIT_CONFIG_NOSYSTEM": "1", "GIT_CONFIG_GLOBAL": gitconfig,
		"RHIZA_TEST_NOW": testNow, "RHIZA_TEST_SEED": testSeed,
	}}
	for k, v := range r.env {
		t.Setenv(k, v)
	}
	r.rhiza("init")
	r.rhiza("add", "intent", "grow", "--record")
	r.rhiza("add", "capability", "store", "--record")

	proposed = r.changeID("Proposed change")
	r.rhiza("add", "capability", "draft")
	r.rhiza("propose")

	blocked = r.changeID("Blocked change")
	r.rhiza("add", "capability", "report")
	r.rhiza("add", "question", "which-format", "--field", "statement=Which format do reports use?")
	r.rhiza("link", "which-format", "blocks", "report")
	r.rhiza("propose")
	r.rhiza("accept", "--change", blocked)
	return r, proposed, blocked
}

func send(m tea.Model, keys ...tea.KeyPressMsg) {
	for _, k := range keys {
		m.Update(k)
	}
}

func typed(s string) []tea.KeyPressMsg {
	var ks []tea.KeyPressMsg
	for _, r := range s {
		ks = append(ks, tea.KeyPressMsg{Code: r, Text: string(r)})
	}
	return ks
}

func (r roadmapRepo) state(cs string) string {
	r.t.Helper()
	var c api.ChangeShowResponse
	if err := json.Unmarshal([]byte(r.rhiza("change", "show", cs, "--json")), &c); err != nil {
		r.t.Fatal(err)
	}
	return string(c.State)
}

func TestRoadmapAppAcceptsAProposedChangeSet(t *testing.T) {
	r, proposed, _ := newRoadmapRepo(t)
	m := roadmapModel(globals{}, callOver(r))
	for i := 0; m.Cursor().String() != proposed; i++ {
		if i > 5 {
			t.Fatal("proposed change set not in the list")
		}
		send(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}
	send(m, tea.KeyPressMsg{Code: 'a', Text: "a"})
	v := m.View().Content
	if !strings.Contains(v, "What-if") || !strings.Contains(v, "proposed -> accepted") || !strings.Contains(v, "y/n") {
		t.Fatalf("no what-if:\n%s", v)
	}
	if got := r.state(proposed); got != "proposed" {
		t.Fatalf("state before y = %s", got)
	}
	send(m, tea.KeyPressMsg{Code: 'y', Text: "y"})
	if got := r.state(proposed); got != "accepted" {
		t.Fatalf("state after y = %s", got)
	}
	if v := m.View().Content; !strings.Contains(v, "[accepted] Proposed change") {
		t.Errorf("list not reloaded:\n%s", v)
	}
}

func TestRoadmapAppAnswersABlockingQuestion(t *testing.T) {
	r, _, blocked := newRoadmapRepo(t)
	m := roadmapModel(globals{}, callOver(r))
	for i := 0; m.Cursor().String() != blocked; i++ {
		if i > 5 {
			t.Fatal("blocked change set not in the list")
		}
		send(m, tea.KeyPressMsg{Code: 'j', Text: "j"})
	}
	send(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	if v := m.View().Content; !strings.Contains(v, "Which format do reports use?") {
		t.Fatalf("detail lacks the question:\n%s", v)
	}
	if out := r.rhiza("ready", "--json"); !strings.Contains(out, "blocked by question") {
		t.Fatalf("not blocked before the answer: %s", out)
	}
	send(m, tea.KeyPressMsg{Code: tea.KeyEnter})
	send(m, typed("Use CSV")...)
	send(m, tea.KeyPressMsg{Code: tea.KeyEnter})

	if out := r.rhiza("show", "answer-which-format", "--change", blocked, "--json"); !strings.Contains(out, `Use CSV`) || !strings.Contains(out, "answers") {
		t.Errorf("decision does not quote the answer and answer the question: %s", out)
	}
	var rd api.ReadyResponse
	if err := json.Unmarshal([]byte(r.rhiza("ready", "--json")), &rd); err != nil {
		t.Fatal(err)
	}
	ready := false
	for _, c := range rd.ChangeSets {
		if c.ID.String() == blocked {
			ready = c.Ready
		}
	}
	if !ready {
		t.Errorf("change set not ready after the answer: %+v", rd)
	}
	r.rhiza("render")
	r.rhiza("check")
	if v := m.View().Content; strings.Contains(v, "blocked by question") {
		t.Errorf("view still shows the block:\n%s", v)
	}
}

// callOver runs an operation in r's repository as the app's runner does.
func callOver(r roadmapRepo) func(g globals, op string, record, dry bool, args []string) ([]byte, error) {
	return func(g globals, op string, record, dry bool, args []string) ([]byte, error) {
		a := []string{op, "--json"}
		if g.changeSet {
			a = append(a, "--change", g.change)
		}
		if dry {
			a = append(a, "--dry-run")
		}
		a = append(a, args...)
		var out, errb bytes.Buffer
		if code := run(a, nil, &out, &errb, r.getenv); code != 0 {
			return nil, &execError{out.String() + errb.String()}
		}
		return out.Bytes(), nil
	}
}

type execError struct{ s string }

func (e *execError) Error() string { return e.s }
