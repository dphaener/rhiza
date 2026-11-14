package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/index"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/web"
)

// webPost posts form to target with the token's bearer header, as a script
// would; withCookie posts as a browser does, with the cookie and the CSRF
// field a page carries.
func webPost(h http.Handler, target string, form url.Values) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodPost, target, strings.NewReader(form.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", "Bearer "+apiTestToken)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

func setWebToken(t *testing.T, ctx *Ctx) {
	t.Helper()
	pinFreeBytes(ctx)
	if ctx.Wd == "" {
		ctx.Wd = ctx.Root
	}
	root, ok := rootcfg.Find(ctx.Wd)
	if !ok {
		t.Fatalf("no root above %s", ctx.Wd)
	}
	if err := rootcfg.SetDaemon(root, web.TokenKey, apiTestToken); err != nil {
		t.Fatal(err)
	}
}

// webGateWF is a plan run's workflow, enough for run.resume to queue its
// continuation.
const webGateWF = `[step.gate]
kind = "command"
run = "true"
`

// webGateFixture is a project, registered in a root at its own folder as
// repo, with plan run RWEBGATE stopped at its gate on two open questions of
// the change set it planned. It returns the project's ctx, that change set
// and the git common folder holding the run.
func webGateFixture(t *testing.T) (*Ctx, id.ID, string) {
	t.Helper()
	cache := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", cache)
	t.Setenv("HOME", cache)
	ctx := initCtx(t)
	call := func(name, req string) any {
		t.Helper()
		resp, err := callOp(t, ctx, name, req)
		if err != nil {
			t.Fatalf("%s %s: %v", name, req, err)
		}
		return resp
	}
	own := call("change.new", `{"intent":"Plan"}`).(ChangeNewResponse).ChangeSet
	call("add", `{"kind":"intent","handle":"web-answer-intent"}`)
	call("propose", `{}`)
	call("accept", `{"reason":"ok"}`)
	call("land", `{}`)
	cs := call("change.new", `{"intent":"Planned"}`).(ChangeNewResponse).ChangeSet
	call("link", fmt.Sprintf(`{"from":%q,"type":"serves","to":"web-answer-intent"}`, cs.String()))
	call("add", `{"kind":"question","handle":"web-answer-q","fields":["statement=Should webhooks be signed?"]}`)
	call("add", `{"kind":"question","handle":"web-answer-r","fields":["statement=Should webhooks retry?"]}`)
	call("propose", `{}`)

	common := t.TempDir()
	runs := run.NewStore(run.Dir(common))
	at := ctx.Now()
	if err := runs.Create("RWEBGATE", run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "plan", Definition: webGateWF, Change: own.String(), Worktree: ctx.Root}}); err != nil {
		t.Fatal(err)
	}
	if err := runs.Append("RWEBGATE", run.Entry{Type: run.TypeStopped, At: at, Step: "gate", Detail: run.OpenStop + "web-answer-q, web-answer-r; answer each with a Decision, then rhiza run resume RWEBGATE"}); err != nil {
		t.Fatal(err)
	}
	ctx.Git = func(args ...string) (string, error) { return common + "\n", nil }
	ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		return []byte(common + "\n"), nil
	}
	ctx.Start = (&shStart{}).start
	p, err := index.DefaultPath(ctx.Root)
	if err != nil {
		t.Fatal(err)
	}
	x, err := index.Open(ctx.Root, p)
	if err != nil {
		t.Fatal(err)
	}
	x.Close()
	underRoot(t, ctx)
	setWebToken(t, ctx)
	return ctx, cs, common
}

