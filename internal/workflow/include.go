package workflow

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// buildWorkflow finishes a workflow step: it needs the name of a workflow,
// takes skip, and refuses every key of a session or a condition, as it
// runs nothing of its own: the included workflow's steps do.
func buildWorkflow(s Step, fs fileStep, fail func(string, ...any) (Step, error)) (Step, error) {
	refused := []struct {
		key string
		set bool
	}{
		{"prompt", fs.Prompt != ""},
		{"agent", fs.Agent != ""},
		{"permission_mode", fs.PermissionMode != ""},
		{"run", fs.Run != ""},
		{"until", fs.Until != ""},
		{"until_when", fs.UntilWhen != ""},
		{"timeout", fs.Timeout != ""},
		{"stall", fs.Stall != ""},
		{"interval", fs.Interval != ""},
		{"max_visits", fs.MaxVisits != 0},
	}
	for _, r := range refused {
		if r.set {
			return fail("a workflow step cannot have %s: it runs the steps of another workflow, which have their own; remove it", r.key)
		}
	}
	if strings.TrimSpace(fs.Workflow) == "" {
		return fail("a workflow step needs workflow = \"<name of a workflow>\"")
	}
	if !namePattern.MatchString(fs.Workflow) {
		return fail("workflow %q: want a workflow's name, like \"build\"", fs.Workflow)
	}
	s.Workflow = fs.Workflow
	s.Skip = fs.Skip
	return s, nil
}

// inlined is a workflow step's workflow, parsed and with its skipped
// steps taken as passed: its steps' targets name its own flat steps, or
// done or failed.
type inlined struct {
	step  Step
	entry string
	steps map[string]Step
	order []string
}

// expand replaces each workflow step of d with the steps of the workflow it
// names, named <workflow>/<step>, so the engine runs one flat machine. A
// target that names the workflow step goes to the included workflow's
// start; the included workflow ending done takes the step's on.pass, and
// ending failed its on.fail (else the run ends). start is the declared
// step the run begins at.
func (d *Definition) expand(start string, resolve Resolver, stack []string, includes map[string]string) error {
	ins := map[string]*inlined{}
	byWorkflow := map[string]string{}
	for _, n := range d.Order {
		s := d.Steps[n]
		if s.Kind != KindWorkflow {
			continue
		}
		if other, dup := byWorkflow[s.Workflow]; dup {
			return fmt.Errorf("steps %q and %q both run workflow %q; a workflow runs another at most once", other, n, s.Workflow)
		}
		byWorkflow[s.Workflow] = n
		in, err := includeWorkflow(s, resolve, stack, includes)
		if err != nil {
			return err
		}
		ins[n] = in
	}
	if len(ins) == 0 {
		d.Start = start
		d.declared = declaredOf(d, nil)
		return nil
	}

	// entry gives where a target is entered, in the flat steps.
	entries := map[string]string{}
	visiting := map[string]bool{}
	var target func(t string) (string, error)
	var entryOf func(n string) (string, error)
	target = func(t string) (string, error) {
		if t == TargetDone || t == TargetFailed {
			return t, nil
		}
		if w, inner, ok := strings.Cut(t, "/"); ok {
			if in, ok := ins[w]; ok {
				if _, ok := in.steps[inner]; !ok {
					return "", fmt.Errorf("target %q: workflow %q has no step %q left to run", t, in.step.Workflow, inner)
				}
				return in.step.Workflow + "/" + inner, nil
			}
		}
		return entryOf(t)
	}
	entryOf = func(n string) (string, error) {
		in, ok := ins[n]
		if !ok {
			return n, nil
		}
		if e, ok := entries[n]; ok {
			return e, nil
		}
		if visiting[n] {
			return "", fmt.Errorf("step %q: every step of workflow %q is skipped and its on.pass leads back here", n, in.step.Workflow)
		}
		visiting[n] = true
		defer delete(visiting, n)
		e, err := in.outer(in.entry, target)
		if err != nil {
			return "", err
		}
		entries[n] = e
		return e, nil
	}

	steps := make(map[string]Step, len(d.Steps))
	var order []string
	for _, n := range d.Order {
		in, ok := ins[n]
		if !ok {
			s := d.Steps[n]
			on := make(map[string]string, len(s.On))
			for o, t := range s.On {
				ft, err := target(t)
				if err != nil {
					return err
				}
				on[o] = ft
			}
			s.On = on
			steps[n] = s
			order = append(order, n)
			continue
		}
		for _, inner := range in.order {
			s := in.steps[inner]
			s.Name = in.step.Workflow + "/" + inner
			on := make(map[string]string, len(s.On))
			for o, t := range s.On {
				ft, err := in.outer(t, target)
				if err != nil {
					return err
				}
				on[o] = ft
			}
			s.On = on
			steps[s.Name] = s
			order = append(order, s.Name)
		}
	}
	declared := declaredOf(d, nil)
	for n := range d.Steps {
		e, err := entryOf(n)
		if err != nil {
			return err
		}
		pass := TargetDone
		if t, ok := d.Steps[n].On[OutcomePass]; ok {
			if pass, err = target(t); err != nil {
				return err
			}
		}
		declared[n] = declStep{entry: e, pass: pass}
	}
	for _, n := range order {
		s := steps[n]
		if s.Kind != KindSignoff {
			continue
		}
		if ts, ok := steps[s.On[OutcomeChanges]]; !ok || ts.Kind != KindChecks {
			return fmt.Errorf("step %q: on.changes names %q, which is not a checks step; want the checks step that fixes what the owner asks for", n, s.On[OutcomeChanges])
		}
	}
	st, err := entryOf(start)
	if err != nil {
		return err
	}
	if st == TargetDone || st == TargetFailed {
		return fmt.Errorf("start names %q, whose workflow has no step left to run", start)
	}
	d.Steps, d.Order, d.Start, d.declared = steps, order, st, declared
	return nil
}

