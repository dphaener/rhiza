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

	"github.com/dphaener/rhiza/internal/daemon"
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
	byTest := map[string][]string{}
	for _, r := range runs {
		for _, name := range r.Flakes {
			if _, ok := byTest[name]; !ok {
				tests = append(tests, name)
			}
			where := "head " + r.Head
			if r.PR > 0 {
				where = fmt.Sprintf("PR #%d at head %s", r.PR, r.Head)
			}
			if !slices.Contains(byTest[name], where) {
				byTest[name] = append(byTest[name], where)
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
		open := false
		for _, r := range reqs {
			if r.Kind != daemon.KindBug || !strings.HasPrefix(r.Text, head+":") {
				continue
			}
			if state, _ := daemon.RequestState(r, d.all); state != daemon.RequestPlanned {
				open = true
				break
			}
		}
		if open {
			continue
		}
		text := fmt.Sprintf("%s: failed then passed on a rerun in %s; reproduce it locally, repeatedly, then fix it or quarantine it", head, strings.Join(byTest[name], ", "))
		if err := d.file(text); err != nil {
			d.log("watchdog: bug for flaky test %s not filed: %v", name, err)
			continue
		}
		d.log("watchdog: filed a bug request for flaky test %s", name)
		reqs = append(reqs, daemon.Request{Kind: daemon.KindBug, Text: text})
	}
	return nil
}
