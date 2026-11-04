package api

import (
	"os"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/session"
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
