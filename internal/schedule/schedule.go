// Package schedule computes the slots of a schedule: the instants a scheduled
// run is due. It is pure: it reads no clock and no file.
//
// A schedule is either a cron expression (the five-field dialect, see
// ParseCron) or a fixed interval (see ParseEvery).
package schedule

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata" // zone names resolve on machines without a zoneinfo database
)

// Schedule yields slots.
type Schedule interface {
	// Next returns the first slot strictly after t. The bool is false when
	// the schedule has no slot in the foreseeable future (a cron expression
	// naming a day that never exists, such as 31 February).
	Next(t time.Time) (time.Time, bool)
}

// MinEvery is the shortest interval an every schedule accepts.
const MinEvery = time.Minute

// horizon bounds the search for a cron slot.
const horizon = 8 // years

// Every is a schedule whose slots are the multiples of its interval since the
// Unix epoch, in UTC.
type Every struct{ d time.Duration }

// ParseEvery reads a Go duration of at least one minute.
func ParseEvery(s string) (Every, error) {
	d, err := time.ParseDuration(strings.TrimSpace(s))
	if err != nil {
		return Every{}, fmt.Errorf("every %q: not a duration", s)
	}
	if d < MinEvery {
		return Every{}, fmt.Errorf("every %q: must be at least %s", s, MinEvery)
	}
	return Every{d}, nil
}

// Interval returns the schedule's duration.
func (e Every) Interval() time.Duration { return e.d }

// Next implements Schedule.
func (e Every) Next(t time.Time) (time.Time, bool) {
	n := t.UnixNano()
	d := int64(e.d)
	k := n / d
	if n < 0 && n%d != 0 {
		k-- // floor
	}
	return time.Unix(0, (k+1)*d).UTC(), true
}

// Cron is a parsed five-field cron expression.
type Cron struct {
	minute, hour, dom, month, dow uint64
	domStar, dowStar              bool
	loc                           *time.Location
}

var macros = map[string]string{
	"@hourly":   "0 * * * *",
	"@daily":    "0 0 * * *",
	"@midnight": "0 0 * * *",
	"@weekly":   "0 0 * * 0",
	"@monthly":  "0 0 1 * *",
	"@yearly":   "0 0 1 1 *",
	"@annually": "0 0 1 1 *",
}

var monthNames = []string{"jan", "feb", "mar", "apr", "may", "jun", "jul", "aug", "sep", "oct", "nov", "dec"}
var dayNames = []string{"sun", "mon", "tue", "wed", "thu", "fri", "sat"}

// ParseCron reads a cron expression: minute hour day-of-month month
// day-of-week, with *, lists, ranges, steps, month and weekday names, or one
// of @hourly, @daily, @weekly, @monthly and @yearly. A leading CRON_TZ=<zone>
// names the IANA zone the expression is read in; without it local is used
// (the machine's local zone when nil). When both day fields are restricted a
// slot matches either.
func ParseCron(expr string, local *time.Location) (*Cron, error) {
	if local == nil {
		local = time.Local
	}
	loc := local
	fields := strings.Fields(expr)
	if len(fields) > 0 && strings.HasPrefix(fields[0], "CRON_TZ=") {
		name := strings.TrimPrefix(fields[0], "CRON_TZ=")
		z, err := time.LoadLocation(name)
		if err != nil || name == "" || name == "Local" {
			return nil, fmt.Errorf("cron %q: unknown time zone %q", expr, name)
		}
		loc = z
		fields = fields[1:]
	}
	if len(fields) == 1 && strings.HasPrefix(fields[0], "@") {
		m, ok := macros[fields[0]]
		if !ok {
			return nil, fmt.Errorf("cron %q: unknown macro %q", expr, fields[0])
		}
		fields = strings.Fields(m)
	}
	if len(fields) != 5 {
		return nil, fmt.Errorf("cron %q: want 5 fields (minute hour day-of-month month day-of-week), got %d", expr, len(fields))
	}
	c := &Cron{loc: loc}
	var err error
	specs := []struct {
		name     string
		lo, hi   int
		names    []string
		namesLo  int
		dst      *uint64
		star     *bool
		sundayHi bool
	}{
		{"minute", 0, 59, nil, 0, &c.minute, nil, false},
		{"hour", 0, 23, nil, 0, &c.hour, nil, false},
		{"day-of-month", 1, 31, nil, 0, &c.dom, &c.domStar, false},
		{"month", 1, 12, monthNames, 1, &c.month, nil, false},
		{"day-of-week", 0, 7, dayNames, 0, &c.dow, &c.dowStar, true},
	}
	for i, s := range specs {
		*s.dst, err = parseField(fields[i], s.lo, s.hi, s.names, s.namesLo)
		if err != nil {
			return nil, fmt.Errorf("cron %q: %s: %v", expr, s.name, err)
		}
		if s.star != nil {
			*s.star = strings.HasPrefix(fields[i], "*")
		}
		if s.sundayHi && *s.dst&(1<<7) != 0 {
			*s.dst |= 1
			*s.dst &^= 1 << 7
		}
	}
	return c, nil
}

