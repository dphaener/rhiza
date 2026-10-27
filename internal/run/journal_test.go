package run

import (
	"encoding/json"
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

func TestFoldAutonomous(t *testing.T) {
	plain := startedEntry("a", 0)
	if st := (Journal{Entries: []Entry{plain}}).State(); st.Autonomous {
		t.Fatalf("a run started without autonomous folds autonomous: %+v", st)
	}
	b, err := json.Marshal(plain)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "autonomous") {
		t.Fatalf("a plain started entry names autonomous: %s", b)
	}
	auto := startedEntry("a", 0)
	auto.Started.Autonomous = true
	b, err = json.Marshal(auto)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"autonomous":true`) {
		t.Fatalf("an autonomous started entry lacks autonomous: %s", b)
	}
	var back Entry
	if err := json.Unmarshal(b, &back); err != nil {
		t.Fatal(err)
	}
	if st := (Journal{Entries: []Entry{back, {Type: TypeEntered, At: at(1), Step: "a"}}}).State(); !st.Autonomous {
		t.Fatalf("an autonomous run folds not autonomous: %+v", st)
	}
}

func TestFoldAccepts(t *testing.T) {
	es := []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a"},
		{Type: TypeAccepted, At: at(2), Change: "CS1", LogEntry: "E1"},
		{Type: TypeRejected, At: at(3), Change: "CS2", LogEntry: "E2"},
	}
	st := (Journal{Entries: es}).State()
	want := []Accept{
		{Change: "CS1", LogEntry: "E1", Verb: TypeAccepted, At: at(2)},
		{Change: "CS2", LogEntry: "E2", Verb: TypeRejected, At: at(3)},
	}
	if !reflect.DeepEqual(st.Accepts, want) {
		t.Fatalf("accepts = %+v, want %+v", st.Accepts, want)
	}
	if st.Current != "a" || st.Ended || st.Stopped != "" {
		t.Fatalf("accepts changed the fold: %+v", st)
	}
	b, err := json.Marshal(es[2])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), `"change":"CS1","log_entry":"E1"`) {
		t.Fatalf("accepted entry = %s", b)
	}
	if b, _ := json.Marshal(es[1]); strings.Contains(string(b), "change") || strings.Contains(string(b), "log_entry") {
		t.Fatalf("an entered entry names change or log_entry: %s", b)
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

func TestFoldAmended(t *testing.T) {
	b := &Block{Kind: "question", Question: "Ship it?"}
	amend := Entry{Type: TypeAmended, At: at(9), Step: "a", Session: "S1", Actor: "claude-code-f436602f", File: "plan/brief.md"}
	cases := []struct {
		name    string
		entries []Entry
		// wantBlocked and wantStopped are the fold without the amendment.
		wantBlocked bool
		wantStopped string
		wantCurrent string
	}{
		{
			name: "running step",
			entries: []Entry{
				startedEntry("a", 0),
				{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
				{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
			},
			wantCurrent: "a",
		},
		{
			name: "blocked step",
			entries: []Entry{
				startedEntry("a", 0),
				{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
				{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
				{Type: TypeBlocked, At: at(3), Step: "a", Session: "S1", Block: b},
			},
			wantBlocked: true,
			wantCurrent: "a",
		},
		{
			name: "stopped step",
			entries: []Entry{
				startedEntry("a", 0),
				{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
				{Type: TypeStopped, At: at(2), Step: "a", Detail: "waiting: test -f go"},
			},
			wantStopped: "waiting: test -f go",
			wantCurrent: "a",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			before := Journal{ID: "r", Entries: c.entries}.State()
			after := Journal{ID: "r", Entries: append(append([]Entry{}, c.entries...), amend)}.State()
			want := Amendment{Step: "a", Session: "S1", Actor: "claude-code-f436602f", File: "plan/brief.md", At: at(9)}
			if len(after.Amendments) != 1 || after.Amendments[0] != want {
				t.Fatalf("amendments = %+v, want [%+v]", after.Amendments, want)
			}
			if len(before.Amendments) != 0 {
				t.Fatalf("amendments before = %+v", before.Amendments)
			}
			if after.Current != c.wantCurrent || before.Current != c.wantCurrent {
				t.Fatalf("current = %q, want %q", after.Current, c.wantCurrent)
			}
			if after.Stopped != before.Stopped || after.Stopped != c.wantStopped {
				t.Fatalf("stopped = %q, was %q, want %q", after.Stopped, before.Stopped, c.wantStopped)
			}
			if (after.Blocked != nil) != c.wantBlocked || !reflect.DeepEqual(after.Blocked, before.Blocked) {
				t.Fatalf("blocked = %+v, was %+v; want blocked %v", after.Blocked, before.Blocked, c.wantBlocked)
			}
			if after.BlockedStep != before.BlockedStep || after.BlockedSession != before.BlockedSession {
				t.Fatalf("blocked step/session changed: %q %q, was %q %q",
					after.BlockedStep, after.BlockedSession, before.BlockedStep, before.BlockedSession)
			}
			if after.Ended != before.Ended || after.Result != before.Result || !reflect.DeepEqual(after.Visits, before.Visits) || !reflect.DeepEqual(after.Session, before.Session) {
				t.Fatalf("amendment changed other fold fields:\nafter %+v\nbefore %+v", after, before)
			}
		})
	}
}

func TestAmendedRoundTrip(t *testing.T) {
	s := newStore(t)
	if err := s.Create("01AMEND", startedEntry("a", 0)); err != nil {
		t.Fatal(err)
	}
	e := Entry{Type: TypeAmended, At: at(1), Step: "a", Session: "S1", Actor: "claude-code-f436602f", File: "plan/brief.md"}
	if err := s.Append("01AMEND", e); err != nil {
		t.Fatal(err)
	}
	j, err := s.Read("01AMEND")
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Entries) != 2 || !reflect.DeepEqual(j.Entries[1], e) {
		t.Fatalf("got %+v, want %+v", j.Entries, e)
	}
	raw, err := os.ReadFile(s.path("01AMEND"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"type":"amended"`, `"actor":"claude-code-f436602f"`, `"file":"plan/brief.md"`} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("journal lacks %s:\n%s", want, raw)
		}
	}
}

