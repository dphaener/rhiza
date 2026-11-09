package api

import (
	"strings"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

// The Overview's inbox reads in the owner's words: an open item is a row
// per question by its full text, and no row carries a run or session ID, a
// question's handle or a command (constraint owner-surfaces-no-ids).
func TestWebInboxViewsOwnerWording(t *testing.T) {
	at := time.Date(2026, 11, 5, 6, 0, 0, 0, time.UTC)
	journal := func(id string, es ...run.Entry) run.Journal {
		started := run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "build", Change: "C-" + id}}
		return run.Journal{ID: id, Entries: append([]run.Entry{started}, es...)}
	}
	enter := func(step string) run.Entry {
		return run.Entry{Type: run.TypeEntered, At: at.Add(time.Second), Step: step, Visit: 1}
	}
	runs := []run.Journal{
		journal("ROPEN", enter("gate"), run.Entry{Type: run.TypeStopped, At: at.Add(time.Minute), Step: "gate",
			Detail: run.OpenStop + "q-retry-cap, q-webhook-auth; answer each with a Decision, then rhiza run resume ROPEN"}),
		journal("RSIGN", enter("signoff"), run.Entry{Type: run.TypeStopped, At: at.Add(2 * time.Minute), Step: "signoff", Detail: run.SignoffWait}),
		journal("RFAIL", enter("open-pr"),
			run.Entry{Type: run.TypeEnded, At: at.Add(3 * time.Minute), Step: "open-pr", Visit: 1, Outcome: "fail", Detail: "exit 3\nRex tab window:01a1a8f1 was already closed"},
			run.Entry{Type: run.TypeFinished, At: at.Add(3 * time.Minute), Step: "open-pr", Detail: `step "open-pr" ended with outcome "fail"`, End: run.Failed}),
		journal("RLIVE", enter("leaves")),
	}
	live := func(id string) bool { return id == "RLIVE" }
	stalled := func(st run.State) []RunStalled {
		return []RunStalled{{Step: "leaves", Session: "SESS42", Since: at.Add(4 * time.Minute).Format(time.RFC3339)}}
	}
	items := inboxItems(runs, "h", nil, live, func(string) (string, bool) { return "", false }, stalled, nil, nil)
	kinds := map[string]bool{}
	for i := range items {
		kinds[items[i].Kind] = true
		items[i].ChangeTitle = "Title of " + strings.ToLower(strings.TrimPrefix(items[i].Change, "C-R"))
		if items[i].Kind == InboxOpen {
			items[i].Questions = []InboxGateQuestion{
				{Handle: "q-retry-cap", Statement: "How many times should a webhook be retried?"},
				{Handle: "q-webhook-auth", Statement: "Should webhooks be signed?", Blocks: "Webhook delivery"},
			}
		}
	}
	for _, k := range []string{InboxOpen, InboxSignoff, InboxFailed, InboxStalled} {
		if !kinds[k] {
			t.Fatalf("no %s item among %+v", k, items)
		}
	}
	got := webInboxViews(items)
	byTitle := map[string]web.InboxItemView{}
	for _, v := range got {
		for _, bad := range []string{"ROPEN", "RSIGN", "RFAIL", "RLIVE", "SESS42", "q-retry-cap", "q-webhook-auth", "rhiza ", "window:"} {
			if strings.Contains(v.Title, bad) || strings.Contains(v.Detail, bad) {
				t.Errorf("%s row shows %q: %+v", v.Kind, bad, v)
			}
		}
		byTitle[v.Title] = v
	}
	if len(got) != 5 {
		t.Fatalf("views = %+v, want two open questions, a sign-off, a failure and a stall", got)
	}
	want := map[string]string{
		"How many times should a webhook be retried?": "Title of open",
		"Should webhooks be signed?":                  "blocks Webhook delivery",
	}
	for title, detail := range want {
		v, ok := byTitle[title]
		if !ok || v.Kind != InboxOpen || v.Detail != detail || v.Run != "ROPEN" {
			t.Errorf("open question %q: %+v", title, v)
		}
	}
	if v := byTitle["Title of sign"]; v.Detail != "waiting for your sign-off" {
		t.Errorf("sign-off: %+v", v)
	}
	if v := byTitle["Title of fail"]; v.Detail != "the run failed at its open-pr step" {
		t.Errorf("failed: %+v", v)
	}
	if v := byTitle["Title of live"]; v.Kind != InboxStalled || v.Detail != "its session has stopped reporting" {
		t.Errorf("stalled: %+v", v)
	}
}

// An open item whose questions could not be read still names its change
// set, without the stop's handles.
func TestWebInboxViewsOpenWithoutQuestions(t *testing.T) {
	got := webInboxViews([]InboxItem{{Kind: InboxOpen, Run: "R1", Change: "C1", ChangeTitle: "Retry webhooks",
		Detail: run.OpenStop + "q-a; answer each with a Decision, then rhiza run resume R1"}})
	if len(got) != 1 || got[0].Title != "Retry webhooks" || got[0].Detail != "open questions wait on your answers" {
		t.Errorf("views = %+v", got)
	}
}

func TestWebInboxViewsPermission(t *testing.T) {
	got := webInboxViews([]InboxItem{{Kind: InboxPermission, Run: "R1", Change: "C1", ChangeTitle: "Retry webhooks",
		Session: "SESS1", Detail: "Allow rm -rf /home/x?", Next: "rhiza inbox open abc1234"}})
	if len(got) != 1 || got[0].Kind != InboxPermission || got[0].Title != "Retry webhooks" || got[0].Detail != "a session waits for your permission" {
		t.Errorf("views = %+v", got)
	}
}
