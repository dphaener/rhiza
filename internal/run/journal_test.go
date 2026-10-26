package run

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

var t0 = time.Date(2026, 10, 26, 12, 0, 0, 0, time.UTC)

func at(n int) time.Time { return t0.Add(time.Duration(n) * time.Minute) }

func startedEntry(start string, n int) Entry {
	return Entry{Type: TypeStarted, At: at(n), Step: start, Started: &Started{
		Workflow: "wf", Definition: "name = \"wf\"", Change: "CS1", Title: "T",
		Worktree: "/w", PermissionMode: "manual", Agent: "claude-code",
	}}
}

func newStore(t *testing.T) *Store { return NewStore(filepath.Join(t.TempDir(), "runs")) }

func TestRoundTrip(t *testing.T) {
	s := newStore(t)
	first := startedEntry("a", 0)
	if err := s.Create("01RUNAAAA", first); err != nil {
		t.Fatal(err)
	}
	if err := s.Create("01RUNAAAA", first); err == nil {
		t.Fatal("second Create succeeded")
	}
	more := []Entry{
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeEnded, At: at(2), Step: "a", Visit: 1, Outcome: "pass", Next: "b"},
	}
	for _, e := range more {
		if err := s.Append("01RUNAAAA", e); err != nil {
			t.Fatal(err)
		}
	}
	j, err := s.Read("01RUNAAAA")
	if err != nil {
		t.Fatal(err)
	}
	want := append([]Entry{first}, more...)
	if j.ID != "01RUNAAAA" || !reflect.DeepEqual(j.Entries, want) {
		t.Fatalf("got %+v, want %+v", j.Entries, want)
	}
	if err := s.Append("01NOPE", more[0]); !errors.As(err, new(*NotFoundError)) {
		t.Fatalf("Append to missing run: %v", err)
	}
	for _, bad := range []string{"", "a/b", ".hidden"} {
		if err := s.Create(bad, first); err == nil {
			t.Errorf("Create(%q) succeeded", bad)
		}
	}
}

func TestReadPrefix(t *testing.T) {
	s := newStore(t)
	for _, id := range []string{"01ABCDEF", "01ABXYZ0", "01ZZZZZZ"} {
		if err := s.Create(id, startedEntry("a", 0)); err != nil {
			t.Fatal(err)
		}
	}
	j, err := s.Read("01ZZ")
	if err != nil || j.ID != "01ZZZZZZ" {
		t.Fatalf("prefix read = %q, %v", j.ID, err)
	}
}

func TestReadAmbiguousPrefix(t *testing.T) {
	s := newStore(t)
	for _, id := range []string{"01ABCDEF", "01ABXYZ0"} {
		if err := s.Create(id, startedEntry("a", 0)); err != nil {
			t.Fatal(err)
		}
	}
	_, err := s.Read("01AB")
	var amb *AmbiguousError
	if !errors.As(err, &amb) || len(amb.Matches) != 2 {
		t.Fatalf("want AmbiguousError with 2 matches, got %v", err)
	}
}

func TestReadNotFoundPrefix(t *testing.T) {
	s := newStore(t)
	if _, err := s.Read("01NOTHERE"); !errors.As(err, new(*NotFoundError)) {
		t.Fatalf("empty store: want NotFoundError, got %v", err)
	}
	if err := s.Create("01ABCDEF", startedEntry("a", 0)); err != nil {
		t.Fatal(err)
	}
	for _, id := range []string{"01ZZ", "01A", "", "a/b"} {
		if _, err := s.Read(id); !errors.As(err, new(*NotFoundError)) {
			t.Errorf("Read(%q): want NotFoundError, got %v", id, err)
		}
	}
}

