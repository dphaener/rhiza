package api

import (
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/canon"
	"github.com/dphaener/rhiza/internal/guide"
)

func statusOf(t *testing.T, ctx *Ctx, req string) StatusResponse {
	t.Helper()
	resp, err := callOp(t, ctx, "status", req)
	if err != nil {
		t.Fatal(err)
	}
	return resp.(StatusResponse)
}

func TestStatusOutsideRepo(t *testing.T) {
	ctx := testCtx(t.TempDir())
	if _, err := callOp(t, ctx, "status", `{}`); ExitCode(err) != 4 {
		t.Errorf("status outside a repo: err %v, want not_found", err)
	}
	r := statusOf(t, ctx, `{"hook":true}`)
	if out := human(t, r); out != "" {
		t.Errorf("hook outside a repo printed %q", out)
	}
	if b, err := canon.Marshal(r); err != nil || string(b) != "{}\n" {
		t.Errorf("hook JSON outside a repo = %q, %v; want {}", b, err)
	}
	if r.ExitCode() != 0 {
		t.Errorf("hook exit code %d", r.ExitCode())
	}
}

func TestStatusHookErrorSilent(t *testing.T) {
	root := t.TempDir()
	ctx := testCtx(root)
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	ctx.ChangeRef = "ZZZZZZZZ"
	if _, err := callOp(t, ctx, "status", `{}`); err == nil {
		t.Fatal("status with a bad --change: no error")
	}
	if out := human(t, statusOf(t, ctx, `{"hook":true}`)); out != "" {
		t.Errorf("hook after an error printed %q", out)
	}
}

// TestStatusHookCap checks the hook's worst case stays under guide.StartCap:
// a long intent, both installs out of date and a long version.
func TestStatusHookCap(t *testing.T) {
	root := t.TempDir()
	ctx := testCtx(root)
	ctx.Home = t.TempDir()
	ctx.Version = "v0.0.1"
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	guideInstall(t, ctx, `{"agent":"claude"}`)
	guideInstall(t, ctx, `{"agent":"claude","project":true}`)
	ctx.Version = "v10.200.3000-rc.12-g0123456789ab-dirty"
	intent := strings.Repeat("Make the onboarding flow explain every step in plain words ", 5)
	if _, err := callOp(t, ctx, "change.new", `{"intent":"`+intent+`"}`); err != nil {
		t.Fatal(err)
	}
	r := statusOf(t, ctx, `{"hook":true}`)
	out := human(t, r)
	t.Logf("hook output (%d bytes, %d tokens):\n%s", len(out), guide.Tokens(out), out)
	if n := guide.Tokens(out); n > guide.StartCap {
		t.Errorf("hook output is %d tokens, cap %d", n, guide.StartCap)
	}
	for _, want := range []string{
		"This repo uses Rhiza",
		"(draft): Make the onboarding",
		"…\n",
		"In flight: 1 draft, 0 proposed, 0 accepted.",
		"`rhiza guide install --agent claude`, `rhiza guide install --agent claude --project`",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("hook output lacks %q", want)
		}
	}
	if len(r.Guide) != 2 || r.Guide[0].Scope != "user" || r.Guide[1].Scope != "project" || r.Guide[0].Current {
		t.Errorf("guide = %+v", r.Guide)
	}
	full := statusOf(t, ctx, `{}`)
	a, _ := canon.Marshal(full)
	b, _ := canon.Marshal(r)
	if string(a) != string(b) {
		t.Errorf("hook JSON differs from status JSON:\n%s\n%s", b, a)
	}
}

func TestTruncate(t *testing.T) {
	for _, tc := range []struct {
		in   string
		n    int
		want string
	}{
		{"short", 10, "short"},
		{"two\nlines", 20, "two lines"},
		{"abcdefghij", 6, "abc…"},
		{"ééééé", 6, "é…"},
	} {
		if got := truncate(tc.in, tc.n); got != tc.want || len(got) > tc.n {
			t.Errorf("truncate(%q, %d) = %q, want %q", tc.in, tc.n, got, tc.want)
		}
	}
}
