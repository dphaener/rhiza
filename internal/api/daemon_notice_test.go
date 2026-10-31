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
