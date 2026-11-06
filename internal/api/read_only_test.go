package api

import (
	"slices"
	"testing"
)

// The operations the web UI uses are marked ReadOnly, and nothing that
// writes the log, settings, files, Rex or GitHub is.
func TestReadOnlyMarker(t *testing.T) {
	uses := []string{"run.ls", "run.summary", "roadmap", "ready", "show", "ls", "tree", "slice",
		"why", "log", "overlaps", "inbox", "daemon.settings", "daemon.status", "status", "search", "run.tail"}
	for _, n := range uses {
		if o, ok := Lookup(n); !ok || !o.ReadOnly {
			t.Errorf("%s is not marked ReadOnly", n)
		}
	}
	for _, o := range Operations() {
		if o.ReadOnly && !slices.Contains(uses, o.Name) {
			t.Errorf("%s is marked ReadOnly but is not one the web UI uses", o.Name)
		}
		if o.ReadOnly && (o.Write || o.Record) {
			t.Errorf("%s is ReadOnly and writes the log", o.Name)
		}
	}
	for _, n := range []string{"land", "accept", "daemon.set", "daemon.start", "daemon.stop", "daemon.install",
		"daemon.uninstall", "daemon", "render", "ship.merge", "ship.land", "ship.open-pr", "ship.close-out",
		"ship.rebase", "ship.command", "add", "set", "mv", "retire", "link", "unlink", "apply", "init",
		"change.new", "change.use", "guide.install", "rex.install", "rex.uninstall", "ci.install", "hook.run",
		"inbox.open", "inbox.approve", "inbox.request-changes", "run.start", "run.log", "run.resume",
		"run.abandon", "session.start", "session.run", "session.hook", "session.done"} {
		if o, ok := Lookup(n); ok && o.ReadOnly {
			t.Errorf("%s is marked ReadOnly", n)
		}
	}
	for _, o := range Operations() {
		if o.Write && o.ReadOnly {
			t.Errorf("%s writes the log and is ReadOnly", o.Name)
		}
	}
}