func TestTornLineIgnored(t *testing.T) {
	s := newStore(t)
	if err := s.Create("01TORN", startedEntry("a", 0)); err != nil {
		t.Fatal(err)
	}
	f, err := os.OpenFile(s.path("01TORN"), os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(`{"type":"ent`)
	f.Close()
	j, err := s.Read("01TORN")
	if err != nil || len(j.Entries) != 1 {
		t.Fatalf("got %d entries, %v", len(j.Entries), err)
	}
}

func TestFoldFinished(t *testing.T) {
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
		{Type: TypeEnded, At: at(3), Step: "a", Visit: 1, Outcome: "pass", Next: "b"},
		{Type: TypeEntered, At: at(4), Step: "b", Visit: 1},
		{Type: TypeEnded, At: at(5), Step: "b", Visit: 1, Outcome: "fail", Next: "a", Detail: "exit 1"},
		{Type: TypeEntered, At: at(6), Step: "a", Visit: 2},
		{Type: TypeEnded, At: at(7), Step: "a", Visit: 2, Outcome: "pass", End: Succeeded},
		{Type: TypeFinished, At: at(8), End: Succeeded},
	}}
	st := j.State()
	if !st.Ended || st.Result != Succeeded || st.Current != "" || st.CurrentVisit != 0 {
		t.Fatalf("finished state wrong: %+v", st)
	}
	if st.Workflow != "wf" || st.Change != "CS1" || st.Worktree != "/w" || st.Definition == "" || !st.StartedAt.Equal(at(0)) {
		t.Fatalf("header wrong: %+v", st)
	}
	if !reflect.DeepEqual(st.Visits, map[string]int{"a": 2, "b": 1}) {
		t.Fatalf("visits = %v", st.Visits)
	}
	if len(st.Steps) != 3 || st.Steps[0].Session != "S1" || st.Steps[1].Detail != "exit 1" || st.Steps[2].Visit != 2 {
		t.Fatalf("steps = %+v", st.Steps)
	}
	if len(st.Session) != 0 {
		t.Fatalf("session map not cleared: %v", st.Session)
	}
}

func TestFoldFailedDetail(t *testing.T) {
	st := Journal{Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeFinished, At: at(1), End: Failed, Detail: `step "a" went over max_visits 1`},
	}}.State()
	if st.Result != Failed || st.Detail == "" || !st.Ended {
		t.Fatalf("%+v", st)
	}
}

func TestFoldStoppedMidStep(t *testing.T) {
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeEnded, At: at(2), Step: "a", Visit: 1, Outcome: "pass", Next: "w"},
		{Type: TypeEntered, At: at(3), Step: "w", Visit: 1},
		{Type: TypeSession, At: at(4), Step: "w", Session: "S9"},
		{Type: TypeStopped, At: at(5), Step: "w", Detail: "waiting: test -f go"},
	}}
	st := j.State()
	if st.Ended || st.Result != "" || st.Current != "w" || st.CurrentVisit != 1 {
		t.Fatalf("state = %+v", st)
	}
	if st.Stopped != "waiting: test -f go" || st.Session["w"] != "S9" {
		t.Fatalf("stopped/session = %q %v", st.Stopped, st.Session)
	}
	if !reflect.DeepEqual(st.Visits, map[string]int{"a": 1, "w": 1}) || len(st.Steps) != 1 {
		t.Fatalf("visits/steps = %v %+v", st.Visits, st.Steps)
	}
	// Entering again clears the stop.
	j.Entries = append(j.Entries, Entry{Type: TypeEnded, At: at(6), Step: "w", Outcome: "pass", Next: "z"})
	st = j.State()
	if st.Current != "z" || st.Stopped != "" || st.Steps[1].Session != "S9" || st.Steps[1].Visit != 1 {
		t.Fatalf("after end: %+v", st)
	}
}

func TestFoldBlockedAndUnblocked(t *testing.T) {
	b := &Block{Kind: "question", Question: "Ship it?", TabSession: "rex", TabWindow: "w1", TabBlock: "b1", TabTitle: "build"}
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
		{Type: TypeBlocked, At: at(3), Step: "a", Session: "S1", Block: b},
	}}
	st := j.State()
	if !reflect.DeepEqual(st.Blocked, b) || st.BlockedStep != "a" || st.BlockedSession != "S1" {
		t.Fatalf("blocked = %+v step %q session %q", st.Blocked, st.BlockedStep, st.BlockedSession)
	}
	if st.Blocked == b {
		t.Errorf("folded block aliases the entry's")
	}
	j.Entries = append(j.Entries, Entry{Type: TypeUnblocked, At: at(4), Step: "a", Session: "S1"})
	st = j.State()
	if st.Blocked != nil || st.BlockedStep != "" || st.BlockedSession != "" || st.Current != "a" {
		t.Fatalf("after unblocked: %+v", st)
	}
	// A blocked entry with no block still reads blocked.
	j.Entries = append(j.Entries, Entry{Type: TypeBlocked, At: at(5), Step: "a", Session: "S1"})
	if st = j.State(); st.Blocked == nil || *st.Blocked != (Block{}) {
		t.Fatalf("bare blocked = %+v", st.Blocked)
	}
}

