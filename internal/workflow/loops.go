package workflow

import (
	"sort"
	"strings"
)

// PlainRedLoop is a red outcome that loops back with no cap: a step's
// on.<outcome> other than pass names, by plain step name rather than
// handback:<agent step>, a step that leads back to the routing step
// through an agent step. handback_cap does not count it, no handback.md is
// delivered and no inbox item parks it; only max_visits ends it.
type PlainRedLoop struct {
	// Step is the routing step and Outcome its red outcome.
	Step, Outcome string
	// Target is the plain step name the outcome names.
	Target string
	// Agent is the first agent step on the way back: the one a
	// handback:<agent step> route would hand back to.
	Agent string
}

// PlainRedLoops is every plain red loop of the definition, in file order
// of the routing step, then by outcome.
func (d *Definition) PlainRedLoops() []PlainRedLoop {
	var out []PlainRedLoop
	for _, n := range d.Order {
		s, ok := d.Steps[n]
		if !ok {
			continue
		}
		outcomes := make([]string, 0, len(s.On))
		for o := range s.On {
			if o != OutcomePass {
				outcomes = append(outcomes, o)
			}
		}
		sort.Strings(outcomes)
		for _, o := range outcomes {
			t := s.On[o]
			if t == TargetDone || t == TargetFailed || strings.HasPrefix(t, HandbackPrefix) {
				continue
			}
			if agent, ok := d.backThroughAgent(t, n); ok {
				out = append(out, PlainRedLoop{Step: n, Outcome: o, Target: t, Agent: agent})
			}
		}
	}
	return out
}

// backThroughAgent finds a path along the steps' on.<outcome> routes from
// from to to that passes through an agent step (from and to included) and
// takes no handback:<step> route, and returns the first agent step on it.
func (d *Definition) backThroughAgent(from, to string) (string, bool) {
	type state struct{ step, agent string }
	seen := map[state]bool{}
	var walk func(step, agent string) (string, bool)
	walk = func(step, agent string) (string, bool) {
		s, ok := d.Steps[step]
		if !ok {
			return "", false
		}
		if agent == "" && s.Kind == KindAgent {
			agent = step
		}
		if step == to {
			return agent, agent != ""
		}
		// Only whether an agent was passed changes what lies ahead.
		k := state{step, agent}
		if agent != "" {
			k.agent = "*"
		}
		if seen[k] {
			return "", false
		}
		seen[k] = true
		outcomes := make([]string, 0, len(s.On))
		for o := range s.On {
			outcomes = append(outcomes, o)
		}
		sort.Strings(outcomes)
		for _, o := range outcomes {
			t := s.On[o]
			// A path through a handback route is counted by handback_cap.
			if t == TargetDone || t == TargetFailed || strings.HasPrefix(t, HandbackPrefix) {
				continue
			}
			if a, ok := walk(t, agent); ok {
				return a, true
			}
		}
		return "", false
	}
	return walk(from, "")
}
