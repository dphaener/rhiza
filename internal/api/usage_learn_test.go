package api

import (
	"os"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/session"
	"github.com/dphaener/rhiza/internal/usage"
)

// usageFake is a fake agent that is usage-limited while its transcript says
// "limited" and that reports fixed token records.
type usageFake struct {
	limitAdapter
}

const usageFakeAgent = "usage-fake"

var usageFakeReset time.Time

func (usageFake) Name() string { return usageFakeAgent }
func (usageFake) UsageLimit(t []byte) (adapter.UsageLimit, bool) {
	_, ok := limitAdapter{}.UsageLimit(t)
	return adapter.UsageLimit{ResetAt: usageFakeReset}, ok
}
func (usageFake) UsageRecords(_ func(string) string, from, to time.Time) ([]adapter.UsageRecord, error) {
	base := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	all := []adapter.UsageRecord{
		{At: base.Add(-7 * time.Hour), Tokens: 999}, // an earlier window
		{At: base, Tokens: 100},
		{At: base.Add(30 * time.Minute), Tokens: 250},
		{At: base.Add(2 * time.Hour), Tokens: 50}, // after the hit
	}
	var out []adapter.UsageRecord
	for _, r := range all {
		if !r.At.Before(from) && r.At.Before(to) {
			out = append(out, r)
		}
	}
	return out, nil
}

var usageFakeOnce = func() func() {
	done := false
	return func() {
		if !done {
			adapter.Register(usageFake{})
			done = true
		}
	}
}()

func TestLimitHitsBecomeLearnedLimits(t *testing.T) {
	usageFakeOnce()
	f, _, _ := limitFixture(t)
	f.ctx.Getenv = func(string) string { return "" }
	hit := time.Date(2026, 11, 4, 10, 45, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return hit }
	a := usageFake{}

	// No reset time: nothing is recorded, nothing learned.
	recordLimitHit(f.ctx, a, adapter.UsageLimit{})
	var state daemon.State
	learnLimits(f.ctx, f.common, &state)
	if len(daemon.ReadLimitHits(f.common)) != 0 || len(state.Learned) != 0 {
		t.Fatalf("a hit with no reset time taught %+v", state.Learned)
	}

	// A reset 3 hours out is the session window; its use at the hit is 350.
	recordLimitHit(f.ctx, a, adapter.UsageLimit{ResetAt: hit.Add(3 * time.Hour)})
	learnLimits(f.ctx, f.common, &state)
	if len(state.Learned) != 1 || state.Learned[0].Kind != "session" || state.Learned[0].Tokens != 350 {
		t.Fatalf("learned %+v: want the session limit 350", state.Learned)
	}

	// A reset 4 days out is the weekly window, which holds the earlier use too.
	recordLimitHit(f.ctx, a, adapter.UsageLimit{ResetAt: hit.Add(4 * 24 * time.Hour)})
	learnLimits(f.ctx, f.common, &state)
	if len(state.Learned) != 2 || state.Learned[1].Kind != "weekly" || state.Learned[1].Tokens != 1349 {
		t.Fatalf("learned %+v: want the weekly limit 1349 beside the session one", state.Learned)
	}

	// The setting wins while set; a window with neither has no limit.
	rows := []Setting{{Key: "rhiza.daemon.usage.session", Value: "42"}, {Key: "rhiza.daemon.usage.weekly"}}
	if n, ok := usageLimit(rows, state.Learned, "session"); !ok || n != 42 {
		t.Errorf("session limit = %d, %v: want the setting 42", n, ok)
	}
	if n, ok := usageLimit(rows, state.Learned, "weekly"); !ok || n != 1349 {
		t.Errorf("weekly limit = %d, %v: want the learned 1349", n, ok)
	}
	if _, ok := usageLimit(nil, nil, "session"); ok {
		t.Error("a window with no limit reported one")
	}
}

// A run driver's session stopped at its limit records the hit.
func TestSessionRunRecordsLimitHit(t *testing.T) {
	usageFakeOnce()
	f, s, path := limitFixture(t)
	hit := time.Date(2026, 11, 4, 10, 45, 0, 0, time.UTC)
	f.ctx.Now = func() time.Time { return hit }
	usageFakeReset = hit.Add(3 * time.Hour)
	if err := os.WriteFile(path, []byte("limited"), 0o644); err != nil {
		t.Fatal(err)
	}
	st := session.NewStore(session.Dir(f.common))
	if _, err := st.Update(s.ID, func(r *session.Record) { r.Agent = usageFakeAgent }); err != nil {
		t.Fatal(err)
	}
	useClock(t, removeClock{path})
	if _, err := callOp(t, f.ctx, "session.run", `{"id":"`+s.ID+`"}`); err != nil {
		t.Fatal(err)
	}
	hits := daemon.ReadLimitHits(f.common)
	if len(hits) != 1 || !hits[0].HitAt.Equal(hit) || !hits[0].ResetAt.Equal(hit.Add(3*time.Hour)) {
		t.Errorf("hits %+v: want one at %s resetting 3h later", hits, hit)
	}
}