// declaredOf is d's declared steps, each entered at itself and going on
// to its on.pass target ("done" without one).
func declaredOf(d *Definition, _ map[string]string) map[string]declStep {
	out := make(map[string]declStep, len(d.Steps))
	for n, s := range d.Steps {
		pass := TargetDone
		if t, ok := s.On[OutcomePass]; ok {
			pass = t
		}
		out[n] = declStep{entry: n, pass: pass}
	}
	return out
}

// outer maps a target inside the included workflow to the flat steps of
// the including one: a step to its prefixed name, done to where the
// workflow step goes on pass (else done), failed to where it goes on fail
// (else failed). target resolves a target of the including workflow.
func (in *inlined) outer(t string, target func(string) (string, error)) (string, error) {
	switch t {
	case TargetDone:
		if next, ok := in.step.On[OutcomePass]; ok {
			return target(next)
		}
		return TargetDone, nil
	case TargetFailed:
		if next, ok := in.step.On[OutcomeFail]; ok {
			return target(next)
		}
		return TargetFailed, nil
	}
	return in.step.Workflow + "/" + t, nil
}

// includeWorkflow resolves and parses the workflow a workflow step names
// and takes its skipped steps as passed.
func includeWorkflow(s Step, resolve Resolver, stack []string, includes map[string]string) (*inlined, error) {
	fail := func(format string, args ...any) (*inlined, error) {
		return nil, fmt.Errorf("step %q: %s", s.Name, fmt.Sprintf(format, args...))
	}
	if i := slices.Index(stack, s.Workflow); i >= 0 {
		return fail("workflow %q includes itself: %s", s.Workflow, strings.Join(append(slices.Clone(stack[i:]), s.Workflow), " -> "))
	}
	if resolve == nil {
		return fail("workflow %q cannot be resolved here", s.Workflow)
	}
	data, err := resolve(s.Workflow)
	if err != nil {
		return fail("workflow %q: %v", s.Workflow, err)
	}
	includes[s.Workflow] = string(data)
	inner, err := parse(s.Workflow, data, resolve, append(slices.Clone(stack), s.Workflow), includes)
	if err != nil {
		return fail("workflow %q: %v", s.Workflow, err)
	}

	skipped := map[string]bool{}
	for _, k := range s.Skip {
		if _, ok := inner.declared[k]; !ok {
			var names []string
			for n := range inner.declared {
				names = append(names, n)
			}
			sort.Strings(names)
			return fail("skip names %q, which workflow %q does not have; want any of %s", k, s.Workflow, strings.Join(names, ", "))
		}
		skipped[k] = true
	}
	// A skipped step is taken as passed: whatever went to it goes where
	// it goes on pass.
	skipAt := map[string]string{}
	for k := range skipped {
		skipAt[inner.declared[k].entry] = inner.declared[k].pass
	}
	resolveSkip := func(t string) (string, error) {
		for range len(skipAt) + 1 {
			next, ok := skipAt[t]
			if !ok {
				return t, nil
			}
			t = next
		}
		return "", fmt.Errorf("the skipped steps of workflow %q lead into each other", s.Workflow)
	}

	in := &inlined{step: s, steps: map[string]Step{}}
	if in.entry, err = resolveSkip(inner.Start); err != nil {
		return fail("%v", err)
	}
	all := map[string]Step{}
	for n, is := range inner.Steps {
		if is.Kind == KindArbiter {
			continue
		}
		on := make(map[string]string, len(is.On))
		for _, o := range kindOutcomes[is.Kind] {
			t, ok := is.On[o]
			switch {
			case ok:
			case o == OutcomePass:
				t = TargetDone
			default:
				t = TargetFailed
			}
			if on[o], err = resolveSkip(t); err != nil {
				return fail("%v", err)
			}
		}
		is.On = on
		all[n] = is
	}
	// Keep what the included workflow can still reach.
	var visit func(n string)
	reached := map[string]bool{}
	visit = func(n string) {
		is, ok := all[n]
		if !ok || reached[n] {
			return
		}
		reached[n] = true
		for _, t := range is.On {
			visit(t)
		}
	}
	visit(in.entry)
	for _, n := range inner.Order {
		if reached[n] {
			in.steps[n] = all[n]
			in.order = append(in.order, n)
		}
	}
	return in, nil
}
