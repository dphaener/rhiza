package schedule

import (
	"strings"
	"testing"
	"time"
)

func at(t *testing.T, s string, loc *time.Location) time.Time {
	t.Helper()
	v, err := time.ParseInLocation("2006-01-02 15:04", s, loc)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

func nexts(t *testing.T, expr string, loc *time.Location, from string, n int) []string {
	t.Helper()
	c, err := ParseCron(expr, loc)
	if err != nil {
		t.Fatalf("%q: %v", expr, err)
	}
	cur := at(t, from, loc)
	var out []string
	for i := 0; i < n; i++ {
		var ok bool
		cur, ok = c.Next(cur)
		if !ok {
			t.Fatalf("%q: no slot", expr)
		}
		out = append(out, cur.In(loc).Format("2006-01-02 15:04"))
	}
	return out
}

func eq(t *testing.T, got, want []string) {
	t.Helper()
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestCronFields(t *testing.T) {
	u := time.UTC
	eq(t, nexts(t, "*/20 * * * *", u, "2026-01-01 00:00", 4),
		[]string{"2026-01-01 00:20", "2026-01-01 00:40", "2026-01-01 01:00", "2026-01-01 01:20"})
	eq(t, nexts(t, "5,35 9 * * *", u, "2026-01-01 09:05", 3),
		[]string{"2026-01-01 09:35", "2026-01-02 09:05", "2026-01-02 09:35"})
	eq(t, nexts(t, "0 9-11 * * *", u, "2026-01-01 08:00", 4),
		[]string{"2026-01-01 09:00", "2026-01-01 10:00", "2026-01-01 11:00", "2026-01-02 09:00"})
	eq(t, nexts(t, "10-30/10 0 * * *", u, "2026-01-01 00:00", 3),
		[]string{"2026-01-01 00:10", "2026-01-01 00:20", "2026-01-01 00:30"})
	// 2026-01-05 is a Monday.
	eq(t, nexts(t, "0 12 * jan-feb mon-wed", u, "2026-01-01 00:00", 3),
		[]string{"2026-01-05 12:00", "2026-01-06 12:00", "2026-01-07 12:00"})
	eq(t, nexts(t, "0 0 * * 7", u, "2026-01-01 00:00", 2),
		[]string{"2026-01-04 00:00", "2026-01-11 00:00"})
	eq(t, nexts(t, "0 0 29 FEB *", u, "2026-01-01 00:00", 1), []string{"2028-02-29 00:00"})
}

func TestCronMacros(t *testing.T) {
	u := time.UTC
	from := "2026-03-04 05:06" // a Wednesday
	eq(t, nexts(t, "@hourly", u, from, 1), []string{"2026-03-04 06:00"})
	eq(t, nexts(t, "@daily", u, from, 1), []string{"2026-03-05 00:00"})
	eq(t, nexts(t, "@weekly", u, from, 1), []string{"2026-03-08 00:00"})
	eq(t, nexts(t, "@monthly", u, from, 1), []string{"2026-04-01 00:00"})
	eq(t, nexts(t, "@yearly", u, from, 1), []string{"2027-01-01 00:00"})
}

func TestCronEitherDayField(t *testing.T) {
	u := time.UTC
	// the 15th or any Friday: 2026-01-02 is a Friday
	eq(t, nexts(t, "0 0 15 * fri", u, "2026-01-01 00:00", 4),
		[]string{"2026-01-02 00:00", "2026-01-09 00:00", "2026-01-15 00:00", "2026-01-16 00:00"})
	// only one restricted: both must match (the star matches all)
	eq(t, nexts(t, "0 0 15 * *", u, "2026-01-01 00:00", 2),
		[]string{"2026-01-15 00:00", "2026-02-15 00:00"})
	// a star with a step is still a star
	// so both must match: a Friday that is day 1, 11, 21 or 31
	eq(t, nexts(t, "0 0 */10 * fri", u, "2026-01-01 00:00", 2),
		[]string{"2026-05-01 00:00", "2026-07-31 00:00"})
}

func TestCronTimeZone(t *testing.T) {
	den, err := time.LoadLocation("America/Denver")
	if err != nil {
		t.Fatal(err)
	}
	c, err := ParseCron("CRON_TZ=America/Denver 0 9 * * *", time.UTC)
	if err != nil {
		t.Fatal(err)
	}
	// 09:00 MST is 16:00 UTC in January, 15:00 UTC in July.
	got, _ := c.Next(time.Date(2026, 1, 10, 0, 0, 0, 0, time.UTC))
	if want := time.Date(2026, 1, 10, 16, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("winter: got %v want %v", got, want)
	}
	got, _ = c.Next(time.Date(2026, 7, 10, 0, 0, 0, 0, time.UTC))
	if want := time.Date(2026, 7, 10, 15, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("summer: got %v want %v", got, want)
	}
	if got.Location() != den && got.Location().String() != "America/Denver" {
		t.Fatalf("zone %v", got.Location())
	}
}

func TestCronDST(t *testing.T) {
	den, _ := time.LoadLocation("America/Denver")
	// Spring forward 2026-03-08: 02:00 -> 03:00; a 02:30 slot does not fire.
	eq(t, nexts(t, "30 2 * * *", den, "2026-03-07 00:00", 3),
		[]string{"2026-03-07 02:30", "2026-03-09 02:30", "2026-03-10 02:30"})
	// Fall back 2026-11-01: 02:00 -> 01:00; a 01:30 slot fires once.
	c, _ := ParseCron("CRON_TZ=America/Denver 30 1 * * *", time.UTC)
	cur := at(t, "2026-10-31 12:00", den)
	var fired []time.Time
	for i := 0; i < 3; i++ {
		cur, _ = c.Next(cur)
		fired = append(fired, cur)
	}
	if n := countOn(fired, den, "2026-11-01"); n != 1 {
		t.Fatalf("fall back fired %d times: %v", n, fired)
	}
	if d := fired[2].Sub(fired[1]); d != 24*time.Hour {
		t.Fatalf("slots %v", fired)
	}
	// An hourly slot in the repeated hour fires once per reading too.
	hc, _ := ParseCron("CRON_TZ=America/Denver 0 * * * *", time.UTC)
	cur = at(t, "2026-11-01 00:30", den)
	var hours []string
	for i := 0; i < 4; i++ {
		cur, _ = hc.Next(cur)
		hours = append(hours, cur.In(den).Format("15:04"))
	}
	eq(t, hours, []string{"01:00", "02:00", "03:00", "04:00"})
}

func countOn(ts []time.Time, loc *time.Location, day string) int {
	n := 0
	for _, t := range ts {
		if t.In(loc).Format("2006-01-02") == day {
			n++
		}
	}
	return n
}

func TestNever(t *testing.T) {
	c, err := ParseCron("0 0 31 feb *", time.UTC)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := c.Next(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)); ok {
		t.Fatal("31 February has a slot")
	}
}

func TestEvery(t *testing.T) {
	e, err := ParseEvery("15m")
	if err != nil {
		t.Fatal(err)
	}
	for _, start := range []string{"2026-05-05 10:01", "2026-05-05 10:15", "2026-05-05 10:59", "2026-05-05 11:00"} {
		cur := at(t, start, time.UTC).Add(7 * time.Second)
		var got []string
		for i := 0; i < 4; i++ {
			cur, _ = e.Next(cur)
			got = append(got, cur.Format("04"))
		}
		for _, m := range got {
			if m != "00" && m != "15" && m != "30" && m != "45" {
				t.Fatalf("start %s: %v", start, got)
			}
		}
	}
	// strictly after: a slot itself yields the next one
	got, _ := e.Next(time.Date(2026, 5, 5, 10, 15, 0, 0, time.UTC))
	if got.Minute() != 30 {
		t.Fatalf("got %v", got)
	}
	// multiples since the epoch: 90m slots are not aligned to the hour
	e90, _ := ParseEvery("90m")
	got, _ = e90.Next(time.Unix(0, 0).UTC())
	if !got.Equal(time.Unix(5400, 0)) {
		t.Fatalf("got %v", got)
	}
}

func TestParseErrors(t *testing.T) {
	for _, expr := range []string{
		"0 0 0 * * *",                    // six fields
		"CRON_TZ=Mars/Olympus 0 9 * * *", // unknown zone
		"* * * *",                        // four
		"60 * * * *", "* 24 * * *", "* * 0 * *", "* * * 13 *", "* * * * 8",
		"*/0 * * * *", "5-1 * * * *", "@sometimes", "x * * * *", "1,,2 * * * *",
	} {
		if _, err := ParseCron(expr, time.UTC); err == nil {
			t.Errorf("%q parsed", expr)
		}
	}
	for _, s := range []string{"30s", "-5m", "0", "59s", "soon", ""} {
		if _, err := ParseEvery(s); err == nil {
			t.Errorf("every %q parsed", s)
		}
	}
	if _, err := ParseEvery("1m"); err != nil {
		t.Error(err)
	}
}
