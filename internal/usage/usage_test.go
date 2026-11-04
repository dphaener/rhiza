package usage

import (
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
)

var t0 = time.Date(2026, 11, 2, 8, 0, 0, 0, time.UTC)

func rec(d time.Duration, n int64) adapter.UsageRecord {
	return adapter.UsageRecord{At: t0.Add(d), Tokens: n}
}

func TestFirstWindow(t *testing.T) {
	for _, k := range Kinds {
		ws := Windows(k, []adapter.UsageRecord{rec(0, 10), rec(time.Hour, 5)}, nil)
		if len(ws) != 1 {
			t.Fatalf("%s: %d windows", k, len(ws))
		}
		w := ws[0]
		if !w.Start.Equal(t0) || !w.Reset.Equal(t0.Add(k.Length())) || w.Used != 15 {
			t.Errorf("%s: %+v", k, w)
		}
	}
}

func TestRecordAfterEndStartsNext(t *testing.T) {
	for _, k := range Kinds {
		l := k.Length()
		ws := Windows(k, []adapter.UsageRecord{rec(0, 10), rec(l, 7), rec(l-time.Second, 1)}, nil)
		if len(ws) != 2 {
			t.Fatalf("%s: %d windows", k, len(ws))
		}
		if ws[0].Used != 11 || ws[1].Used != 7 || !ws[1].Start.Equal(t0.Add(l)) {
			t.Errorf("%s: %+v", k, ws)
		}
	}
}

func TestGapStartsAtNextRecord(t *testing.T) {
	for _, k := range Kinds {
		gap := 3 * k.Length()
		ws := Windows(k, []adapter.UsageRecord{rec(0, 1), rec(gap, 2)}, nil)
		if len(ws) != 2 || !ws[1].Start.Equal(t0.Add(gap)) {
			t.Fatalf("%s: %+v", k, ws)
		}
		if _, ok := Current(k, []adapter.UsageRecord{rec(0, 1)}, nil, t0.Add(gap)); ok {
			t.Errorf("%s: window live after its reset", k)
		}
		if w, ok := Current(k, []adapter.UsageRecord{rec(0, 1)}, nil, t0.Add(time.Minute)); !ok || w.Used != 1 {
			t.Errorf("%s: current %+v %v", k, w, ok)
		}
	}
}

func TestLimitResetMovesStart(t *testing.T) {
	for _, k := range Kinds {
		l := k.Length()
		reset := t0.Add(l - time.Hour)
		lim := []Limit{{HitAt: t0.Add(30 * time.Minute), ResetAt: reset}}
		ws := Windows(k, []adapter.UsageRecord{rec(0, 4), rec(time.Minute, 6), rec(l-time.Hour, 9)}, lim)
		// The limit moves the window to end at reset and start a length earlier,
		// before the first record, so the first record stays inside it.
		if len(ws) != 2 {
			t.Fatalf("%s: %+v", k, ws)
		}
		if !ws[0].Reset.Equal(reset) || !ws[0].Start.Equal(reset.Add(-l)) || ws[0].Used != 10 {
			t.Errorf("%s: first %+v", k, ws[0])
		}
		if !ws[1].Start.Equal(reset) || ws[1].Used != 9 {
			t.Errorf("%s: second %+v", k, ws[1])
		}
	}
}

func TestLimitIgnoredWhenResetNotAfterHit(t *testing.T) {
	lim := []Limit{{HitAt: t0.Add(time.Minute), ResetAt: t0}}
	ws := Windows(Session, []adapter.UsageRecord{rec(0, 1)}, lim)
	if len(ws) != 1 || !ws[0].Reset.Equal(t0.Add(5*time.Hour)) {
		t.Errorf("%+v", ws)
	}
}

func TestNoRecords(t *testing.T) {
	if Windows(Session, nil, nil) != nil || Windows(Kind("x"), []adapter.UsageRecord{rec(0, 1)}, nil) != nil {
		t.Error("expected no windows")
	}
}
