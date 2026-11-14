package web

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type actionCalls struct{ answers, approves []string }

// actionPage is a project's handler with the token gate, one open question
// and one sign-off waiting, and fake actions that record what they run.
func actionPage(token string) (http.Handler, *actionCalls) {
	calls := &actionCalls{}
	inbox := func() []InboxItemView {
		return []InboxItemView{
			{Kind: "open", Title: "Should webhooks be signed?", Detail: "Webhooks", Run: "RGATE", Change: "C1", Question: "q-sign", At: "2026-11-01T10:00:00Z"},
			{Kind: "signoff", Title: "Webhooks", Detail: "waiting for your sign-off", Run: "RSIGN", Change: "C1", At: "2026-11-02T10:00:00Z"},
		}
	}
	a := Actions{
		Answer: func(run, question, text string) (string, error) {
			if text == "boom" {
				return "", &ActionError{Status: http.StatusConflict, Message: "that question is answered already"}
			}
			if text == "plain" {
				return "", errors.New("plain failure")
			}
			calls.answers = append(calls.answers, run+"|"+question+"|"+text)
			return "Answered. The run continues.", nil
		},
		Approve: func(run string) (string, error) {
			calls.approves = append(calls.approves, run)
			return "Approved. The run continues.", nil
		},
	}
	h := HandlerWithInbox(nil, nil, inbox, nil)
	return WithToken(func() string { return token }, WithActions(a, h)), calls
}

func TestWebAnswerFormsCarryCSRFAndAreAbsentWithoutToken(t *testing.T) {
	h, _ := actionPage(testToken)
	cookie := map[string]string{"Cookie": TokenCookie + "=" + testToken}
	body := do(h, "GET", "/", cookie).Body.String()
	for _, want := range []string{
		`action="` + AnswerPath + `"`, `action="` + ApprovePath + `"`,
		`name="csrf" value="` + csrfFor(testToken) + `"`, `name="question" value="q-sign"`, `name="run" value="RGATE"`, `name="run" value="RSIGN"`,
		`name="text"`, `>Answer</button>`, `>Approve</button>`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page lacks %q", want)
		}
	}
	if strings.Contains(body, testToken) {
		t.Error("page shows the token")
	}
	for name, hdr := range map[string]map[string]string{"no token": nil, "wrong cookie": {"Cookie": TokenCookie + "=nope"}} {
		if b := do(h, "GET", "/", hdr).Body.String(); strings.Contains(b, "<form") || strings.Contains(b, `name="csrf"`) {
			t.Errorf("%s: page carries a form", name)
		}
	}
	none, _ := actionPage("")
	if b := do(none, "GET", "/", nil).Body.String(); strings.Contains(b, "<form") {
		t.Error("page carries a form with no token set")
	}
}

func TestWebAnswerPostConfirms(t *testing.T) {
	h, calls := actionPage(testToken)
	cookie := map[string]string{"Cookie": TokenCookie + "=" + testToken}
	csrf := csrfFor(testToken)
	rec := postTo(h, AnswerPath, cookie, url.Values{CSRFField: {csrf}, "run": {"RGATE"}, "question": {"q-sign"}, "text": {"Yes, sign them."}})
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Answered. The run continues.") {
		t.Fatalf("answer = %d %s", rec.Code, rec.Body)
	}
	if len(calls.answers) != 1 || calls.answers[0] != "RGATE|q-sign|Yes, sign them." {
		t.Errorf("answers = %v", calls.answers)
	}
	rec = postTo(h, ApprovePath, cookie, url.Values{CSRFField: {csrf}, "run": {"RSIGN"}})
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Approved. The run continues.") {
		t.Fatalf("approve = %d %s", rec.Code, rec.Body)
	}
	if len(calls.approves) != 1 || calls.approves[0] != "RSIGN" {
		t.Errorf("approves = %v", calls.approves)
	}
	// A bearer header needs no field.
	rec = postTo(h, ApprovePath, map[string]string{"Authorization": "Bearer " + testToken}, url.Values{"run": {"RSIGN"}})
	if rec.Code != 200 {
		t.Errorf("bearer approve = %d", rec.Code)
	}
	// A refusal shows its reason and status.
	rec = postTo(h, AnswerPath, cookie, url.Values{CSRFField: {csrf}, "run": {"RGATE"}, "question": {"q"}, "text": {"boom"}})
	if rec.Code != http.StatusConflict || !strings.Contains(rec.Body.String(), "that question is answered already") {
		t.Errorf("refusal = %d %s", rec.Code, rec.Body)
	}
	rec = postTo(h, AnswerPath, cookie, url.Values{CSRFField: {csrf}, "run": {"RGATE"}, "question": {"q"}, "text": {"plain"}})
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "plain failure") {
		t.Errorf("plain refusal = %d %s", rec.Code, rec.Body)
	}
}

func TestWebApproveRefusesMissingCSRF(t *testing.T) {
	h, calls := actionPage(testToken)
	cookie := map[string]string{"Cookie": TokenCookie + "=" + testToken}
	for name, form := range map[string]url.Values{
		"missing": {"run": {"RSIGN"}},
		"wrong":   {CSRFField: {"nope"}, "run": {"RSIGN"}},
	} {
		if rec := postTo(h, ApprovePath, cookie, form); rec.Code != http.StatusForbidden {
			t.Errorf("approve %s = %d, want 403", name, rec.Code)
		}
		if rec := postTo(h, AnswerPath, cookie, form); rec.Code != http.StatusForbidden {
			t.Errorf("answer %s = %d, want 403", name, rec.Code)
		}
	}
	if len(calls.answers)+len(calls.approves) != 0 {
		t.Errorf("an action ran: %+v", calls)
	}
	// No token set: nothing runs either.
	off, offCalls := actionPage("")
	if rec := postTo(off, ApprovePath, nil, url.Values{"run": {"RSIGN"}}); rec.Code != http.StatusForbidden || len(offCalls.approves) != 0 {
		t.Errorf("approve with no token set = %d, ran %v", rec.Code, offCalls.approves)
	}
}

func postTo(h http.Handler, path string, hdr map[string]string, form url.Values) *httptest.ResponseRecorder {
	r := httptest.NewRequest("POST", path, strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range hdr {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

func TestWebAnswerMergedInboxFormsPostIntoTheProject(t *testing.T) {
	src := func() []ProjectInbox {
		return []ProjectInbox{{Name: "alpha", Items: []InboxItemView{
			{Kind: "open", Title: "Q?", Run: "R1", Question: "q1", At: "2026-11-01T10:00:00Z"},
			{Kind: "signoff", Title: "S", Run: "R2", At: "2026-11-02T10:00:00Z"},
		}}}
	}
	h := WithToken(func() string { return testToken }, RootWithInbox(Handler(nil), "", func(string) http.Handler { return nil }, nil, src))
	body := do(h, "GET", "/inbox", map[string]string{"Cookie": TokenCookie + "=" + testToken}).Body.String()
	for _, want := range []string{`action="/p/alpha` + AnswerPath + `"`, `action="/p/alpha` + ApprovePath + `"`, `name="csrf"`} {
		if !strings.Contains(body, want) {
			t.Errorf("merged inbox lacks %q", want)
		}
	}
	if b := do(h, "GET", "/inbox", nil).Body.String(); strings.Contains(b, "<form") {
		t.Error("merged inbox carries a form without the token")
	}
}