// A retried entry reopens a failed run at its step with no open visit,
// keeps the failure in the history, and marks the next visit of that step
// as the retry. Visits keep counting; the cap base moves to the retry.
func TestFoldRetried(t *testing.T) {
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeEnded, At: at(2), Step: "a", Visit: 1, Outcome: "pass", Next: "b"},
		{Type: TypeEntered, At: at(3), Step: "b", Visit: 1},
		{Type: TypeEnded, At: at(4), Step: "b", Visit: 1, Outcome: "fail", End: Failed, Detail: "exit 1"},
		{Type: TypeFinished, At: at(5), Step: "b", End: Failed, Detail: `step "b" ended with outcome "fail"`},
	}}
	st := j.State()
	if step, err := st.FailedStep(); err != nil || step != "b" {
		t.Fatalf("FailedStep = %q, %v", step, err)
	}
	j.Entries = append(j.Entries, Entry{Type: TypeRetried, At: at(6), Step: "b", Detail: st.Detail})
	st = j.State()
	if st.Ended || st.Result != "" || st.Detail != "" || st.Current != "b" || st.CurrentVisit != 0 || st.FinishedStep != "" {
		t.Fatalf("retried state wrong: %+v", st)
	}
	if _, err := st.FailedStep(); !errors.Is(err, ErrNotFailed) {
		t.Fatalf("FailedStep of a reopened run: %v", err)
	}
	want := []Retry{{Step: "b", Detail: `step "b" ended with outcome "fail"`, At: at(6)}}
	if !reflect.DeepEqual(st.Retries, want) {
		t.Fatalf("retries = %+v", st.Retries)
	}
	if st.CapVisits("b") != 0 || st.CapVisits("a") != 0 || st.Visits["b"] != 1 {
		t.Fatalf("cap visits a=%d b=%d, visits %v", st.CapVisits("a"), st.CapVisits("b"), st.Visits)
	}
	j.Entries = append(j.Entries,
		Entry{Type: TypeEntered, At: at(7), Step: "b", Visit: 2},
		Entry{Type: TypeEnded, At: at(8), Step: "b", Visit: 2, Outcome: "pass", End: Succeeded},
		Entry{Type: TypeFinished, At: at(9), Step: "b", End: Succeeded},
	)
	st = j.State()
	if !st.Ended || st.Result != Succeeded || st.CapVisits("b") != 1 || st.Visits["b"] != 2 {
		t.Fatalf("finished after retry: %+v", st)
	}
	var retried []bool
	for _, s := range st.Steps {
		retried = append(retried, s.Retried)
	}
	if !reflect.DeepEqual(retried, []bool{false, false, true}) {
		t.Fatalf("retried marks = %v", retried)
	}
	if _, err := st.FailedStep(); !errors.Is(err, ErrNotFailed) {
		t.Fatalf("FailedStep of a succeeded run: %v", err)
	}
}

