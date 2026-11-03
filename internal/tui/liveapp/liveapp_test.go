package liveapp

import (
	"fmt"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/tui/launcher"
)

type fakeData struct {
	ls     string
	daemon string
	status map[string]string
	change map[string]string
	calls  int
	shows  int
	opens  []string
	logs   [][]string
	open   string
}

func (f *fakeData) run(op string, args []string) ([]byte, error) {
	f.calls++
	switch op {
	case "run.ls":
		return []byte(f.ls), nil
	case "daemon.status":
		return []byte(f.daemon), nil
	case "run.open":
		f.opens = append(f.opens, args[0])
		return []byte(f.open), nil
	case "run.log":
		f.logs = append(f.logs, args)
		return []byte(`{"run":"` + args[0] + `","lines":[],"opened":true}`), nil
	case "run.status":
		return []byte(f.status[args[0]]), nil
	case "change.show":
		f.shows++
		return []byte(f.change[args[0]]), nil
	}
	return nil, nil
}

const ls1 = `{"runs":[
{"id":"R2","change":"C2","title":"Second","state":"running","current":"review","started_at":"2026-01-01T00:00:00Z","live":true,"workflow":"ship-change"},
{"id":"R1","change":"C1","title":"First","state":"succeeded","started_at":"2025-12-31T23:00:00Z","live":false,"workflow":"ship-change"}]}`

func newFake() *fakeData {
	return &fakeData{
		ls:     ls1,
		daemon: `{"running":true,"autonomous":false,"runs":[{"id":"R2","change":"C2","state":"running"}],"unjournaled":[]}`,
		status: map[string]string{
			"R2": `{"id":"R2","leaves":[{"leaf":"a","outcome":"passed"},{"leaf":"b","outcome":"passed"},{"leaf":"c"}],"checks":[{"check":"q","outcome":"passed"},{"check":"s","outcome":"failed"},{"check":"t"}]}`,
			"R1": `{"id":"R1","leaves":[{"leaf":"a","outcome":"passed"}],"checks":[]}`,
		},
		change: map[string]string{
			"C1": `{"leaves":[{"handle":"a"}]}`,
			"C2": `{"leaves":[{"handle":"a"},{"handle":"b"},{"handle":"c"},{"handle":"d"}]}`,
		},
	}
}

// load builds the model and runs its leaf fetch, as Init's command would.
func load(f *fakeData, now Clock) *Model {
	m := NewModel(f.run, now)
	settle(m)
	return m
}

// settle runs the pending leaf fetch and hands its message to Update.
func settle(m *Model) {
	if c := m.fetchLeaves(); c != nil {
		m.Update(c())
	}
}

// leavesFrom runs a batch's commands and returns its LeavesMsg; the tick's
// command is left asleep.
func leavesFrom(t *testing.T, cmd tea.Cmd) tea.Msg {
	t.Helper()
	if cmd == nil {
		t.Fatal("no command")
	}
	batch, ok := cmd().(tea.BatchMsg)
	if !ok {
		t.Fatal("not a batch")
	}
	got := make(chan tea.Msg, len(batch))
	for _, c := range batch {
		if c != nil {
			go func() { got <- c() }()
		}
	}
	for {
		select {
		case msg := <-got:
			if _, ok := msg.(LeavesMsg); ok {
				return msg
			}
		case <-time.After(time.Second):
			t.Fatal("no leaf fetch in the batch")
		}
	}
}

func clock(s string) Clock {
	t, _ := time.Parse(time.RFC3339, s)
	return func() time.Time { return t }
}

func TestViewShowsRunProgress(t *testing.T) {
	m := load(newFake(), clock("2026-01-01T00:05:07Z"))
	v := m.View().Content
	for _, want := range []string{
		"Daemon: running · 1 live runs",
		"> Second · running",
		"step review · leaves 2/4 · checks 1 passed, 1 failed, 1 running · elapsed 5m07s",
		"First · succeeded",
		"leaves 1/1 · checks none · elapsed 1h05m07s",
	} {
		if !strings.Contains(v, want) {
			t.Errorf("view lacks %q:\n%s", want, v)
		}
	}
}

