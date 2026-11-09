package api

// This file files a bug for each flaky test (decision ci-unstable-gate-rules):
// a test the CI stability record holds a flake for queues one bug request
// naming the test, its flaky runs and their heads, asking to reproduce it
// locally and fix or quarantine it. An open request naming the same test
// queues no second one. fileBugs = false in the policy turns it off.

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// flakyBugHead is the start of the text of a bug request about test: two
// flakes of one test are the same bug.
func flakyBugHead(test string) string {
	return "flaky test " + test
}

// watchdogFlakyTests queues the bug request for each test the record holds
// a flake for, unless the policy turns bug filing off or an open request
// already names the test. A failure to queue is logged: the next poll
// tries again.
func watchdogFlakyTests(pol WatchdogPolicy, d watchdogDeps) error {
	if !pol.FileBugs || d.file == nil || d.requests == nil || d.flakes == nil {
		return nil
	}
	runs, err := d.flakes()
	if err != nil {
		return err
	}
	var tests []string
	byTest := map[string][]CIRun{}
	for _, r := range runs {
		for _, name := range r.Flakes {
			if _, ok := byTest[name]; !ok {
				tests = append(tests, name)
			}
			if !slices.ContainsFunc(byTest[name], func(o CIRun) bool { return flakyWhere(o) == flakyWhere(r) }) {
				byTest[name] = append(byTest[name], r)
			}
		}
	}
	if len(tests) == 0 {
		return nil
	}
	reqs, err := d.requests()
	if err != nil {
		return err
	}
	for _, name := range tests {
		head := flakyBugHead(name)
		if flakyBugFiled(head, byTest[name], reqs, d.all) {
			continue
		}
		var wheres []string
		for _, r := range byTest[name] {
			wheres = append(wheres, flakyWhere(r))
		}
		text := fmt.Sprintf("%s: failed then passed on a rerun in %s; reproduce it locally, repeatedly, then fix it or quarantine it", head, strings.Join(wheres, ", "))
		if err := d.file(text); err != nil {
			d.log("watchdog: bug for flaky test %s not filed: %v", name, err)
			continue
		}
		d.log("watchdog: filed a bug request for flaky test %s", name)
		reqs = append(reqs, daemon.Request{Kind: daemon.KindBug, Text: text})
	}
	return nil
}

// flakyBugFiled is true when a bug request already covers every flaky
// run of the test: one is still open (queued or planning), or each run is
// named in a closed (planned or failed) request or was recorded before
// such a request's run finished. A later flake files the test again.
func flakyBugFiled(head string, runs []CIRun, reqs []daemon.Request, all []run.Journal) bool {
	var texts []string
	var closed time.Time
	for _, r := range reqs {
		if r.Kind != daemon.KindBug || !strings.HasPrefix(r.Text, head+":") {
			continue
		}
		if state, _ := daemon.RequestState(r, all); state == daemon.RequestQueued || state == daemon.RequestPlanning {
			return true
		}
		texts = append(texts, r.Text)
		if at := runFinishedAt(r.Run, all); at.After(closed) {
			closed = at
		}
	}
	if len(texts) == 0 {
		return false
	}
	for _, cr := range runs {
		where := flakyWhere(cr)
		if slices.ContainsFunc(texts, func(t string) bool {
			return strings.Contains(t, where+",") || strings.Contains(t, where+";")
		}) {
			continue
		}
		if at, err := time.Parse(time.RFC3339, cr.Time); err == nil && !closed.IsZero() && !at.After(closed) {
			continue
		}
		return false
	}
	return true
}

// runFinishedAt is when the run's journal says it finished; zero if it
// has not or is not in all.
func runFinishedAt(id string, all []run.Journal) (at time.Time) {
	for _, j := range all {
		if j.ID != id {
			continue
		}
		for _, e := range j.Entries {
			if e.Type == run.TypeFinished {
				at = e.At
			}
		}
	}
	return at
}

// flakyWhere names a flaky run in a bug request's text.
func flakyWhere(r CIRun) string {
	if r.PR > 0 {
		return fmt.Sprintf("PR #%d at head %s", r.PR, r.Head)
	}
	return "head " + r.Head
}