// FailedStep falls back to the last ended step when the finished entry
// names none, and has nothing for an unfinished run.
func TestFailedStep(t *testing.T) {
	st := Journal{Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeEnded, At: at(2), Step: "a", Visit: 1, Outcome: OutcomeError, Detail: "boom"},
		{Type: TypeFinished, At: at(3), End: Failed, Detail: "boom"},
	}}.State()
	if step, err := st.FailedStep(); err != nil || step != "a" {
		t.Fatalf("FailedStep = %q, %v", step, err)
	}
	for name, st := range map[string]State{
		"unfinished":     Journal{Entries: []Entry{startedEntry("a", 0)}}.State(),
		"no step at all": Journal{Entries: []Entry{startedEntry("a", 0), {Type: TypeFinished, At: at(1), End: Failed}}}.State(),
	} {
		if _, err := st.FailedStep(); !errors.Is(err, ErrNotFailed) {
			t.Errorf("%s: err = %v", name, err)
		}
	}
}

func TestRetriedRoundTrip(t *testing.T) {
	s := newStore(t)
	if err := s.Create("01RUNRRRR", startedEntry("a", 0)); err != nil {
		t.Fatal(err)
	}
	e := Entry{Type: TypeRetried, At: at(1), Step: "a", Detail: "why"}
	if err := s.Append("01RUNRRRR", e); err != nil {
		t.Fatal(err)
	}
	j, err := s.Read("01RUNRRRR")
	if err != nil {
		t.Fatal(err)
	}
	if got := j.Entries[1]; !reflect.DeepEqual(got, e) {
		t.Fatalf("got %+v, want %+v", got, e)
	}
}

// Message entries fold into Messages in journal order, with every field.
func TestFoldMessages(t *testing.T) {
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
		{Type: TypeMessage, At: at(3), Step: "a", Session: "S1", Actor: "claude-code-f436602f",
			To: "reviewer", Text: "Please look at the diff.", Summary: "review ask"},
		{Type: TypeMessage, At: at(4), Step: "a", Session: "S1", Actor: "claude-code-f436602f",
			To: "team-lead", Text: `{"type":"shutdown_response"}`, Structured: true},
	}}
	st := j.State()
	want := []Message{
		{Step: "a", Session: "S1", Actor: "claude-code-f436602f", To: "reviewer",
			Text: "Please look at the diff.", Summary: "review ask", At: at(3)},
		{Step: "a", Session: "S1", Actor: "claude-code-f436602f", To: "team-lead",
			Text: `{"type":"shutdown_response"}`, Structured: true, At: at(4)},
	}
	if !reflect.DeepEqual(st.Messages, want) {
		t.Fatalf("messages = %+v, want %+v", st.Messages, want)
	}
}

// The last tokens entry for a session wins.
func TestFoldTokensLastWins(t *testing.T) {
	first := &TokenUse{Session: TokenCounts{Input: 10, Output: 20, CacheCreation: 30, CacheRead: 40}}
	last := &TokenUse{
		Session: TokenCounts{Input: 11, Output: 22, CacheCreation: 33, CacheRead: 44},
		Helpers: TokenCounts{Input: 1, Output: 2, CacheCreation: 3, CacheRead: 4},
	}
	other := &TokenUse{Session: TokenCounts{Input: 7}}
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
		{Type: TypeTokens, At: at(3), Step: "a", Session: "S1", Tokens: first},
		{Type: TypeTokens, At: at(4), Step: "b", Session: "S2", Tokens: other},
		{Type: TypeTokens, At: at(5), Step: "a", Session: "S1", Tokens: last},
		{Type: TypeTokens, At: at(6), Step: "a", Session: "S1"},
	}}
	st := j.State()
	want := map[string]TokenUse{"S1": *last, "S2": *other}
	if !reflect.DeepEqual(st.Tokens, want) {
		t.Fatalf("tokens = %+v, want %+v", st.Tokens, want)
	}
	if fresh := (Journal{ID: "r"}).State(); fresh.Tokens == nil {
		t.Fatal("Tokens map is nil on a fresh fold")
	}
}

