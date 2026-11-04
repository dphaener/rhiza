package workflow

import (
	"fmt"
	"regexp"

	"github.com/BurntSushi/toml"
)

// Problem is one thing wrong with a workflow file. Step names the step it
// is in, or is empty for a problem of the file as a whole.
type Problem struct {
	Step    string
	Message string
}

var stepInMessage = regexp.MustCompile(`^step "([^"]+)"`)

func problemOf(step string, err error) Problem {
	msg := err.Error()
	if step == "" {
		if m := stepInMessage.FindStringSubmatch(msg); m != nil {
			step = m[1]
		}
	}
	return Problem{Step: step, Message: msg}
}

// Validate returns the workflow's name (the file's own, else name) and checks a workflow as ParseWith would, but reports every problem
// it can find rather than the first: each unknown key, each step that does
// not build, each target that leads nowhere. Only when those are clean does
// it expand the workflows it includes, as run start does, and report what
// stops that (a missing include, a cycle, a bad skip). No problems means
// ParseWith accepts the workflow.
func Validate(name string, data []byte, resolve Resolver) (string, []Problem) {
	var f fileDef
	md, err := toml.Decode(string(data), &f)
	if err != nil {
		return name, []Problem{{Message: err.Error()}}
	}
	var out []Problem
	for _, k := range md.Undecoded() {
		out = append(out, problemOf("", unknownKey(k)))
	}
	if f.Name != "" {
		name = f.Name
	}
	if name == "" {
		out = append(out, Problem{Message: "workflow has no name: set name = \"...\""})
	}
	if len(f.Step) == 0 {
		return name, append(out, Problem{Message: "workflow has no steps: add a [step.<name>] table"})
	}

	order := stepOrder(md.Keys(), f.Step)
	steps := make(map[string]Step, len(order))
	built := map[string]bool{}
	arbiter := ""
	for _, n := range order {
		s, err := buildStep(n, f.Step[n])
		if err != nil {
			out = append(out, problemOf(n, err))
			// Keep the step so a target naming it is not also reported.
			steps[n] = Step{Name: n, Kind: f.Step[n].Kind}
			continue
		}
		built[n] = true
		if s.Kind == KindArbiter {
			if arbiter != "" {
				out = append(out, problemOf(n, fmt.Errorf("steps %q and %q are both arbiter steps; a workflow has at most one", arbiter, n)))
			}
			arbiter = n
		}
		steps[n] = s
	}
	for _, n := range order {
		if !built[n] {
			continue
		}
		if err := checkTargets(steps[n], steps); err != nil {
			out = append(out, problemOf(n, err))
		}
	}
	if len(out) > 0 {
		return name, out
	}
	if _, err := ParseWith(name, data, resolve); err != nil {
		return name, []Problem{problemOf("", err)}
	}
	return name, nil
}