func parseField(f string, lo, hi int, names []string, namesLo int) (uint64, error) {
	var bits uint64
	for _, item := range strings.Split(f, ",") {
		if item == "" {
			return 0, fmt.Errorf("empty list item in %q", f)
		}
		rng, step := item, 1
		if r, st, ok := strings.Cut(item, "/"); ok {
			rng = r
			n, err := strconv.Atoi(st)
			if err != nil || n < 1 {
				return 0, fmt.Errorf("bad step %q", st)
			}
			step = n
		}
		var a, b int
		switch {
		case rng == "*":
			a, b = lo, hi
			if hi == 7 {
				b = 6 // * in day-of-week is 0-6; 7 is only Sunday's alias
			}
		case strings.Contains(rng, "-"):
			x, y, _ := strings.Cut(rng, "-")
			var err error
			if a, err = parseValue(x, lo, hi, names, namesLo); err != nil {
				return 0, err
			}
			if b, err = parseValue(y, lo, hi, names, namesLo); err != nil {
				return 0, err
			}
			if a > b {
				return 0, fmt.Errorf("range %q runs backwards", rng)
			}
		default:
			var err error
			if a, err = parseValue(rng, lo, hi, names, namesLo); err != nil {
				return 0, err
			}
			b = a
			if strings.Contains(item, "/") {
				b = hi // a/n runs from a to the end
			}
		}
		for v := a; v <= b; v += step {
			bits |= 1 << uint(v)
		}
	}
	return bits, nil
}

func parseValue(s string, lo, hi int, names []string, namesLo int) (int, error) {
	if n, err := strconv.Atoi(s); err == nil {
		if n < lo || n > hi {
			return 0, fmt.Errorf("%d out of range %d-%d", n, lo, hi)
		}
		return n, nil
	}
	for i, n := range names {
		if strings.EqualFold(s, n) {
			return i + namesLo, nil
		}
	}
	return 0, fmt.Errorf("bad value %q", s)
}

func (c *Cron) dayMatches(t time.Time) bool {
	dom := c.dom&(1<<uint(t.Day())) != 0
	dow := c.dow&(1<<uint(t.Weekday())) != 0
	if !c.domStar && !c.dowStar {
		return dom || dow
	}
	return dom && dow
}

// repeated reports whether t's wall-clock reading already occurred before it,
// because a DST change set the clock back: a slot there fired the first time.
func repeated(t time.Time) bool {
	start, _ := t.ZoneBounds()
	if start.IsZero() {
		return false
	}
	_, off := t.Zone()
	_, prev := start.Add(-time.Second).In(t.Location()).Zone()
	return prev > off && t.Sub(start) < time.Duration(prev-off)*time.Second
}

// Next implements Schedule. A slot the clock never shows (spring forward)
// does not fire; one it shows twice (fall back) fires once.
func (c *Cron) Next(after time.Time) (time.Time, bool) {
	t := after.In(c.loc).Truncate(time.Minute).Add(time.Minute)
	limit := after.AddDate(horizon, 0, 0)
	for t.Before(limit) {
		switch {
		case c.month&(1<<uint(t.Month())) == 0:
			t = time.Date(t.Year(), t.Month()+1, 1, 0, 0, 0, 0, c.loc)
		case !c.dayMatches(t):
			t = time.Date(t.Year(), t.Month(), t.Day()+1, 0, 0, 0, 0, c.loc)
		case c.hour&(1<<uint(t.Hour())) == 0:
			t = t.Add(time.Duration(60-t.Minute()) * time.Minute)
		case c.minute&(1<<uint(t.Minute())) == 0:
			t = t.Add(time.Minute)
		case repeated(t):
			t = t.Add(time.Minute)
		default:
			return t, true
		}
	}
	return time.Time{}, false
}
