package api

import (
	"testing"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/session"
)

// Two bug requests planned at once must not share a branch or worktree: the
// change set a bug's plan run opens takes its display title, and so its slug,
// from the bug's text, not from a bare "Fix".
func TestBugRequestTitlesAreDistinct(t *testing.T) {
	a := daemon.Request{Kind: daemon.KindBug, Text: "Request plan runs fail at close-out: the daemon writes"}
	b := daemon.Request{Kind: daemon.KindBug, Text: "rhiza run resume refuses a plan run whose record change set landed"}
	ta := changeset.Open{Kind: changeset.KindIntention, Intent: a.Intent()}.DisplayTitle()
	tb := changeset.Open{Kind: changeset.KindIntention, Intent: b.Intent()}.DisplayTitle()
	if ta == "Fix" || tb == "Fix" {
		t.Fatalf("display titles %q and %q: a bug's title is a bare Fix", ta, tb)
	}
	if session.Slug(ta) == session.Slug(tb) {
		t.Fatalf("display titles %q and %q share slug %q", ta, tb, session.Slug(ta))
	}
	if ta != "Fix Request plan runs fail at" {
		t.Fatalf("display title %q, want %q", ta, "Fix Request plan runs fail at")
	}
}
