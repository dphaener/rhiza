// Package dispatch reads a repo's .rhiza/dispatch.toml: the rules mapping a
// trigger and a request kind to the workflow a run starts.
package dispatch

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/BurntSushi/toml"

	"github.com/dphaener/rhiza/internal/schedule"
)

// File is the rules file's path under a repo root.
const File = ".rhiza/dispatch.toml"

// The kinds a rule maps, and the tables that hold the rules.
const (
	KindFeature   = "feature"
	KindBug       = "bug"
	KindMilestone = "milestone"

	TableRequest = "request"
	TableReady   = "ready"

	// TableSchedules is the array of tables holding the scheduled triggers.
	TableSchedules = "schedules"
)

// The values catch_up and overlap take; the first of each is the default.
const (
	CatchUpOnce  = "once"
	CatchUpSkip  = "skip"
	OverlapSkip  = "skip"
	OverlapQueue = "queue"
)

// Schedule is one [[schedules]] entry: Every or Cron says when, Workflow or
// Command says what.
type Schedule struct {
	Name     string
	Every    string
	Cron     string
	Workflow string
	Command  string
	CatchUp  string
	Overlap  string
}

// The built-in defaults: today's behaviour.
const (
	DefaultRequest = "plan"
	DefaultReady   = "ship-change"
)

// Kinds lists the kinds [request] may name.
var Kinds = []string{KindFeature, KindBug, KindMilestone}

// ReadyKinds lists the kinds [ready] may name. A milestone never becomes a
// ready change set, so a [ready] milestone rule would never fire.
var ReadyKinds = []string{KindFeature, KindBug}

// TableKinds is the kinds table may name.
func TableKinds(table string) []string {
	if table == TableReady {
		return ReadyKinds
	}
	return Kinds
}

// Rules maps each kind to a workflow name, per trigger: Request has an entry
// for every kind in Kinds, Ready for every kind in ReadyKinds.
type Rules struct {
	Request map[string]string
	Ready   map[string]string
	// Set holds "<table>.<kind>" for each rule the file sets, as against
	// the built-in default.
	Set map[string]bool
	// Schedules are the file's [[schedules]] entries, in file order.
	Schedules []Schedule
}

// RequestPick is ForRequest and where it came from: the rule's name
// (request.bug) or "default" when the file does not set it.
func (r Rules) RequestPick(kind string) (name, source string) {
	return r.Request[kind], r.source(TableRequest, kind)
}

// ReadyPick is ForReady and where it came from.
func (r Rules) ReadyPick(kind string) (name, source string) {
	return r.ForReady(kind), r.source(TableReady, kind)
}

func (r Rules) source(table, kind string) string {
	if r.Set[table+"."+kind] {
		return "rule " + table + "." + kind
	}
	return "default"
}

// Defaults returns the built-in rules.
func Defaults() Rules {
	r := Rules{Request: map[string]string{}, Ready: map[string]string{}, Set: map[string]bool{}}
	for _, k := range Kinds {
		r.Request[k] = DefaultRequest
	}
	for _, k := range ReadyKinds {
		r.Ready[k] = DefaultReady
	}
	return r
}

// ForRequest is the workflow for a queued request of kind.
func (r Rules) ForRequest(kind string) string { return r.Request[kind] }

// ForReady is the workflow for an accepted, ready change set of kind. A kind
// [ready] cannot name gets the default.
func (r Rules) ForReady(kind string) string {
	if name, ok := r.Ready[kind]; ok {
		return name
	}
	return DefaultReady
}

// Load reads root's dispatch.toml. A missing file gives the defaults.
func Load(root string) (Rules, error) {
	data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(File)))
	if errors.Is(err, os.ErrNotExist) {
		return Defaults(), nil
	}
	if err != nil {
		return Rules{}, fmt.Errorf("read %s: %w", File, err)
	}
	r, err := Parse(data)
	if err != nil {
		return Rules{}, fmt.Errorf("%s: %w", File, err)
	}
	return r, nil
}

// Parse reads the file's content. Anything not [request], [ready] or
// [[schedules]], any key
// that is not a kind, and any value that is not a bare workflow name is an
// error naming the table and key; [ready] takes no milestone key.
func Parse(data []byte) (Rules, error) {
	var raw map[string]any
	if _, err := toml.Decode(string(data), &raw); err != nil {
		return Rules{}, err
	}
	r := Defaults()
	names := make([]string, 0, len(raw))
	for t := range raw {
		names = append(names, t)
	}
	sort.Strings(names)
	for _, t := range names {
		var dst map[string]string
		switch t {
		case TableSchedules:
			ss, err := parseSchedules(raw[t])
			if err != nil {
				return Rules{}, err
			}
			r.Schedules = ss
			continue
		case TableRequest:
			dst = r.Request
		case TableReady:
			dst = r.Ready
		default:
			return Rules{}, fmt.Errorf("unknown table [%s]: want [%s], [%s] or [[%s]]", t, TableRequest, TableReady, TableSchedules)
		}
		tbl, ok := raw[t].(map[string]any)
		if !ok {
			return Rules{}, fmt.Errorf("[%s] must be a table", t)
		}
		keys := make([]string, 0, len(tbl))
		for k := range tbl {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if !isKind(t, k) {
				return Rules{}, fmt.Errorf("[%s] unknown key %q: want one of %s", t, k, strings.Join(TableKinds(t), ", "))
			}
			s, ok := tbl[k].(string)
			if !ok {
				return Rules{}, fmt.Errorf("[%s] %s: want a workflow name as a string, got %T", t, k, tbl[k])
			}
			if !bareName(s) {
				return Rules{}, fmt.Errorf("[%s] %s: %q is not a bare workflow name (no path, no .toml)", t, k, s)
			}
			dst[k] = s
			r.Set[t+"."+k] = true
		}
	}
	return r, nil
}

