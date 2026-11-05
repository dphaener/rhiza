package api

import (
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// TestAddQuestionLenientRefs: a run session's question is an engine write
// of free text (decision text-ref-engine-writes): an @word naming no node
// never refuses it, is stored as a literal @ and shown as @@word, while an
// @<existing handle> is stored as @{<full ID>} and shown by its handle.
func TestAddQuestionLenientRefs(t *testing.T) {
	f := askFixture(t, true)
	cs := currentOf(t, f.ctx)
	// The question is written in the session's worktree, so the node it
	// names is added there.
	wctx := testCtx(f.worktree)
	wctx.Git, wctx.ChangeSet = f.git.run, cs
	store := mustWrite(t, wctx, "add", `{"kind":"capability","handle":"store"}`).Nodes[0]
	r := &stepRunner{ctx: f.ctx, change: cs.String(), sessions: f.st}
	q, err := r.AddQuestion(f.sid, "Use @store, not @ghost?")
	if err != nil {
		t.Fatalf("AddQuestion refused an unresolvable @word: %v", err)
	}
	m := ctxWorkspace(t, wctx).model(cs)
	n, ok := m.HandleTaken(q.Handle)
	if !ok {
		t.Fatalf("no question %s", q.Handle)
	}
	node, _ := m.Node(n)
	if got, want := fieldText(node.Fields["statement"]), "Use @{"+store.String()+"}, not @ghost?"; got != want {
		t.Errorf("stored statement %q, want %q", got, want)
	}
	sctx := testCtx(f.worktree)
	sctx.ChangeSet, sctx.IndexPath = cs, filepath.Join(t.TempDir(), "index.db")
	resp, err := callOp(t, sctx, "show", `{"ref":"`+q.Handle+`","view":"what-if:`+cs.String()+`"}`)
	if err != nil {
		t.Fatal(err)
	}
	out := human(t, resp)
	if !strings.Contains(out, "Use @store, not @@ghost?") {
		t.Errorf("shown:\n%s", out)
	}
	if regexp.MustCompile(`@\{`).MatchString(out) {
		t.Errorf("shown with a raw @{ID}:\n%s", out)
	}
}
