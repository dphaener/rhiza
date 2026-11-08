package api

// This file is the watchdog's bug filing (decision watchdog-scope-and-defaults):
// a recovery caused by a bug (a stall, a transient failure other than a CI
// flake, a PR its run left unmerged or without CI) queues one bug request
// with the evidence, unless an open request already names the same recovery
// kind and step. fileBugs = false in the policy turns it off.

import (
	"fmt"
	"strings"

	"github.com/dphaener/rhiza/internal/daemon"
)

// The recovery kinds a bug request names.
const (
	bugStall     = "stall"
	bugTransient = "transient failure"
	bugUnmerged  = "unmerged PR"
	bugNoCI      = "PR without CI"
)

// watchdogBugMaxLog bounds the log lines quoted in a bug request.
const watchdogBugMaxLog = 400

// watchdogBugHead is the start of the text of a bug request about kind at
// step: two recoveries with the same head are the same bug.
func watchdogBugHead(kind, step string) string {
	return fmt.Sprintf("watchdog %s at step %s", kind, step)
}

// watchdogLogLines is the last lines of detail, on one line, bounded.
func watchdogLogLines(detail string) string {
	var lines []string
	for _, l := range strings.Split(detail, "\n") {
		if l = strings.TrimSpace(l); l != "" {
			lines = append(lines, l)
		}
	}
	if len(lines) > 3 {
		lines = lines[len(lines)-3:]
	}
	s := strings.Join(lines, " | ")
	if len(s) > watchdogBugMaxLog {
		s = s[len(s)-watchdogBugMaxLog:]
	}
	return s
}

// fileBug queues the bug request for a recovery of runID at step, unless
// the policy turns it off or an open request already names the same kind
// and step. detail is the log lines behind the recovery. A failure to queue
// is logged: the next poll's recovery tries again.
func (d watchdogDeps) fileBug(pol WatchdogPolicy, runID, kind, step, detail string) {
	if !pol.FileBugs || d.file == nil || d.requests == nil {
		return
	}
	reqs, err := d.requests()
	if err != nil {
		d.log("watchdog: bug for run %s not filed: %v", runID, err)
		return
	}
	head := watchdogBugHead(kind, step)
	for _, r := range reqs {
		if r.Kind != daemon.KindBug || !strings.HasPrefix(r.Text, head+":") {
			continue
		}
		if state, _ := daemon.RequestState(r, d.all); state != daemon.RequestPlanned {
			return
		}
	}
	var rids []string
	for _, r := range reqs {
		if r.Run == runID {
			rids = append(rids, r.ID)
		}
	}
	text := fmt.Sprintf("%s: run %s", head, runID)
	if len(rids) > 0 {
		text += fmt.Sprintf(", request %s", strings.Join(rids, ", "))
	}
	if l := watchdogLogLines(detail); l != "" {
		text += "; log: " + l
	}
	if err := d.file(text); err != nil {
		d.log("watchdog: bug for run %s not filed: %v", runID, err)
		return
	}
	d.log("watchdog: run %s step %s: filed a bug request for the %s", runID, step, kind)
}