// A tokens entry for a session with no session entry adds the session,
// with the step from that entry.
func TestFoldTokensAddsSession(t *testing.T) {
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
		{Type: TypeTokens, At: at(3), Step: "b", Session: "S9", Tokens: &TokenUse{Session: TokenCounts{Input: 5}}},
		{Type: TypeSession, At: at(4), Step: "c", Session: "S9"},
	}}
	st := j.State()
	want := []RunSession{{Step: "a", Session: "S1"}, {Step: "b", Session: "S9"}}
	if !reflect.DeepEqual(st.Sessions, want) {
		t.Fatalf("sessions = %+v, want %+v", st.Sessions, want)
	}
}

// Sessions keeps one entry per distinct session, in first-seen order,
// across a revisited and a retried step.
func TestFoldTokensSessionsOrder(t *testing.T) {
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
		{Type: TypeTokens, At: at(3), Step: "a", Session: "S1", Tokens: &TokenUse{}},
		{Type: TypeEnded, At: at(4), Step: "a", Visit: 1, Outcome: "pass", Next: "b"},
		{Type: TypeEntered, At: at(5), Step: "b", Visit: 1},
		{Type: TypeSession, At: at(6), Step: "b", Session: "S2"},
		{Type: TypeEnded, At: at(7), Step: "b", Visit: 1, Outcome: "fail", Next: "a"},
		// Revisit of a: the same session resumed, then a second session.
		{Type: TypeEntered, At: at(8), Step: "a", Visit: 2},
		{Type: TypeSession, At: at(9), Step: "a", Session: "S1"},
		{Type: TypeSession, At: at(10), Step: "a", Session: "S3"},
		{Type: TypeEnded, At: at(11), Step: "a", Visit: 2, Outcome: "fail", End: Failed},
		{Type: TypeFinished, At: at(12), Step: "a", End: Failed},
		{Type: TypeRetried, At: at(13), Step: "a", Detail: "failed"},
		{Type: TypeEntered, At: at(14), Step: "a", Visit: 3},
		{Type: TypeSession, At: at(15), Step: "a", Session: "S3"},
		{Type: TypeTokens, At: at(16), Step: "a", Session: "S3", Tokens: &TokenUse{}},
	}}
	st := j.State()
	want := []RunSession{{Step: "a", Session: "S1"}, {Step: "b", Session: "S2"}, {Step: "a", Session: "S3"}}
	if !reflect.DeepEqual(st.Sessions, want) {
		t.Fatalf("sessions = %+v, want %+v", st.Sessions, want)
	}
}

