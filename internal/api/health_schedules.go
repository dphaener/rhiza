package api

// This file holds the schedules part of health: one line per [[schedules]]
// entry of a project with its last run, result and next run. The part is
// warn, never bad, and absent for a project with no schedules.

import (
	"fmt"
	"os"
	"time"

	"github.com/dphaener/rhiza/internal/dispatch"
	"github.com/dphaener/rhiza/internal/run"
)

// healthSchedules is the schedules part of the project ctx works in, and
// false when the project has none (or its file does not read: the check part
// reports that). daemonRunning says whether a daemon is there to fire them.
func healthSchedules(ctx *Ctx, daemonRunning bool) (HealthPart, bool) {
	rules, err := dispatch.Load(ctx.Root)
	if err != nil || len(rules.Schedules) == 0 {
		return HealthPart{}, false
	}
	var runs []run.Journal
	if common, err := gitCommonDir(ctx); err == nil {
		runs, _ = run.NewStore(run.Dir(common)).List()
	}
	host, _ := os.Hostname()
	now := ctx.Now()
	p := HealthPart{Part: HealthSchedules, State: HealthOK}
	var problems []string
	for _, sc := range rules.Schedules {
		var last *run.Journal
		var ls run.State
		for i := range runs {
			s := runs[i].State()
			if s.Trigger != run.TriggerSchedule || s.TriggerName != sc.Name {
				continue
			}
			if last == nil || s.StartedAt.After(ls.StartedAt) {
				last, ls = &runs[i], s
			}
		}
		line := sc.Name + ": last run "
		switch {
		case last == nil:
			line += "none yet"
		default:
			result := ls.Result
			switch {
			case !ls.Ended && ls.Liveness(host, nil) == run.LiveRunning:
				result = "running"
			case !ls.Ended:
				result = "interrupted"
			}
			line += fmt.Sprintf("%s (run %s), %s", ls.StartedAt.UTC().Format(time.RFC3339), last.ID, result)
			switch {
			case ls.Result == run.Failed:
				problems = append(problems, sc.Name+" failed")
				if p.Next == "" {
					p.Next = "rhiza run log " + last.ID
				}
			case len(ls.Warnings) > 0:
				problems = append(problems, sc.Name+" changed the main checkout")
				line += ": " + ls.Warnings[0]
				if p.Next == "" {
					p.Next = "rhiza run status " + last.ID
				}
			}
		}
		next := "none"
		if sch, err := scheduleOf(sc); err == nil {
			if t, ok := sch.Next(now); ok {
				next = t.UTC().Format(time.RFC3339)
			}
		}
		p.Detail = append(p.Detail, line+"; next run "+next)
	}
	n := len(rules.Schedules)
	p.Summary = fmt.Sprintf("%d schedule", n)
	if n != 1 {
		p.Summary += "s"
	}
	if len(problems) > 0 {
		p.State = HealthWarn
		p.Summary += "; "
		for i, s := range problems {
			if i > 0 {
				p.Summary += ", "
			}
			p.Summary += s
		}
	}
	if !daemonRunning {
		p.State = HealthWarn
		p.Summary += "; no daemon is running to fire them"
		if p.Next == "" {
			p.Next = "rhiza daemon --autonomous"
		}
	}
	return p, true
}
