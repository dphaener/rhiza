package web_test

import (
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/dphaener/rhiza/internal/api"
	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/web"
)

func fixtureCtx(t *testing.T) *api.Ctx {
	t.Helper()
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }
	root := t.TempDir()
	if out, err := exec.Command("git", "-C", root, "init", "-q").CombinedOutput(); err != nil {
		t.Skipf("git init: %v %s", err, out)
	}
	ctx := &api.Ctx{
		Root: root,
		Now:  clock,
		IDs:  id.NewGenerator(clock, rand.New(rand.NewSource(3))),
		Git: func(args ...string) (string, error) {
			c := exec.Command("git", args...)
			c.Dir = root
			out, err := c.Output()
			return string(out), err
		},
		IndexPath: filepath.Join(t.TempDir(), "index.db"),
		GitConfig: func(k string) (string, error) {
			if k == "user.email" {
				return "test@example.com", nil
			}
			return "", errors.New("unset")
		},
	}
	call := func(op, req string) {
		t.Helper()
		o, _ := api.Lookup(op)
		r, err := api.DecodeRequest(o, []byte(req))
		if err == nil {
			_, err = o.Run(ctx, r)
		}
		if err != nil {
			t.Fatalf("%s: %v", op, err)
		}
	}
	call("init", `{}`)
	call("change.new", `{"intent":"A change"}`)
	return ctx
}

// want is the bytes `rhiza <op> --json` prints for the request.
func want(t *testing.T, ctx *api.Ctx, op, req string) string {
	t.Helper()
	o, _ := api.Lookup(op)
	r, err := api.DecodeRequest(o, []byte(req))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := o.Run(ctx, r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := canon.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}

func serve(h http.Handler, method, target string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(method, target, nil))
	return rec
}

func TestAPIAnswersTheHandlersJSON(t *testing.T) {
	ctx := fixtureCtx(t)
	h := web.Handler(api.ReadOnlyHandler(ctx))
	for _, c := range []struct{ target, op, req string }{
		{"/api/run.ls", "run.ls", `{}`},
		{"/api/run.ls?limit=5", "run.ls", `{"limit":5}`},
		{"/api/roadmap", "roadmap", `{}`},
	} {
		rec := serve(h, "GET", c.target)
		if rec.Code != 200 {
			t.Errorf("%s: %d %s", c.target, rec.Code, rec.Body)
			continue
		}
		if got, w := rec.Body.String(), want(t, ctx, c.op, c.req); got != w {
			t.Errorf("%s:\n got %s\nwant %s", c.target, got, w)
		}
	}
}

func TestAPIRefusesWrites(t *testing.T) {
	h := web.Handler(api.ReadOnlyHandler(fixtureCtx(t)))
	for _, op := range []string{"land", "accept", "daemon.set", "daemon.start", "render", "ship.merge", "add", "apply", "nope"} {
		if rec := serve(h, "GET", "/api/"+op); rec.Code != 404 && rec.Code != 405 {
			t.Errorf("GET /api/%s: %d", op, rec.Code)
		}
	}
	for _, m := range []string{"POST", "PUT", "DELETE", "PATCH"} {
		for _, op := range []string{"run.ls", "render", "daemon.set"} {
			if rec := serve(h, m, "/api/"+op); rec.Code != 405 {
				t.Errorf("%s /api/%s: %d", m, op, rec.Code)
			}
		}
	}
	if rec := serve(h, "GET", "/api/run.ls?bogus=1"); rec.Code != 400 {
		t.Errorf("unknown parameter: %d", rec.Code)
	}
	if rec := serve(web.Handler(nil), "GET", "/api/run.ls"); rec.Code != 404 {
		t.Errorf("no api: %d", rec.Code)
	}
}
