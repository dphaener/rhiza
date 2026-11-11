package api

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
)

// TestInboxNoticeNamesNoHandleOrID has a case per inbox kind, each item
// carrying IDs and handles in every field it has: no notice holds any of
// them, and change sets and questions show by title and text.
func TestInboxNoticeNamesNoHandleOrID(t *testing.T) {
	gate := []InboxGateQuestion{{Handle: "q-handle", Statement: "Full statement?", Restatement: "Plain words?", Blocks: "node-handle"}}
	// ids fills the fields that carry IDs and handles on an item of kind.
	ids := func(kind string, it InboxItem) InboxItem {
		it.Kind = kind
		it.ID = "abc1234"
		it.Run = "01RUNID"
		it.Change = "01CHANGEID"
		it.Step = "build/build"
		it.Session = "sess-uuid"
		it.Tab = "tab-title"
		it.TabBlock = "block:blk"
		it.TabSession = "session:rexsess"
		it.Handle = "q-handle"
		it.Blocks = "node-handle"
		it.Other = "01OTHERID"
		it.Target = "node-handle"
		it.Next = "rhiza show q-handle"
		if it.Detail == "" {
			it.Detail = "leaf node-handle in 01RUNID"
		}
		return it
	}
	cases := []struct {
		name string
		it   InboxItem
		want string
	}{
		{"question", ids(InboxQuestion, InboxItem{Question: "Which way?"}), "T: a question waits on you: Which way?"},
		{"open", ids(InboxOpen, InboxItem{Questions: gate}), "T: a plan run stopped for open questions, the first: Plain words?"},
		{"open statement", ids(InboxOpen, InboxItem{Questions: []InboxGateQuestion{{Handle: "q-handle", Statement: "Full statement?"}}}), "T: a plan run stopped for open questions, the first: Full statement?"},
		{"open none", ids(InboxOpen, InboxItem{}), "T: a plan run stopped for open questions"},
		{"signoff", ids(InboxSignoff, InboxItem{}), "T: a sign-off waits on you"},
		{"failed", ids(InboxFailed, InboxItem{}), "T: a run failed"},
		{"gave-up", ids(InboxGaveUp, InboxItem{}), "T: the daemon gave up on a run"},
		{"skipped", ids(InboxSkipped, InboxItem{}), "T: an interrupted run was not resumed"},
		{"blocker", ids(InboxBlocker, InboxItem{Question: "Which store?", BlocksTitle: "Export orders"}), "T: a question blocks it: Which store?"},
		{"overlap", ids(InboxOverlap, InboxItem{OtherTitle: "Other title"}), "T: an overlap with Other title waits on you"},
		{"overlap untitled", ids(InboxOverlap, InboxItem{}), "T: an overlap with another change set waits on you"},
		{"stalled", ids(InboxStalled, InboxItem{}), "T: a session stalled"},
		{"disk-low", ids(InboxDiskLow, InboxItem{}), "disk space is low: the daemon holds new starts and resumes"},
		{"load-high", ids(InboxLoadHigh, InboxItem{}), "machine load is high: the daemon holds new starts and resumes"},
		{"memory-low", ids(InboxMemoryLow, InboxItem{}), "available memory is low: the daemon holds new starts and resumes"},
		{"usage-reserve", ids(InboxUsageReserve, InboxItem{}), "usage reserve: the daemon holds new launches until the window resets"},
		{"main-stale", ids(InboxMainStale, InboxItem{Detail: "git pull failed"}), "T: main was not updated after the merge: git pull failed"},
		{"accept", ids(InboxAccept, InboxItem{}), "T: a plan waits for your acceptance"},
		{"interview", ids(InboxInterview, InboxItem{}), "T: a milestone interview waits on you"},
	}
	kinds := map[string]bool{}
	for _, c := range cases {
		kinds[c.it.Kind] = true
		got := inboxNotice("T", c.it)
		if got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
		for _, bad := range []string{"abc1234", "01RUNID", "01CHANGEID", "01OTHERID", "q-handle", "node-handle", "sess-uuid", "blk", "rexsess"} {
			if strings.Contains(got, bad) {
				t.Errorf("%s: %q contains %q", c.name, got, bad)
			}
		}
	}
	for _, k := range []string{InboxQuestion, InboxOpen, InboxSignoff, InboxFailed, InboxGaveUp, InboxSkipped, InboxBlocker, InboxOverlap, InboxStalled, InboxDiskLow, InboxUsageReserve, InboxLoadHigh, InboxMemoryLow, InboxMainStale, InboxAccept, InboxInterview} {
		if !kinds[k] {
			t.Errorf("no case for kind %s", k)
		}
	}
}

// A load-high or memory-low hold that has lasted 31 minutes notifies the
// owner once, however many polls see it.
func TestDaemonNotifyMachineHold(t *testing.T) {
	for _, tc := range []struct{ reason, msg string }{
		{daemon.SkipLoadHigh, "machine load is high: the daemon holds new starts and resumes"},
		{daemon.SkipMemoryLow, "available memory is low: the daemon holds new starts and resumes"},
	} {
		t.Run(tc.reason, func(t *testing.T) {
			ctx := testCtx(t.TempDir())
			var rex bytes.Buffer
			ctx.Rex = &rex
			rs := run.NewStore(t.TempDir())
			state := daemon.State{MachineHold: &daemon.MachineHold{
				Since: ctx.Now().Add(-31 * time.Minute), Reason: tc.reason, Detail: "reading over limit",
			}}
			for range 2 {
				if err := daemonNotify(ctx, rs, &state, ""); err != nil {
					t.Fatal(err)
				}
			}
			var blocked []map[string]string
			for _, r := range rexReports(t, rex.String()) {
				if r["state"] == "blocked" {
					blocked = append(blocked, r)
				}
			}
			if len(blocked) != 1 || blocked[0]["msg"] != tc.msg {
				t.Fatalf("blocked reports %v, want one %q", blocked, tc.msg)
			}
		})
	}
}