// A POST through the project's web handler, and through the root's merged
// inbox at /p/<project>/, answers a plan run's open questions: the real
// question.answer writes each Decision into the run's worktree, and once
// every question the gate stopped on is answered, run.resume queues the
// run's continuation, so the run continues as the page says.
func TestWebAnswerAction(t *testing.T) {
	queued := func(t *testing.T, common string) bool {
		t.Helper()
		q, err := daemon.ReadContinues(common)
		if err != nil {
			t.Fatal(err)
		}
		return slices.ContainsFunc(q, func(c daemon.Continue) bool { return c.Run == "RWEBGATE" })
	}
	answer := func(t *testing.T, h http.Handler, prefix string, ctx *Ctx, cs id.ID, common string) {
		t.Helper()
		rec := webPost(h, prefix+web.AnswerPath, url.Values{"run": {"RWEBGATE"}, "question": {"web-answer-q"}, "text": {"Yes, sign them."}})
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Answered.") {
			t.Fatalf("answer = %d %s", rec.Code, rec.Body)
		}
		if d, err := ctxWorkspace(t, ctx).model(cs).Resolve("answer-web-answer-q"); err != nil || d.Kind != op.KindDecision {
			t.Errorf("no answering Decision in %s: %+v, %v", cs, d, err)
		}
		// One question is still open: the run waits for it, and says so.
		if queued(t, common) || strings.Contains(rec.Body.String(), "The run continues.") {
			t.Errorf("run continued with a question open: %s", rec.Body)
		}
		// The same question again is refused in the owner's words and
		// writes no second Decision.
		rec = webPost(h, prefix+web.AnswerPath, url.Values{"run": {"RWEBGATE"}, "question": {"web-answer-q"}, "text": {"Yes, again."}})
		if rec.Code != http.StatusBadRequest || strings.Contains(rec.Body.String(), "RWEBGATE") {
			t.Errorf("second answer = %d %s", rec.Code, rec.Body)
		}
		if _, err := ctxWorkspace(t, ctx).model(cs).Resolve("answer-web-answer-q-2"); err == nil {
			t.Errorf("a second answer wrote a second Decision")
		}
		// A blank answer is refused in the owner's words and writes nothing.
		rec = webPost(h, prefix+web.AnswerPath, url.Values{"run": {"RWEBGATE"}, "question": {"web-answer-r"}, "text": {"  "}})
		if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "blank") {
			t.Errorf("blank answer = %d %s", rec.Code, rec.Body)
		}
		// The last question answered, the run's continuation is queued.
		rec = webPost(h, prefix+web.AnswerPath, url.Values{"run": {"RWEBGATE"}, "question": {"web-answer-r"}, "text": {"Yes, three times."}})
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Answered. The run continues.") {
			t.Fatalf("last answer = %d %s", rec.Code, rec.Body)
		}
		if !queued(t, common) {
			t.Errorf("no continuation queued after every question was answered")
		}
	}
	t.Run("project", func(t *testing.T) {
		ctx, cs, common := webGateFixture(t)
		answer(t, WebProjectHandler(ctx), "", ctx, cs, common)
	})
	t.Run("root", func(t *testing.T) {
		ctx, cs, common := webGateFixture(t)
		rootCtx := *ctx
		rootCtx.Start = nil // the project's git is the fixture's fake, unbounded
		h := WebRootHandler(&rootCtx, ctx.Root, "repo", WebProjectHandler(ctx))

		// The merged inbox lists the gate with an answer form that posts
		// into the project.
		r := httptest.NewRequest(http.MethodGet, "/inbox", nil)
		r.Header.Set("Authorization", "Bearer "+apiTestToken)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		if body := rec.Body.String(); !strings.Contains(body, `action="/p/repo`+web.AnswerPath+`"`) {
			t.Fatalf("merged inbox has no answer form:\n%s", body)
		}
		answer(t, h, "/p/repo", ctx, cs, common)
	})
}

// A POST through the project's handler, and through the root's merged inbox
// at /p/<project>/, approves a run waiting at its sign-off: the real
// run.approve journals it and queues the continuation.
func TestWebApproveAction(t *testing.T) {
	signoffs := func(rs *run.Store) int {
		return len(raceSignoffs(t, rs))
	}
	t.Run("project", func(t *testing.T) {
		f, rs := signoffRaceFixture(t)
		setWebToken(t, f.ctx)
		h := WebProjectHandler(f.ctx)
		rec := webPost(h, web.ApprovePath, url.Values{"run": {signoffRaceRun}})
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Approved. The run continues.") {
			t.Fatalf("approve = %d %s", rec.Code, rec.Body)
		}
		if n := signoffs(rs); n != 1 {
			t.Errorf("signoffs = %d, want 1", n)
		}
		// Again: the sign-off is done, so it is refused and nothing is journalled twice.
		rec = webPost(h, web.ApprovePath, url.Values{"run": {signoffRaceRun}})
		if rec.Code != http.StatusBadRequest || strings.Contains(rec.Body.String(), signoffRaceRun) {
			t.Errorf("second approve = %d %s", rec.Code, rec.Body)
		}
		if n := signoffs(rs); n != 1 {
			t.Errorf("signoffs = %d after a second approve", n)
		}
	})
	t.Run("root", func(t *testing.T) {
		f, rs := signoffRaceFixture(t)
		root := filepath.Dir(f.ctx.Root)
		setWebToken(t, f.ctx)
		rootCtx := *f.ctx
		rootCtx.Root, rootCtx.Wd = root, root
		rootCtx.Exec = nil // the projects' git is f.git, as the fixture's
		home := WebProjectHandler(f.ctx)
		h := WebRootHandler(&rootCtx, root, "repo", home)

		// The merged inbox lists the sign-off with an approve form that
		// posts into the project, and shows no run ID.
		r := httptest.NewRequest(http.MethodGet, "/inbox", nil)
		r.Header.Set("Authorization", "Bearer "+apiTestToken)
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, r)
		body := rec.Body.String()
		if !strings.Contains(body, `action="/p/repo`+web.ApprovePath+`"`) || !strings.Contains(body, `name="csrf"`) {
			t.Fatalf("merged inbox has no approve form:\n%s", body)
		}
		rec = webPost(h, "/p/repo"+web.ApprovePath, url.Values{"run": {signoffRaceRun}})
		if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Approved. The run continues.") {
			t.Fatalf("approve = %d %s", rec.Code, rec.Body)
		}
		if n := signoffs(rs); n != 1 {
			t.Errorf("signoffs = %d, want 1", n)
		}
	})
	t.Run("refused without the token", func(t *testing.T) {
		f, rs := signoffRaceFixture(t)
		setWebToken(t, f.ctx)
		r := httptest.NewRequest(http.MethodPost, web.ApprovePath, strings.NewReader("run="+signoffRaceRun))
		r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rec := httptest.NewRecorder()
		WebProjectHandler(f.ctx).ServeHTTP(rec, r)
		if rec.Code != http.StatusUnauthorized || signoffs(rs) != 0 {
			t.Errorf("approve without the token = %d, signoffs %d", rec.Code, signoffs(rs))
		}
	})
}

