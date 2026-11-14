package workflow

import (
	"regexp"
	"sort"
	"strings"
)

// mergeCommand matches a run that calls rhiza ship merge or gh pr merge,
// as a whole command rather than part of a longer word.
var mergeCommand = regexp.MustCompile(`(^|[^\w-])(rhiza\s+ship\s+merge|gh\s+pr\s+merge)($|[^\w-])`)

// IsMergeStep says the step merges: a command step whose run calls rhiza
// ship merge or gh pr merge.
func IsMergeStep(s Step) bool {
	return s.Kind == KindCommand && mergeCommand.MatchString(s.Run)
}

// gateCommand matches a run that calls rhiza run gate or rhiza ship gate,
// as a whole command, and failedFlag its --failed flag.
var (
	gateCommand = regexp.MustCompile(`(^|[^\w-])rhiza\s+(run|ship)\s+gate($|[^\w-])`)
	failedFlag  = regexp.MustCompile(`(^|\s)--failed($|[\s=])`)
)

// isGateStep says the step gates a merge: a gate step, or a command step
// whose run calls rhiza run gate or rhiza ship gate, not with --failed
// (which re-runs only what failed before).
func isGateStep(s Step) bool {
	if s.Kind == KindGate {
		return true
	}
	return s.Kind == KindCommand && gateCommand.MatchString(s.Run) && !failedFlag.MatchString(s.Run)
}

// UngatedMerge is a merge step the run can reach from its start step along
// the steps' on.<outcome> routes without entering a gate step (isGateStep).
type UngatedMerge struct {
	// Step is the merge step.
	Step string
	// Path is one shortest such path, from the start step to Step, both
	// included.
	Path []string
}

// UngatedMerges is every merge step (IsMergeStep) reachable from the start
// step along a path that enters no gate step, each with one shortest such
// path, in file order. A merge step reached only through a gate step is
// not one.
func (d *Definition) UngatedMerges() []UngatedMerge {
	start, ok := d.Steps[d.Start]
	if !ok || isGateStep(start) {
		return nil
	}
	prev := map[string]string{d.Start: ""}
	queue := []string{d.Start}
	for len(queue) > 0 {
		n := queue[0]
		queue = queue[1:]
		s := d.Steps[n]
		outcomes := make([]string, 0, len(s.On))
		for o := range s.On {
			outcomes = append(outcomes, o)
		}
		sort.Strings(outcomes)
		for _, o := range outcomes {
			t := strings.TrimPrefix(s.On[o], HandbackPrefix)
			if t == TargetDone || t == TargetFailed {
				continue
			}
			next, ok := d.Steps[t]
			if !ok || isGateStep(next) {
				continue
			}
			if _, seen := prev[t]; seen {
				continue
			}
			prev[t] = n
			queue = append(queue, t)
		}
	}
	var out []UngatedMerge
	for _, n := range d.Order {
		s, ok := d.Steps[n]
		if !ok || !IsMergeStep(s) {
			continue
		}
		if _, reached := prev[n]; !reached {
			continue
		}
		var path []string
		for p := n; p != ""; p = prev[p] {
			path = append([]string{p}, path...)
		}
		out = append(out, UngatedMerge{Step: n, Path: path})
	}
	return out
}