// Message and tokens entries change nothing else in the fold: interleaved
// with a journal's other entries, Current, Visits, Steps, Blocked, Stopped
// and the rest come out as they do without them.
func TestFoldMessagesAndTokensChangeNothingElse(t *testing.T) {
	b := &Block{Kind: "question", Question: "Ship it?"}
	msg := func(n int, step, sess string) Entry {
		return Entry{Type: TypeMessage, At: at(n), Step: step, Session: sess, Actor: "claude-code-f436602f", To: "x", Text: "hi"}
	}
	tok := func(n int, step, sess string) Entry {
		return Entry{Type: TypeTokens, At: at(n), Step: step, Session: sess, Tokens: &TokenUse{Session: TokenCounts{Input: int64(n)}}}
	}
	cases := []struct {
		name    string
		entries []Entry
		extra   map[int][]Entry // inserted after the entry at that index
	}{
		{
			name: "finished after a revisit",
			entries: []Entry{
				startedEntry("a", 0),
				{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
				{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
				{Type: TypeEnded, At: at(3), Step: "a", Visit: 1, Outcome: "pass", Next: "b"},
				{Type: TypeEntered, At: at(4), Step: "b", Visit: 1},
				{Type: TypeEnded, At: at(5), Step: "b", Visit: 1, Outcome: "fail", Next: "a"},
				{Type: TypeEntered, At: at(6), Step: "a", Visit: 2},
				{Type: TypeEnded, At: at(7), Step: "a", Visit: 2, Outcome: "pass", End: Succeeded},
				{Type: TypeFinished, At: at(8), End: Succeeded},
			},
			extra: map[int][]Entry{
				2: {msg(20, "a", "S1"), tok(21, "a", "S1")},
				3: {tok(22, "a", "S1")},
				5: {msg(23, "b", "S2")},
				8: {tok(24, "a", "S1")},
			},
		},
		{
			name: "blocked step",
			entries: []Entry{
				startedEntry("a", 0),
				{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
				{Type: TypeSession, At: at(2), Step: "a", Session: "S1"},
				{Type: TypeBlocked, At: at(3), Step: "a", Session: "S1", Block: b},
			},
			extra: map[int][]Entry{
				2: {tok(20, "a", "S1")},
				3: {msg(21, "a", "S1"), tok(22, "a", "S1")},
			},
		},
		{
			name: "stopped step",
			entries: []Entry{
				startedEntry("a", 0),
				{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
				{Type: TypeStopped, At: at(2), Step: "a", Detail: "waiting: test -f go"},
			},
			extra: map[int][]Entry{
				2: {msg(20, "a", "S1"), tok(21, "a", "S1")},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var with []Entry
			for i, e := range c.entries {
				with = append(with, e)
				with = append(with, c.extra[i]...)
			}
			before := Journal{ID: "r", Entries: c.entries}.State()
			after := Journal{ID: "r", Entries: with}.State()
			if len(after.Messages) == 0 && len(after.Tokens) == 0 {
				t.Fatal("the interleaved entries folded to nothing")
			}
			after.Messages, after.Sessions, after.Tokens = nil, nil, nil
			before.Messages, before.Sessions, before.Tokens = nil, nil, nil
			if !reflect.DeepEqual(after, before) {
				t.Fatalf("message/tokens entries changed the fold:\nafter  %+v\nbefore %+v", after, before)
			}
		})
	}
}

// A tokens entry round-trips through the store with the documented keys,
// and a line written before these fields existed still decodes.
func TestTokensRoundTrip(t *testing.T) {
	s := newStore(t)
	if err := s.Create("01TOKENS", startedEntry("a", 0)); err != nil {
		t.Fatal(err)
	}
	e := Entry{Type: TypeTokens, At: at(1), Step: "a", Session: "S1", Tokens: &TokenUse{
		Session: TokenCounts{Input: 1, Output: 2, CacheCreation: 3, CacheRead: 4},
		Helpers: TokenCounts{Input: 5, Output: 6, CacheCreation: 7, CacheRead: 8},
	}}
	m := Entry{Type: TypeMessage, At: at(2), Step: "a", Session: "S1", Actor: "claude-code-f436602f",
		To: "reviewer", Text: "hi", Summary: "greet", Structured: true}
	for _, x := range []Entry{e, m} {
		if err := s.Append("01TOKENS", x); err != nil {
			t.Fatal(err)
		}
	}
	j, err := s.Read("01TOKENS")
	if err != nil {
		t.Fatal(err)
	}
	if len(j.Entries) != 3 || !reflect.DeepEqual(j.Entries[1], e) || !reflect.DeepEqual(j.Entries[2], m) {
		t.Fatalf("got %+v, want %+v and %+v", j.Entries, e, m)
	}
	raw, err := os.ReadFile(s.path("01TOKENS"))
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		`"type":"tokens"`,
		`"tokens":{"session":{"input":1,"output":2,"cache_creation":3,"cache_read":4},"helpers":{"input":5,"output":6,"cache_creation":7,"cache_read":8}}`,
		`"type":"message"`, `"to":"reviewer"`, `"text":"hi"`, `"summary":"greet"`, `"structured":true`,
	} {
		if !strings.Contains(string(raw), want) {
			t.Errorf("journal lacks %s:\n%s", want, raw)
		}
	}
	// The started line carries none of the new keys.
	first := strings.SplitN(string(raw), "\n", 2)[0]
	for _, key := range []string{`"to"`, `"text"`, `"summary"`, `"structured"`, `"tokens"`} {
		if strings.Contains(first, key) {
			t.Errorf("started line has %s: %s", key, first)
		}
	}

	old := `{"type":"ended","at":"2026-10-26T12:03:00Z","step":"a","visit":1,"outcome":"pass","next":"b","session":"S1"}`
	var got Entry
	if err := json.Unmarshal([]byte(old), &got); err != nil {
		t.Fatal(err)
	}
	wantOld := Entry{Type: TypeEnded, At: at(3), Step: "a", Visit: 1, Outcome: "pass", Next: "b", Session: "S1"}
	if !reflect.DeepEqual(got, wantOld) {
		t.Fatalf("old line decoded to %+v, want %+v", got, wantOld)
	}
}

// Each findings entry opens the next round, an empty one too; a verdict
// entry adds its verdicts to the latest round, and the latest verdict on a
// finding wins. Rounds change nothing else in the fold.
func TestFoldRounds(t *testing.T) {
	f1 := Finding{ID: "f1", Node: "N1", Handle: "cap-a", Kind: KindMissingAcceptance, Reason: "no acceptance"}
	f2 := Finding{ID: "f2", Node: "N2", Kind: KindContradiction, Reason: "says both"}
	es := []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a"},
		{Type: TypeVerdict, At: at(2), Verdicts: []Verdict{{Finding: "f1", Verdict: Confirmed, Reason: "orphan"}}},
		{Type: TypeFindings, At: at(3), Findings: []Finding{f1, f2}},
		{Type: TypeVerdict, At: at(4), Verdicts: []Verdict{
			{Finding: "f1", Verdict: Confirmed, Reason: "it has none"},
			{Finding: "f2", Verdict: Confirmed, Reason: "first look"},
		}},
		{Type: TypeVerdict, At: at(5), Verdicts: []Verdict{{Finding: "f2", Verdict: Dropped, Reason: "they agree"}}},
		{Type: TypeFindings, At: at(6)},
	}
	st := (Journal{Entries: es}).State()
	if len(st.Rounds) != 2 {
		t.Fatalf("rounds = %+v, want 2", st.Rounds)
	}
	r1 := st.Rounds[0]
	if r1.Number != 1 || !r1.At.Equal(at(3)) || !reflect.DeepEqual(r1.Findings, []Finding{f1, f2}) || len(r1.Verdicts) != 3 {
		t.Fatalf("round 1 = %+v", r1)
	}
	if v, ok := r1.Verdict("f2"); !ok || v.Verdict != Dropped || v.Reason != "they agree" {
		t.Fatalf("f2 verdict = %+v, %v", v, ok)
	}
	if _, ok := r1.Finding("f3"); ok {
		t.Fatal("round 1 holds f3")
	}
	if got := r1.Confirmed(); !reflect.DeepEqual(got, []Finding{f1}) {
		t.Fatalf("confirmed = %+v, want f1", got)
	}
	r2, ok := st.LatestRound()
	if !ok || r2.Number != 2 || len(r2.Findings) != 0 || len(r2.Confirmed()) != 0 {
		t.Fatalf("latest round = %+v, %v", r2, ok)
	}
	if st.Current != "a" || st.Ended || st.Stopped != "" {
		t.Fatalf("rounds changed the fold: %+v", st)
	}
	if _, ok := (Journal{Entries: es[:3]}).State().LatestRound(); ok {
		t.Fatal("a journal with no findings entry has a round")
	}
	b, err := json.Marshal(es[3])
	if err != nil {
		t.Fatal(err)
	}
	want := `{"type":"findings","at":"2026-10-26T12:03:00Z","findings":[` +
		`{"id":"f1","node":"N1","handle":"cap-a","kind":"missing-acceptance","reason":"no acceptance"},` +
		`{"id":"f2","node":"N2","kind":"contradiction","reason":"says both"}]}`
	if string(b) != want {
		t.Fatalf("findings entry = %s\nwant %s", b, want)
	}
	var back Entry
	if err := json.Unmarshal(b, &back); err != nil || !reflect.DeepEqual(back, es[3]) {
		t.Fatalf("round trip = %+v, %v", back, err)
	}
	if !ValidKind(KindOutsideIntent) || ValidKind("style") {
		t.Fatal("ValidKind")
	}
}

// Two sessions block and unblock independently: each has its own block,
// Blocked is the latest open one, and an unblocked entry ends only its
// session's block.
func TestFoldBlocksPerSession(t *testing.T) {
	j := Journal{ID: "r", Entries: []Entry{
		startedEntry("a", 0),
		{Type: TypeEntered, At: at(1), Step: "a", Visit: 1},
		{Type: TypeBlocked, At: at(2), Step: "a", Session: "S1", Block: &Block{Kind: "permission"}},
		{Type: TypeBlocked, At: at(3), Step: "a", Session: "S2", Block: &Block{Kind: "question", Question: "Q?"}},
		{Type: TypeAsked, At: at(4), Step: "a", Session: "S2", Question: "Q1"},
	}}
	st := j.State()
	if len(st.Blocks) != 2 || st.BlockedSession != "S2" || st.Blocked.Kind != "question" {
		t.Fatalf("two blocks: %+v latest %q", st.Blocks, st.BlockedSession)
	}
	if _, ok := st.BlockAsked("S1"); ok {
		t.Errorf("S1's block was asked")
	}
	if a, ok := st.BlockAsked("S2"); !ok || a.Question != "Q1" {
		t.Errorf("S2's block asked = %+v %v", a, ok)
	}
	j.Entries = append(j.Entries, Entry{Type: TypeUnblocked, At: at(5), Step: "a", Session: "S2"})
	st = j.State()
	if b, ok := st.BlockOf("S1"); !ok || b.Block.Kind != "permission" || st.BlockedSession != "S1" || len(st.Blocks) != 1 {
		t.Fatalf("after S2 unblocked: %+v latest %q", st.Blocks, st.BlockedSession)
	}
	if _, ok := st.BlockOf("S2"); ok {
		t.Errorf("S2 still blocked")
	}
	j.Entries = append(j.Entries,
		Entry{Type: TypeBlocked, At: at(6), Step: "a", Session: "S2", Block: &Block{Kind: "permission"}},
		Entry{Type: TypeUnblocked, At: at(7), Step: "a", Session: "S1"})
	if st = j.State(); len(st.Blocks) != 1 || st.BlockedSession != "S2" {
		t.Fatalf("after S1 unblocked: %+v", st.Blocks)
	}
	if _, ok := st.BlockAsked("S2"); ok {
		t.Errorf("S2's new block carries the old block's asked entry")
	}
	j.Entries = append(j.Entries, Entry{Type: TypeEnded, At: at(8), Step: "a", Visit: 1, Outcome: "pass", Next: "b"})
	if st = j.State(); len(st.Blocks) != 0 || st.Blocked != nil {
		t.Fatalf("after ended: %+v", st.Blocks)
	}
}

// A findings entry with Leaf set folds onto the leaf, by attempt, and
// opens no plan review round; a plan findings entry still does.
func TestFoldLeafFindings(t *testing.T) {
	j := Journal{ID: "R", Entries: []Entry{
		{Type: TypeLeafSession, Step: "build", Leaf: "L1", Handle: "leaf-a", Role: "verifier", Attempt: 1, Session: "S1"},
		{Type: TypeFindings, Step: "build", Leaf: "L1", Handle: "leaf-a", Attempt: 1, Findings: []Finding{{ID: "f1", Reason: "r1"}}},
		{Type: TypeFindings, Step: "build", Leaf: "L1", Handle: "leaf-a", Attempt: 2, Findings: []Finding{{ID: "f1", Reason: "r2"}}},
		{Type: TypeFindings, Findings: []Finding{{ID: "f1", Node: "N", Kind: "contradiction", Reason: "plan"}}},
	}}
	st := j.State()
	if len(st.Rounds) != 1 || st.Rounds[0].Findings[0].Reason != "plan" {
		t.Fatalf("Rounds = %+v", st.Rounds)
	}
	l, ok := st.Leaf("build", "L1")
	if !ok {
		t.Fatal("no leaf")
	}
	if got := l.AttemptFindings(1); len(got) != 1 || got[0].Reason != "r1" {
		t.Errorf("attempt 1 findings = %+v", got)
	}
	if got := l.AttemptFindings(2); len(got) != 1 || got[0].Reason != "r2" {
		t.Errorf("attempt 2 findings = %+v", got)
	}
}

// Findings and verdict entries naming a check go to that check's rounds,
// each round tagged with its cycle (the entry's, else the cycle of the
// check's latest session); a verdict joins its check's latest round. They
// leave the plan review rounds alone, and a findings entry with no check
// still opens a plan review round.
func TestFoldCheckRounds(t *testing.T) {
	fa := Finding{ID: "f1", Check: "alpha", Severity: SeverityMajor, Leaf: "L1", LeafHandle: "leaf-a", Place: "a.go:3", Reason: "off by one"}
	fb := Finding{ID: "f1", Check: "beta", Severity: SeverityNit, Place: "README.md:1", Reason: "typo"}
	pf := Finding{ID: "f1", Node: "N1", Kind: KindContradiction, Reason: "says both"}
	es := []Entry{
		startedEntry("review", 0),
		{Type: TypeEntered, At: at(1), Step: "review"},
		{Type: TypeCheckCycle, At: at(1), Step: "review", Cycle: 1, Checks: []string{"alpha", "beta"}},
		{Type: TypeCheckSession, At: at(2), Step: "review", Check: "alpha", Role: "check", Cycle: 1, Session: "s1"},
		{Type: TypeCheckSession, At: at(2), Step: "review", Check: "beta", Role: "check", Cycle: 1, Session: "s2"},
		{Type: TypeFindings, At: at(3), Step: "review", Check: "alpha", Cycle: 1, Findings: []Finding{fa}},
		{Type: TypeFindings, At: at(4), Step: "review", Check: "beta", Findings: []Finding{fb}},
		{Type: TypeVerdict, At: at(5), Step: "review", Check: "alpha", Verdicts: []Verdict{{Finding: "f1", Verdict: Dropped, Reason: "fine"}}},
		{Type: TypeFindings, At: at(6), Findings: []Finding{pf}},
		{Type: TypeFindings, At: at(7), Step: "review", Check: "beta", Findings: []Finding{}},
	}
	st := (Journal{Entries: es}).State()
	if len(st.Rounds) != 1 || !reflect.DeepEqual(st.Rounds[0].Findings, []Finding{pf}) || len(st.Rounds[0].Verdicts) != 0 {
		t.Fatalf("plan rounds = %+v", st.Rounds)
	}
	a, ok := st.Check("review", "alpha")
	if !ok || len(a.Rounds) != 1 {
		t.Fatalf("alpha = %+v, %v", a, ok)
	}
	ra, ok := a.LatestRound()
	if !ok || ra.Number != 1 || ra.Cycle != 1 || !ra.At.Equal(at(3)) || !reflect.DeepEqual(ra.Findings, []Finding{fa}) {
		t.Fatalf("alpha round = %+v, %v", ra, ok)
	}
	if v, ok := ra.Verdict("f1"); !ok || v.Verdict != Dropped {
		t.Fatalf("alpha verdict = %+v, %v", v, ok)
	}
	b, _ := st.Check("review", "beta")
	if len(b.Rounds) != 2 || b.Rounds[0].Cycle != 1 || b.Rounds[1].Number != 2 || len(b.Rounds[1].Findings) != 0 || len(b.Rounds[0].Verdicts) != 0 {
		t.Fatalf("beta rounds = %+v", b.Rounds)
	}
	if rd, ok := b.RoundFor(1); !ok || rd.Number != 2 {
		t.Fatalf("beta RoundFor(1) = %+v, %v", rd, ok)
	}
	if _, ok := b.RoundFor(2); ok {
		t.Fatal("beta has a round in cycle 2")
	}
	if _, ok := (CheckState{}).LatestRound(); ok {
		t.Fatal("an empty check has a round")
	}
	raw, err := json.Marshal(es[5])
	if err != nil {
		t.Fatal(err)
	}
	want := `{"type":"findings","at":"2026-10-26T12:03:00Z","step":"review","findings":[` +
		`{"id":"f1","check":"alpha","severity":"major","leaf":"L1","leaf_handle":"leaf-a","place":"a.go:3","reason":"off by one"}],` +
		`"check":"alpha","cycle":1}`
	if string(raw) != want {
		t.Fatalf("findings entry = %s\nwant %s", raw, want)
	}
	if !ValidSeverity(SeverityBlocker) || ValidSeverity("critical") || !Serious(SeverityMajor) || Serious(SeverityMinor) || Serious(SeverityNit) {
		t.Fatal("severities")
	}
}