// A POST through the project's web handler, and through the root's at
// /p/<project>/, writes rhiza.daemon.paused and the project's
// rhiza.project.paused through daemon.set.
func TestWebPauseAction(t *testing.T) {
	read := func(t *testing.T, ctx *Ctx, key string) string {
		t.Helper()
		out, err := runDaemonSettings(ctx, DaemonSettingsRequest{})
		if err != nil {
			t.Fatal(err)
		}
		for _, s := range out.(DaemonSettingsResponse).Settings {
			if s.Key == key {
				return s.Value
			}
		}
		t.Fatalf("no setting %s", key)
		return ""
	}
	check := func(t *testing.T, ctx *Ctx, h http.Handler, prefix string) {
		t.Helper()
		for _, c := range []struct{ scope, key string }{{"daemon", "rhiza.daemon.paused"}, {"project", "rhiza.project.paused"}} {
			rec := webPost(h, prefix+web.PausePath, url.Values{"scope": {c.scope}, "paused": {"true"}})
			if rec.Code != 200 || read(t, ctx, c.key) != "true" {
				t.Fatalf("pause %s = %d %s; %s = %s", c.scope, rec.Code, rec.Body, c.key, read(t, ctx, c.key))
			}
			rec = webPost(h, prefix+web.PausePath, url.Values{"scope": {c.scope}, "paused": {"false"}})
			if rec.Code != 200 || read(t, ctx, c.key) != "false" {
				t.Fatalf("resume %s = %d %s", c.scope, rec.Code, rec.Body)
			}
		}
		if rec := webPost(h, prefix+web.PausePath, url.Values{"scope": {"other"}, "paused": {"true"}}); rec.Code != http.StatusBadRequest {
			t.Errorf("bad scope = %d", rec.Code)
		}
	}
	// realGit runs git in the fixture's project, so git config is written.
	realGit := func(f *sessionFixture) {
		if out, err := exec.Command("git", "init", "-q", f.ctx.Root).CombinedOutput(); err != nil {
			t.Fatalf("git init: %v %s", err, out)
		}
		f.ctx.Git = func(args ...string) (string, error) {
			out, err := exec.Command("git", append([]string{"-C", f.ctx.Root}, args...)...).Output()
			return string(out), err
		}
	}
	t.Run("project", func(t *testing.T) {
		f, _ := signoffRaceFixture(t)
		realGit(f)
		setWebToken(t, f.ctx)
		check(t, f.ctx, WebProjectHandler(f.ctx), "")
	})
	t.Run("root", func(t *testing.T) {
		f, _ := signoffRaceFixture(t)
		realGit(f)
		root := filepath.Dir(f.ctx.Root)
		setWebToken(t, f.ctx)
		rootCtx := *f.ctx
		rootCtx.Root, rootCtx.Wd = root, root
		rootCtx.Exec = nil
		h := WebRootHandler(&rootCtx, root, "repo", WebProjectHandler(f.ctx))
		check(t, f.ctx, h, "/p/repo")
	})
}