// A hit teaches only the window kind its reset time gives: a session-kind hit
// leaves a learned weekly limit as it was, and a weekly-kind hit leaves a
// learned session limit as it was.
func TestLimitHitTeachesOnlyItsKind(t *testing.T) {
	usageFakeOnce()
	hit := time.Date(2026, 11, 4, 10, 45, 0, 0, time.UTC)
	earlier := hit.Add(-24 * time.Hour)
	cases := []struct {
		name, hitKind, other string
		reset                time.Time
		hitTokens            int64
	}{
		{"session hit over weekly", "session", "weekly", hit.Add(3 * time.Hour), 350},
		{"weekly hit over session", "weekly", "session", hit.Add(4 * 24 * time.Hour), 1349},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f, _, _ := limitFixture(t)
			f.ctx.Getenv = func(string) string { return "" }
			f.ctx.Wd = t.TempDir() // outside any root, so hits stay in this repo
			f.ctx.Now = func() time.Time { return hit }
			state := daemon.State{Learned: []daemon.LearnedLimit{{Kind: c.other, Tokens: 77, HitAt: earlier}}}
			recordLimitHit(f.ctx, usageFake{}, adapter.UsageLimit{ResetAt: c.reset})
			learnLimits(f.ctx, f.common, &state)
			got := map[string]daemon.LearnedLimit{}
			for _, l := range state.Learned {
				got[l.Kind] = l
			}
			if o := got[c.other]; o.Tokens != 77 || !o.HitAt.Equal(earlier) {
				t.Errorf("the %s limit changed to %+v", c.other, o)
			}
			if h := got[c.hitKind]; h.Tokens != c.hitTokens || !h.HitAt.Equal(hit) {
				t.Errorf("the %s limit is %+v: want %d at the hit", c.hitKind, h, c.hitTokens)
			}
		})
	}
}

func TestLearnedLimitRatchets(t *testing.T) {
	hitAt := time.Date(2026, 11, 1, 10, 0, 0, 0, time.UTC)
	start := time.Date(2026, 11, 4, 10, 0, 0, 0, time.UTC)
	reset := start.Add(7 * 24 * time.Hour)
	win := func(used int64) map[usage.Kind]usage.Window {
		return map[usage.Kind]usage.Window{usage.Weekly: {Kind: usage.Weekly, Start: start, Reset: reset, Used: used}}
	}
	learned := func() *daemon.State {
		return &daemon.State{Learned: []daemon.LearnedLimit{{Kind: "weekly", Tokens: 1000, HitAt: hitAt}}}
	}

	// Use past the limit with no hit raises it, keeping the hit time.
	st := learned()
	raiseLearned(st, win(1500), nil)
	if l := st.Learned[0]; l.Tokens != 1500 || !l.HitAt.Equal(hitAt) {
		t.Errorf("learned %+v: want 1500 at the old hit time", l)
	}

	// Use under the limit leaves it.
	st = learned()
	raiseLearned(st, win(900), nil)
	if st.Learned[0].Tokens != 1000 {
		t.Errorf("learned %+v: want 1000 left alone", st.Learned[0])
	}

	// A hit of the kind inside the window stops the ratchet; one of another
	// kind, or outside the window, does not.
	inWin := daemon.LimitHit{HitAt: start.Add(time.Hour), ResetAt: reset}
	st = learned()
	raiseLearned(st, win(1500), []daemon.LimitHit{inWin})
	if st.Learned[0].Tokens != 1000 {
		t.Errorf("learned %+v: a window holding a hit must not ratchet", st.Learned[0])
	}
	other := []daemon.LimitHit{
		{HitAt: start.Add(time.Hour), ResetAt: start.Add(3 * time.Hour)},
		{HitAt: hitAt, ResetAt: hitAt.Add(6 * 24 * time.Hour)},
	}
	st = learned()
	raiseLearned(st, win(1500), other)
	if st.Learned[0].Tokens != 1500 {
		t.Errorf("learned %+v: hits elsewhere must not stop the ratchet", st.Learned[0])
	}

	// A setting is still what usageLimit returns.
	rows := []Setting{{Key: "rhiza.daemon.usage.weekly", Value: "42"}}
	if n, ok := usageLimit(rows, st.Learned, "weekly"); !ok || n != 42 {
		t.Errorf("weekly limit = %d, %v: want the setting 42", n, ok)
	}
}
