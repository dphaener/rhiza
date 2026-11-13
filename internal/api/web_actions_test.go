package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"

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

// A POST through the project's web handler answers a plan run's open
// question: the real question.answer writes the Decision into the run's
// worktree.
func TestWebAnswerAction(t *testing.T) {
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
	call("propose", `{}`)

	common := t.TempDir()
	runs := run.NewStore(run.Dir(common))
	at := ctx.Now()
	if err := runs.Create("RWEBGATE", run.Entry{Type: run.TypeStarted, At: at, Started: &run.Started{Workflow: "plan", Change: own.String(), Worktree: ctx.Root}}); err != nil {
		t.Fatal(err)
	}
	if err := runs.Append("RWEBGATE", run.Entry{Type: run.TypeStopped, At: at, Step: "gate", Detail: run.OpenStop + "web-answer-q; answer each with a Decision, then rhiza run resume RWEBGATE"}); err != nil {
		t.Fatal(err)
	}
	ctx.Git = func(args ...string) (string, error) { return common + "\n", nil }
	ctx.Exec = func(name string, args []string, dir string, env []string) ([]byte, error) {
		return []byte(common + "\n"), nil
	}
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

	h := WebProjectHandler(ctx)
	rec := webPost(h, web.AnswerPath, url.Values{"run": {"RWEBGATE"}, "question": {"web-answer-q"}, "text": {"Yes, sign them."}})
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "Answered. The run continues.") {
		t.Fatalf("answer = %d %s", rec.Code, rec.Body)
	}
	if d, err := ctxWorkspace(t, ctx).model(cs).Resolve("answer-web-answer-q"); err != nil || d.Kind != op.KindDecision {
		t.Errorf("no answering Decision in %s: %+v, %v", cs, d, err)
	}
	// A blank answer is refused in the owner's words and writes nothing.
	rec = webPost(h, web.AnswerPath, url.Values{"run": {"RWEBGATE"}, "question": {"web-answer-q"}, "text": {"  "}})
	if rec.Code != http.StatusBadRequest || !strings.Contains(rec.Body.String(), "blank") {
		t.Errorf("blank answer = %d %s", rec.Code, rec.Body)
	}
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