func TestTickUpdatesViewKeepingCursorOnRun(t *testing.T) {
	f := newFake()
	m := load(f, clock("2026-01-01T00:05:07Z"))
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	if r, _ := m.Selected(); r.ID != "R1" {
		t.Fatalf("selected %q", r.ID)
	}
	// A newer run appears on top and R2's results change.
	f.ls = `{"runs":[{"id":"R3","change":"C3","title":"Third","state":"running","current":"build","started_at":"2026-01-01T00:05:00Z","live":true,"workflow":"ship-change"},` + ls1[len(`{"runs":[`):]
	f.status["R3"] = `{"id":"R3","leaves":[],"checks":[]}`
	f.status["R2"] = `{"id":"R2","leaves":[{"leaf":"a","outcome":"passed"},{"leaf":"b","outcome":"passed"},{"leaf":"c","outcome":"passed"}],"checks":[]}`
	f.change["C2"] = `{"leaves":[{"handle":"a"},{"handle":"b"},{"handle":"c"}]}`
	f.daemon = `{"running":false,"runs":[],"unjournaled":[]}`
	shows := f.shows
	_, cmd := m.Update(TickMsg{})
	if f.shows != shows {
		t.Fatalf("the tick ran change.show %d times on the update loop", f.shows-shows)
	}
	m.Update(leavesFrom(t, cmd))
	if f.shows != shows+1 {
		t.Fatalf("change.show ran %d times for the one new change set", f.shows-shows)
	}
	v := m.View().Content
	if r, _ := m.Selected(); r.ID != "R1" {
		t.Fatalf("cursor moved to %q", r.ID)
	}
	// C2's total is kept from the first fetch until r fetches it again.
	for _, want := range []string{"Daemon: not running", "Third · running", "leaves 3/4 · checks none"} {
		if !strings.Contains(v, want) {
			t.Errorf("view lacks %q:\n%s", want, v)
		}
	}
	_, cmd = m.Update(tea.KeyPressMsg{Code: 'r', Text: "r"})
	if cmd == nil {
		t.Fatal("r fetched no leaf totals")
	}
	m.Update(cmd())
	if v := m.View().Content; !strings.Contains(v, "leaves 3/3 · checks none") {
		t.Errorf("view lacks refetched leaves 3/3:\n%s", v)
	}
}

func TestEmptyAndRegistered(t *testing.T) {
	f := newFake()
	f.ls = `{"runs":[]}`
	m := NewModel(f.run, nil)
	if !strings.Contains(m.View().Content, "No runs") {
		t.Fatalf("view: %q", m.View().Content)
	}
	found := false
	for _, a := range launcher.Apps() {
		found = found || a.Name == Name
	}
	if !found {
		t.Fatal("live app not in the launcher registry")
	}
}

func TestResultIsRunLsBytes(t *testing.T) {
	f := newFake()
	m := NewModel(f.run, nil)
	if out, err := m.Result(); err != nil || string(out) != ls1 {
		t.Fatalf("result %q, %v", out, err)
	}
}

func enter(m *Model) { m.Update(tea.KeyPressMsg{Code: tea.KeyEnter}) }

func TestEnterFocusesLiveRunTab(t *testing.T) {
	f := newFake()
	m := NewModel(f.run, nil)
	f.open = `{"change":"C2","live":true,"run":"R2","opened":true,"tab":"block:b1"}`
	enter(m)
	if len(f.opens) != 1 || f.opens[0] != "C2" {
		t.Fatalf("run.open calls: %v", f.opens)
	}
	if v := m.View().Content; !strings.Contains(v, "Focused Rex tab block:b1 of run R2") {
		t.Errorf("view lacks focused notice:\n%s", v)
	}
	f.open = `{"change":"C2","live":true,"run":"R2","opened":false,"tab_gone":true}`
	enter(m)
	if v := m.View().Content; !strings.Contains(v, "The tab of run R2 is gone") {
		t.Errorf("view lacks tab gone notice:\n%s", v)
	}
}

func TestEnterOnFinishedRunCallsNothing(t *testing.T) {
	f := newFake()
	m := NewModel(f.run, nil)
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	enter(m)
	if len(f.opens) != 0 {
		t.Fatalf("run.open called: %v", f.opens)
	}
	if v := m.View().Content; !strings.Contains(v, "is not live") {
		t.Errorf("view lacks not-live notice:\n%s", v)
	}
}