func TestFoldBlockClearedByEndedAndFinished(t *testing.T) {
	b := &Block{Kind: "permission", Question: "Run rm?"}
	base := []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
		{Type: TypeBlocked, At: at(3), Step: "a", Session: "S1", Block: b},
	}
	ended := Journal{Entries: append(append([]Entry{}, base...),
		Entry{Type: TypeEnded, At: at(4), Step: "a", Visit: 1, Outcome: "pass", Next: "b"})}.State()
	if ended.Blocked != nil || ended.BlockedSession != "" || ended.Current != "b" {
		t.Fatalf("after ended: %+v", ended)
	}
	// Another step's ended entry leaves the block alone.
	other := Journal{Entries: append(append([]Entry{}, base...),
		Entry{Type: TypeEnded, At: at(4), Step: "z", Outcome: "pass", Next: "a"})}.State()
	if other.Blocked == nil {
		t.Fatalf("another step's end cleared the block")
	}
	finished := Journal{Entries: append(append([]Entry{}, base...),
		Entry{Type: TypeFinished, At: at(4), End: Failed, Detail: "killed"})}.State()
	if finished.Blocked != nil || finished.BlockedStep != "" || !finished.Ended {
		t.Fatalf("after finished: %+v", finished)
	}
}

func TestBlockEntryRoundTrip(t *testing.T) {
	s := newStore(t)
	if err := s.Create("01BLOCK", startedEntry("a", 0)); err != nil {
		t.Fatal(err)
	}
	b := &Block{Kind: "question", Question: "Ship it?", TabSession: "rex", TabWindow: "w1", TabBlock: "b1", TabTitle: "build"}
	more := []Entry{
		{Type: TypeBlocked, At: at(1), Step: "a", Session: "S1", Block: b},
		{Type: TypeUnblocked, At: at(2), Step: "a", Session: "S1"},
	}
	for _, e := range more {
		if err := s.Append("01BLOCK", e); err != nil {
			t.Fatal(err)
		}
	}
	j, err := s.Read("01BLOCK")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(j.Entries[1:], more) {
		t.Fatalf("got %+v, want %+v", j.Entries[1:], more)
	}
	raw, err := os.ReadFile(s.path("01BLOCK"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"type":"blocked"`, `"block":{"kind":"question","question":"Ship it?","tab_session":"rex","tab_window":"w1","tab_block":"b1","tab_title":"build"}`, `"type":"unblocked"`} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("journal lacks %s:\n%s", want, raw)
		}
	}
	if strings.Count(string(raw), `"block":`) != 1 {
		t.Errorf("unblocked entry carries a block:\n%s", raw)
	}
}

func TestFoldFresh(t *testing.T) {
	st := Journal{Entries: []Entry{startedEntry("first", 0)}}.State()
	if st.Current != "first" || st.Ended || st.CurrentVisit != 0 {
		t.Fatalf("%+v", st)
	}
}

func TestListOrder(t *testing.T) {
	s := newStore(t)
	for _, c := range []struct {
		id string
		n  int
	}{{"01C", 5}, {"01B", 1}, {"01A", 5}, {"01D", 0}} {
		if err := s.Create(c.id, startedEntry("a", c.n)); err != nil {
			t.Fatal(err)
		}
	}
	// Lock files and stray files are not runs.
	if _, ok, err := s.TryLock("01B"); err != nil || !ok {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(s.dir, "notes.txt"), nil, 0o644)
	js, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, j := range js {
		got = append(got, j.ID)
	}
	if want := []string{"01D", "01B", "01A", "01C"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("order = %v, want %v", got, want)
	}
	if js, err := newStore(t).List(); err != nil || len(js) != 0 {
		t.Fatalf("empty list = %v, %v", js, err)
	}
}

func TestTryLockExclusion(t *testing.T) {
	s := newStore(t)
	if s.Live("01LOCK") {
		t.Fatal("live before any lock")
	}
	unlock, ok, err := s.TryLock("01LOCK")
	if err != nil || !ok {
		t.Fatalf("first TryLock = %v, %v", ok, err)
	}
	_, ok2, err := s.TryLock("01LOCK")
	if err != nil || ok2 {
		t.Fatalf("second TryLock while held = %v, %v", ok2, err)
	}
	if !s.Live("01LOCK") {
		t.Fatal("not live while held")
	}
	if _, ok, _ := s.TryLock("01OTHER"); !ok {
		t.Fatal("other run's lock refused")
	}
	unlock()
	if s.Live("01LOCK") {
		t.Fatal("live after unlock")
	}
	unlock2, ok, err := s.TryLock("01LOCK")
	if err != nil || !ok {
		t.Fatalf("TryLock after unlock = %v, %v", ok, err)
	}
	unlock2()
}
