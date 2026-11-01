package api

import (
	"strings"
	"testing"
)

func TestInboxNoticeOpen(t *testing.T) {
	got := inboxNotice("Plan reports", InboxItem{Kind: InboxOpen})
	if got != "Plan reports: a plan run stopped for open questions" || strings.Contains(got, "gave up") {
		t.Fatalf("notice: %q", got)
	}
}

func TestInboxNoticeOverlap(t *testing.T) {
	got := inboxNotice("Edit one", InboxItem{Kind: InboxOverlap, Other: "CSB"})
	if got != "Edit one: an overlap with change set CSB waits on you" {
		t.Fatalf("notice: %q", got)
	}
}