func TestLeavesCountChangeSetLeavesOnceEach(t *testing.T) {
	f := newFake()
	// a and b passed; a checks step's fixer then fixed b again, an entry
	// with no outcome. c and d are not conducted yet.
	f.status["R2"] = `{"id":"R2","leaves":[{"step":"build","leaf":"a","outcome":"passed"},{"step":"build","leaf":"b","outcome":"passed"},{"step":"review","leaf":"b"}],"checks":[]}`
	m := load(f, clock("2026-01-01T00:05:07Z"))
	if v := m.View().Content; !strings.Contains(v, "leaves 2/4 · checks none") {
		t.Errorf("view lacks leaves 2/4:\n%s", v)
	}
	// Every leaf passed, with a fixer entry: the count reads complete.
	f.status["R2"] = `{"id":"R2","leaves":[{"step":"build","leaf":"a","outcome":"passed"},{"step":"build","leaf":"b","outcome":"passed"},{"step":"build","leaf":"c","outcome":"passed"},{"step":"build","leaf":"d","outcome":"passed"},{"step":"review","leaf":"b"}],"checks":[]}`
	m.Reload()
	if v := m.View().Content; !strings.Contains(v, "leaves 4/4 · checks none") {
		t.Errorf("view lacks leaves 4/4:\n%s", v)
	}
	// Without change.show the total is the distinct leaves the run lists.
	delete(f.change, "C2")
	m = load(f, clock("2026-01-01T00:05:07Z"))
	if v := m.View().Content; !strings.Contains(v, "leaves 4/4 · checks none") {
		t.Errorf("view lacks fallback leaves 4/4:\n%s", v)
	}
}

func TestViewScrollsToCursorKeepingDaemonLine(t *testing.T) {
	f := newFake()
	var runs []string
	for i := 0; i < 20; i++ {
		id := fmt.Sprintf("R%02d", i)
		runs = append(runs, fmt.Sprintf(`{"id":%q,"change":"C1","title":"Run %02d","state":"succeeded","started_at":"2026-01-01T00:00:00Z","workflow":"ship-change"}`, id, i))
		f.status[id] = `{"id":"x","leaves":[],"checks":[]}`
	}
	f.ls = `{"runs":[` + strings.Join(runs, ",") + `]}`
	// Every height, with and without a notice: the frame fits, so the
	// daemon line stays on top.
	for h := 7; h <= 30; h++ {
		for _, notice := range []bool{false, true} {
			m := load(f, clock("2026-01-01T00:05:07Z"))
			m.Update(tea.WindowSizeMsg{Width: 80, Height: h})
			if notice {
				enter(m) // Run 00 is not live
			}
			lines := func() int { return strings.Count(m.View().Content, "\n") + 1 }
			v := m.View().Content
			if notice && !strings.Contains(v, "is not live") {
				t.Fatalf("height %d: no notice:\n%s", h, v)
			}
			if !strings.HasPrefix(v, "Daemon: running") || !strings.Contains(v, "> Run 00") || strings.Contains(v, "Run 19") || lines() > h {
				t.Fatalf("height %d, notice %v: top view (%d lines):\n%s", h, notice, lines(), v)
			}
			for i := 0; i < 19; i++ {
				m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
			}
			v = m.View().Content
			if !strings.HasPrefix(v, "Daemon: running") || !strings.Contains(v, "> Run 19") || strings.Contains(v, "Run 00") || lines() > h {
				t.Fatalf("height %d, notice %v: bottom view (%d lines):\n%s", h, notice, lines(), v)
			}
		}
	}
}

// A check whose latest cycle is open has an empty outcome in run.status,
// and the checks summary counts it running, not as its earlier outcome.
func TestChecksSummaryCountsOpenCycleRunning(t *testing.T) {
	cs := []api.RunCheck{{Check: "q", Outcome: "passed"}, {Check: "s", OpenCycle: 2, PriorCycle: 1, PriorOutcome: "failed"}}
	if got, want := checksSummary(cs), "1 passed, 1 running"; got != want {
		t.Errorf("checksSummary = %q, want %q", got, want)
	}
}

func TestLKeyOpensSelectedRunLogLayer(t *testing.T) {
	f := newFake()
	m := NewModel(f.run, nil)
	m.Update(tea.KeyPressMsg{Code: 'j', Text: "j"})
	m.Update(tea.KeyPressMsg{Code: 'l', Text: "l"})
	if len(f.logs) != 1 || strings.Join(f.logs[0], " ") != "R1 --layer" {
		t.Fatalf("run.log calls: %v", f.logs)
	}
	if v := m.View().Content; !strings.Contains(v, "Opened the log of run R1") {
		t.Errorf("view lacks notice:\n%s", v)
	}
}
