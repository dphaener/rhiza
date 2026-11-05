package api

import (
	"strings"
	"testing"
)

func TestInboxNoticeNamesNoHandleOrID(t *testing.T) {
	gate := []InboxGateQuestion{{Handle: "q-handle", Statement: "Full statement?", Restatement: "Plain words?"}}
	cases := []struct {
		name string
		it   InboxItem
		want string
	}{
		{"question", InboxItem{Kind: InboxQuestion, ID: "abc1234", Change: "01CHANGEID", Question: "Which way?"}, "T: a question waits on you: Which way?"},
		{"open", InboxItem{Kind: InboxOpen, ID: "abc1234", Change: "01CHANGEID", Questions: gate}, "T: a plan run stopped for open questions, the first: Plain words?"},
		{"open statement", InboxItem{Kind: InboxOpen, ID: "abc1234", Change: "01CHANGEID", Questions: []InboxGateQuestion{{Handle: "q-handle", Statement: "Full statement?"}}}, "T: a plan run stopped for open questions, the first: Full statement?"},
		{"open none", InboxItem{Kind: InboxOpen, ID: "abc1234"}, "T: a plan run stopped for open questions"},
		{"overlap", InboxItem{Kind: InboxOverlap, ID: "abc1234", Change: "01CHANGEID", Other: "01OTHERID", OtherTitle: "Other title", Target: "node-handle"}, "T: an overlap with Other title waits on you"},
		{"overlap untitled", InboxItem{Kind: InboxOverlap, ID: "abc1234", Other: "01OTHERID"}, "T: an overlap with another change set waits on you"},
	}
	for _, c := range cases {
		got := inboxNotice("T", c.it)
		if got != c.want {
			t.Errorf("%s: %q, want %q", c.name, got, c.want)
		}
		for _, bad := range []string{"abc1234", "01CHANGEID", "01OTHERID", "q-handle", "node-handle"} {
			if strings.Contains(got, bad) {
				t.Errorf("%s: %q contains %q", c.name, got, bad)
			}
		}
	}
}
