package api

// This file holds the daemon's part of scheduled triggers: on each poll it
// reads the project's [[schedules]], decides (daemon.ScheduleDue) which are
// due and starts each one's run as a detached rhiza run start --schedule
// process. A scheduled run opens no Rex tab or pane of its own, and its
// start, end and failure go to the daemon log (stderr), never to a tab.

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/dispatch"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/schedule"
)

// daemonSchedules is what the daemon keeps for one project's schedules
// between polls. Everything else it needs it reads from the runs' journals.
type daemonSchedules struct {
	// since is when the daemon first saw each schedule's entry.
	since map[string]time.Time
	// starting holds a schedule's run started and not yet journaled.
	starting map[string]time.Time
	// problems is the last problem logged per schedule ("" for the file),
	// so a standing one is logged once, not on every poll.
	problems map[string]string
	// ended holds the runs whose end was logged; nil until the first look,
	// which marks every run already ended as logged.
	ended map[string]bool
}

// dueSchedule is a schedule that fires now.
type dueSchedule struct {
	sc   dispatch.Schedule
	slot time.Time
}

func (d *daemonSchedules) problem(ctx *Ctx, project, name, msg string) {
	if d.problems == nil {
		d.problems = map[string]string{}
	}
	if d.problems[name] == msg {
		return
	}
	d.problems[name] = msg
	if msg != "" {
		daemonLogf(ctx, "project %s: %s", project, msg)
	}
}

func daemonLogf(ctx *Ctx, format string, a ...any) {
	if ctx.Stderr != nil {
		fmt.Fprintf(ctx.Stderr, "rhiza: daemon: "+format+"\n", a...)
	}
}

// scheduleOf is the timing of a schedule entry.
func scheduleOf(sc dispatch.Schedule) (schedule.Schedule, error) {
	if sc.Every != "" {
		return schedule.ParseEvery(sc.Every)
	}
	return schedule.ParseCron(sc.Cron, nil)
}

// scheduleStateOf rebuilds the state of schedule name from the runs'
// journals: its latest scheduled run decides LastStart, LastEnd and Running.
// The slot a run fired for is taken as its start, which is never before the
// slot: a slot between the start and now is a later one.
func (d *daemonSchedules) scheduleStateOf(name string, runs []run.Journal, host string, now time.Time) daemon.ScheduleState {
	st := daemon.ScheduleState{Since: d.since[name]}
	var last *run.Journal
	var lastState run.State
	for i := range runs {
		s := runs[i].State()
		if s.Trigger != run.TriggerSchedule || s.TriggerName != name {
			continue
		}
		if last == nil || s.StartedAt.After(lastState.StartedAt) {
			last, lastState = &runs[i], s
		}
	}
	if last != nil {
		st.LastStart, st.LastSlot = lastState.StartedAt, lastState.StartedAt
		switch {
		case lastState.Liveness(host, nil) == run.LiveRunning:
			st.Running = true
		case lastState.Ended:
			for _, e := range last.Entries {
				if e.Type == run.TypeFinished {
					st.LastEnd = e.At
				}
			}
			if st.LastEnd.IsZero() {
				st.LastEnd = st.LastStart
			}
		default:
			// Interrupted: its process is gone; the next slot goes on.
			st.LastEnd = now
		}
	}
	if t, ok := d.starting[name]; ok {
		if now.Sub(t) < daemonStartGrace && (last == nil || !lastState.StartedAt.After(t.Add(-time.Second))) {
			st.Running, st.LastStart, st.LastSlot = true, t, t
		} else {
			delete(d.starting, name)
		}
	}
	return st
}