func isKind(table, k string) bool {
	for _, x := range TableKinds(table) {
		if x == k {
			return true
		}
	}
	return false
}

func bareName(s string) bool {
	return s != "" && !strings.ContainsAny(s, "/\\ \t\n") && !strings.HasSuffix(s, ".toml") && s != "." && s != ".."
}

var scheduleKeys = []string{"name", "every", "cron", "workflow", "command", "catch_up", "overlap"}

// parseSchedules reads the [[schedules]] array. Every error names the
// schedule (by name, else by position) and the key.
func parseSchedules(v any) ([]Schedule, error) {
	entries, ok := v.([]map[string]any)
	if !ok {
		return nil, fmt.Errorf("[%s] must be an array of tables: write [[%s]]", TableSchedules, TableSchedules)
	}
	var out []Schedule
	seen := map[string]bool{}
	for i, e := range entries {
		label := fmt.Sprintf("schedule #%d", i+1)
		str := func(k string) (string, bool, error) {
			x, ok := e[k]
			if !ok {
				return "", false, nil
			}
			s, ok := x.(string)
			if !ok {
				return "", true, fmt.Errorf("[[%s]] %s: %s: want a string, got %T", TableSchedules, label, k, x)
			}
			return s, true, nil
		}
		name, has, err := str("name")
		if err != nil {
			return nil, err
		}
		if !has || strings.TrimSpace(name) == "" {
			return nil, fmt.Errorf("[[%s]] %s: name: required", TableSchedules, label)
		}
		label = fmt.Sprintf("schedule %q", name)
		if seen[name] {
			return nil, fmt.Errorf("[[%s]] %s: name: duplicate, a name is unique in the repo", TableSchedules, label)
		}
		seen[name] = true
		keys := make([]string, 0, len(e))
		for k := range e {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			known := false
			for _, w := range scheduleKeys {
				known = known || w == k
			}
			if !known {
				return nil, fmt.Errorf("[[%s]] %s: unknown key %q: want one of %s", TableSchedules, label, k, strings.Join(scheduleKeys, ", "))
			}
		}
		sc := Schedule{Name: name, CatchUp: CatchUpOnce, Overlap: OverlapSkip}
		vals := map[string]*string{"every": &sc.Every, "cron": &sc.Cron, "workflow": &sc.Workflow, "command": &sc.Command}
		present := map[string]bool{}
		for _, k := range []string{"every", "cron", "workflow", "command"} {
			s, has, err := str(k)
			if err != nil {
				return nil, err
			}
			if has && strings.TrimSpace(s) == "" {
				return nil, fmt.Errorf("[[%s]] %s: %s: must not be empty", TableSchedules, label, k)
			}
			*vals[k], present[k] = s, has
		}
		switch {
		case present["every"] && present["cron"]:
			return nil, fmt.Errorf("[[%s]] %s: every and cron: set exactly one, not both", TableSchedules, label)
		case !present["every"] && !present["cron"]:
			return nil, fmt.Errorf("[[%s]] %s: every, cron: set exactly one, got neither", TableSchedules, label)
		case present["workflow"] && present["command"]:
			return nil, fmt.Errorf("[[%s]] %s: workflow and command: set exactly one, not both", TableSchedules, label)
		case !present["workflow"] && !present["command"]:
			return nil, fmt.Errorf("[[%s]] %s: workflow, command: set exactly one, got neither", TableSchedules, label)
		}
		if present["every"] {
			if _, err := schedule.ParseEvery(sc.Every); err != nil {
				return nil, fmt.Errorf("[[%s]] %s: every: %v", TableSchedules, label, err)
			}
		} else if _, err := schedule.ParseCron(sc.Cron, nil); err != nil {
			return nil, fmt.Errorf("[[%s]] %s: cron: %v", TableSchedules, label, err)
		}
		if present["workflow"] && !bareName(sc.Workflow) {
			return nil, fmt.Errorf("[[%s]] %s: workflow: %q is not a bare workflow name (no path, no .toml)", TableSchedules, label, sc.Workflow)
		}
		for _, o := range []struct {
			key  string
			dst  *string
			want []string
		}{{"catch_up", &sc.CatchUp, []string{CatchUpOnce, CatchUpSkip}}, {"overlap", &sc.Overlap, []string{OverlapSkip, OverlapQueue}}} {
			s, has, err := str(o.key)
			if err != nil {
				return nil, err
			}
			if !has {
				continue
			}
			ok := false
			for _, w := range o.want {
				ok = ok || w == s
			}
			if !ok {
				return nil, fmt.Errorf("[[%s]] %s: %s: %q is not one of %s", TableSchedules, label, o.key, s, strings.Join(o.want, ", "))
			}
			*o.dst = s
		}
		out = append(out, sc)
	}
	return out, nil
}
