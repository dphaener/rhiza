package api

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/session"
)

// After a check session is judged, the log entry files its actor wrote in
// the run's worktree (a QA check's Evidence, here an add and a link) are
// committed on HEAD as one chore(rhiza) commit naming the check, with the
// run's Change: and the session's Actor: trailers and no Leaf: trailer,
// with the entries recording that actor. Another actor's uncommitted
// write stays as it was. A second call commits nothing and returns "".
func TestCommitCheckWrites(t *testing.T) {
	r, git, ctx, _, ada := evidenceCommitRunner(t)
	if err := r.sessions.Write(session.Record{ID: "S-qa", Actor: ada.String(), Change: r.change, Check: "qa", Role: "check", Cycle: 1}); err != nil {
		t.Fatal(err)
	}
	v := *ctx
	v.Actor = ada
	ev := mustWrite(t, &v, "add", `{"kind":"evidence","handle":"leaf-a-qa-1","fields":["statement=ran it","result=PASS"]}`)
	ln := mustWrite(t, &v, "link", `{"from":"leaf-a-qa-1","type":"verifies","to":"leaf-a"}`)
	other := mustWrite(t, ctx, "add", `{"kind":"capability","handle":"leaf-b"}`)
	otherPath := fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, other.Entry)
	want := []string{
		fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, ev.Entry),
		fmt.Sprintf(".rhiza/log/%s/%s.json", r.change, ln.Entry),
	}
	for _, f := range strings.Split(git("ls-files", "--others", "--", ".rhiza/log"), "\n") {
		if f != "" && !slices.Contains(want, f) && f != otherPath {
			want = append(want, f)
		}
	}
	slices.Sort(want)
	old := git("rev-parse", "HEAD")

	c, err := r.CommitCheckWrites("S-qa")
	if err != nil {
		t.Fatal(err)
	}
	if c == "" || git("rev-parse", "HEAD") != c {
		t.Fatalf("commit %q, HEAD %s", c, git("rev-parse", "HEAD"))
	}
	if parents := git("rev-list", "--parents", "-n1", c); parents != c+" "+old {
		t.Errorf("parents = %q, want one new commit on %s", parents, old)
	}
	files := strings.Split(git("show", "--name-only", "--format=", c), "\n")
	if !reflect.DeepEqual(files, want) {
		t.Errorf("committed files = %v, want %v", files, want)
	}
	msg := git("log", "-1", "--format=%B", c)
	if subject := strings.SplitN(msg, "\n", 2)[0]; subject != "chore(rhiza): record qa check writes" {
		t.Errorf("subject = %q", subject)
	}
	trailers := git("log", "-1", "--format=%(trailers:unfold)", c)
	if w := "Change: " + r.change + "\nActor: " + nodeName(ctxWorkspace(t, ctx).model(id.ID{}), ada); trailers != w {
		t.Errorf("trailers = %q, want %q", trailers, w)
	}
	if s := git("status", "--porcelain", "--untracked-files=all", "--", ".rhiza"); s != "?? "+otherPath {
		t.Errorf(".rhiza status = %q, want only %s untracked", s, otherPath)
	}

	again, err := r.CommitCheckWrites("S-qa")
	if err != nil || again != "" {
		t.Errorf("second call = %q, %v; want no commit", again, err)
	}
	if head := git("rev-parse", "HEAD"); head != c {
		t.Errorf("second call moved HEAD to %s", head)
	}
}

// A check session whose actor wrote nothing makes no commit; a critic's
// commit names its role.
func TestCommitCheckWritesNone(t *testing.T) {
	r, git, ctx, _, ada := evidenceCommitRunner(t)
	if err := r.sessions.Write(session.Record{ID: "S-c", Actor: ada.String(), Change: r.change, Check: "quality", Role: "critic", Cycle: 1}); err != nil {
		t.Fatal(err)
	}
	// ada's actor record is committed already (here ada recorded itself,
	// where a session's launcher records its actor).
	git("add", "-A", ".rhiza")
	git("commit", "-q", "-m", "ada")
	mustWrite(t, ctx, "add", `{"kind":"capability","handle":"leaf-b"}`)
	old := git("rev-parse", "HEAD")
	c, err := r.CommitCheckWrites("S-c")
	if err != nil || c != "" {
		t.Errorf("CommitCheckWrites = %q, %v; want no commit", c, err)
	}
	if head := git("rev-parse", "HEAD"); head != old {
		t.Errorf("HEAD moved to %s", head)
	}
	v := *ctx
	v.Actor = ada
	mustWrite(t, &v, "add", `{"kind":"capability","handle":"leaf-c"}`)
	if c, err = r.CommitCheckWrites("S-c"); err != nil || c == "" {
		t.Fatalf("CommitCheckWrites = %q, %v", c, err)
	}
	if subject := git("log", "-1", "--format=%s", c); subject != "chore(rhiza): record quality critic writes" {
		t.Errorf("subject = %q", subject)
	}
}