// due reads the project's schedules and returns those firing at now, in file
// order. A file that does not read, and a schedule whose workflow does not
// resolve, start nothing and are logged.
func (d *daemonSchedules) due(ctx *Ctx, project string, runs []run.Journal) []dueSchedule {
	now := ctx.Now()
	if d.since == nil {
		d.since, d.starting = map[string]time.Time{}, map[string]time.Time{}
	}
	rules, err := dispatch.Load(ctx.Root)
	if err != nil {
		d.problem(ctx, project, "", fmt.Sprintf("schedules not read, none started: %s: %v", dispatch.File, err))
		return nil
	}
	d.problem(ctx, project, "", "")
	host, _ := os.Hostname()
	var out []dueSchedule
	for _, sc := range rules.Schedules {
		if _, ok := d.since[sc.Name]; !ok {
			d.since[sc.Name] = now
		}
		sch, err := scheduleOf(sc)
		if err != nil {
			d.problem(ctx, project, sc.Name, fmt.Sprintf("schedule %s not started: %v", sc.Name, err))
			continue
		}
		dec := daemon.ScheduleDue(sc, sch, d.scheduleStateOf(sc.Name, runs, host, now), now)
		if !dec.Fire {
			d.problem(ctx, project, sc.Name, "")
			continue
		}
		if _, _, _, err := scheduleWorkflow(ctx.Root, sc); err != nil {
			d.problem(ctx, project, sc.Name, fmt.Sprintf("schedule %s not started: %v", sc.Name, err))
			continue
		}
		d.problem(ctx, project, sc.Name, "")
		out = append(out, dueSchedule{sc: sc, slot: dec.Slot})
	}
	return out
}

var unsafeLogName = regexp.MustCompile(`[^A-Za-z0-9._-]+`)

// start detaches the run of a due schedule and logs it. No Rex call is made:
// the child's agent steps, if any, open their sessions in the project's
// Rex session as in any daemon run.
func (d *daemonSchedules) start(ctx *Ctx, project string, ds dueSchedule, mode string) error {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return err
	}
	if !ctx.repoSession {
		// The session's name only: ensureRepoSession would call Rex.
		if name, err := repoSessionName(ctx); err == nil {
			ctx.RexSession = name
		}
	}
	args := []string{"run", "start", "--foreground", "--schedule", ds.sc.Name, "--permission-mode", mode}
	log := "schedule-" + unsafeLogName.ReplaceAllString(ds.sc.Name, "_")
	if err := daemonDetachLog(ctx, common, "a run of schedule "+ds.sc.Name, log, args); err != nil {
		return err
	}
	d.starting[ds.sc.Name] = ctx.Now()
	daemonLogf(ctx, "project %s: started a run of schedule %s (slot %s)", project, ds.sc.Name, ds.slot.UTC().Format(time.RFC3339))
	return nil
}

// logEnds logs the end and result of each scheduled run that ended since the
// last look; a failure's detail comes with it.
func (d *daemonSchedules) logEnds(ctx *Ctx, project string, runs []run.Journal) {
	first := d.ended == nil
	if first {
		d.ended = map[string]bool{}
	}
	for _, j := range runs {
		st := j.State()
		if st.Trigger != run.TriggerSchedule || !st.Ended || d.ended[j.ID] {
			continue
		}
		d.ended[j.ID] = true
		if first {
			continue
		}
		if st.Result == run.Failed {
			daemonLogf(ctx, "project %s: run %s of schedule %s failed: %s", project, j.ID, st.TriggerName, st.Detail)
		} else {
			daemonLogf(ctx, "project %s: run %s of schedule %s %s", project, j.ID, st.TriggerName, st.Result)
		}
	}
}

// startingCount is how many of the schedules' starts are not journaled yet:
// each holds a run slot.
func (d *daemonSchedules) startingCount(now time.Time, runs []run.Journal) int {
	n := 0
	for name, t := range d.starting {
		if now.Sub(t) >= daemonStartGrace {
			continue
		}
		if slices.ContainsFunc(runs, func(j run.Journal) bool {
			s := j.State()
			return s.Trigger == run.TriggerSchedule && s.TriggerName == name && !s.StartedAt.Before(t.Add(-time.Second))
		}) {
			continue
		}
		n++
	}
	return n
}
