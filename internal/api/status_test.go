package api

import (
	"os"
	"path/filepath"
	"reflect"
	"slices"
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
	t.Parallel()
	ctx := testCtx(t, t.TempDir())
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
	t.Parallel()
	root := t.TempDir()
	ctx := testCtx(t, root)
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
	t.Parallel()
	root := t.TempDir()
	ctx := testCtx(t, root)
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

// TestStatusGuideFiles checks an install at the running version is current
// only while every file it writes is on disk as written.
func TestStatusGuideFiles(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	ctx := testCtx(t, root)
	ctx.Home = t.TempDir()
	if _, err := callOp(t, ctx, "init", `{}`); err != nil {
		t.Fatal(err)
	}
	guideInstall(t, ctx, `{"agent":"claude"}`)
	guideInstall(t, ctx, `{"agent":"claude","project":true}`)
	dir := filepath.Join(ctx.Home, ".claude", "skills", "rhiza")
	if r := statusOf(t, ctx, `{}`); len(r.Guide) != 2 || !r.Guide[0].Current || !r.Guide[1].Current {
		t.Fatalf("after install guide = %+v", r.Guide)
	}
	for _, p := range pluginPaths[1:] {
		t.Run(p, func(t *testing.T) {
			f := filepath.Join(dir, filepath.FromSlash(p))
			for _, tc := range []struct {
				name  string
				spoil func() error
			}{
				{"Missing", func() error { return os.Remove(f) }},
				{"Changed", func() error { return os.WriteFile(f, []byte("changed\n"), 0o644) }},
			} {
				if err := tc.spoil(); err != nil {
					t.Fatal(err)
				}
				r := statusOf(t, ctx, `{}`)
				if len(r.Guide) != 2 || r.Guide[0].Current || r.Guide[0].Version != ctx.Version || !r.Guide[1].Current {
					t.Errorf("%s: guide = %+v", tc.name, r.Guide)
				}
				if !slices.Contains(r.Next, "rhiza guide install --agent claude") {
					t.Errorf("%s: next = %v", tc.name, r.Next)
				}
				guideInstall(t, ctx, `{"agent":"claude"}`)
				if r := statusOf(t, ctx, `{}`); !r.Guide[0].Current {
					t.Errorf("%s: reinstall not current: %+v", tc.name, r.Guide)
				}
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	t.Parallel()
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

// TestStatusIndexedReadsNoLog: with an index, status and its hook count
// the change sets from the projection, reading the log zero times, and
// give the counts and change set status gives over a throwaway index.
func TestStatusIndexedReadsNoLog(t *testing.T) {
	t.Parallel()
	ctx, _, accepted := readLog(t)
	newChange(t, ctx, `{"intent":"Later"}`)
	if _, err := callOp(t, ctx, "change.use", `{"ref":"`+accepted.String()+`"}`); err != nil {
		t.Fatal(err)
	}
	statusOf(t, ctx, `{}`) // catch the index up with the log
	reads := countingCtx(ctx)
	for _, req := range []string{`{}`, `{"hook":true}`} {
		indexed := statusOf(t, ctx, req)
		if *reads != 0 {
			t.Errorf("status %s over the index read the log %d times, want 0", req, *reads)
		}
		folded := *ctx
		// No index path on purpose: the call builds a throwaway index.
		folded.IndexPath = ""
		logReads := countingCtx(&folded)
		logged := statusOf(t, &folded, req)
		if *logReads != 0 {
			t.Errorf("status %s over a throwaway index read the log %d times, want 0", req, *logReads)
		}
		if indexed.Attention != logged.Attention || !reflect.DeepEqual(indexed.ChangeSet, logged.ChangeSet) {
			t.Errorf("status %s over the index = %+v %+v, over the log %+v %+v",
				req, indexed.Attention, indexed.ChangeSet, logged.Attention, logged.ChangeSet)
		}
		if indexed.ChangeSet == nil || indexed.ChangeSet.ID != accepted || indexed.Attention.Draft != 1 || indexed.Attention.Accepted != 1 {
			t.Errorf("status %s = %+v %+v, want the accepted set current, 1 draft and 1 accepted", req, indexed.Attention, indexed.ChangeSet)
		}
		if a, b := human(t, indexed), human(t, logged); a != b {
			t.Errorf("status %s text over the index:\n%s\nover the log:\n%s", req, a, b)
		}
	}
}

// TestStatusHookBadIndexSilent: an index that cannot be opened, with no
// temporary folder for a throwaway one, makes the hook silent with exit 0.
func TestStatusHookBadIndexSilent(t *testing.T) {
	t.Parallel()
	ctx := readCtx(t)
	newChange(t, ctx, `{"intent":"Here"}`)
	// The index's folder is a regular file, so the index cannot be opened.
	if err := os.WriteFile(ctx.IndexPath, []byte("not a folder"), 0o644); err != nil {
		t.Fatal(err)
	}
	ctx.IndexPath = filepath.Join(ctx.IndexPath, "index.db")
	// The temporary folder is a regular file too.
	withTempDir(ctx, filepath.Dir(ctx.IndexPath))
	if _, err := callOp(t, ctx, "status", `{}`); err == nil {
		t.Error("status over a broken index: no error")
	}
	r := statusOf(t, ctx, `{"hook":true}`)
	if out := human(t, r); out != "" {
		t.Errorf("hook over a broken index printed %q", out)
	}
	if r.ExitCode() != 0 {
		t.Errorf("hook exit code %d", r.ExitCode())
	}
}
