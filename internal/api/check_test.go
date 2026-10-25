package api

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/genesis/internal/check"
)

// fakeGit answers git calls from a table keyed by the space-joined
// arguments; any other call fails. It records every call.
type fakeGit struct {
	out   map[string]string
	calls []string
}

func (g *fakeGit) run(args ...string) (string, error) {
	k := strings.Join(args, " ")
	g.calls = append(g.calls, k)
	if out, ok := g.out[k]; ok {
		return out, nil
	}
	return "", errors.New("exit status 1")
}

const (
	gitWorkTree = "rev-parse --is-inside-work-tree"
	diffArgs    = "diff --name-status --no-renames -z --relative "
)

func verifyArgs(rev string) string { return "rev-parse --verify --quiet " + rev + "^{commit}" }

func TestCheckBase(t *testing.T) {
	inTree := map[string]string{gitWorkTree: "true\n"}
	with := func(kv ...string) map[string]string {
		m := map[string]string{}
		for k, v := range inTree {
			m[k] = v
		}
		for i := 0; i < len(kv); i += 2 {
			m[kv[i]] = kv[i+1]
		}
		return m
	}
	for _, c := range []struct {
		name               string
		git                map[string]string // nil: no git at all
		base               string
		wantBase, wantFrom string
		wantSkip           string
		wantUsage          bool
	}{
		{name: "NoGit", wantSkip: "not a git work tree"},
		{name: "NotWorkTree", git: map[string]string{}, wantSkip: "not a git work tree"},
		{name: "WorkTreeFalse", git: map[string]string{gitWorkTree: "false\n"}, wantSkip: "not a git work tree"},
		{name: "NoCandidate", git: with(), wantSkip: "no base: none of origin/HEAD, origin/main, main resolves"},
		{name: "Main", git: with(verifyArgs("main"), "m\n", "merge-base HEAD main", "abc123\n"),
			wantBase: "abc123", wantFrom: "main"},
		{name: "OriginHEADFirst", git: with(verifyArgs("origin/HEAD"), "o\n", "merge-base HEAD origin/HEAD", "fff\n",
			verifyArgs("main"), "m\n", "merge-base HEAD main", "abc123\n"),
			wantBase: "fff", wantFrom: "origin/HEAD"},
		{name: "OriginMain", git: with(verifyArgs("origin/main"), "o\n", "merge-base HEAD origin/main", "eee\n"),
			wantBase: "eee", wantFrom: "origin/main"},
		// main resolves but shares no history with HEAD (or HEAD has no
		// commit yet): no base.
		{name: "NoMergeBase", git: with(verifyArgs("main"), "m\n"),
			wantSkip: "no base: none of origin/HEAD, origin/main, main resolves"},
		{name: "Explicit", git: with(verifyArgs("HEAD~1"), "h\n"), base: "HEAD~1", wantBase: "HEAD~1"},
		{name: "ExplicitUnknown", git: with(), base: "nope", wantUsage: true},
		{name: "ExplicitNoGit", base: "HEAD", wantSkip: "not a git work tree"},
	} {
		t.Run(c.name, func(t *testing.T) {
			var git func(...string) (string, error)
			if c.git != nil {
				git = (&fakeGit{out: c.git}).run
			}
			base, from, skip, err := checkBase(git, c.base)
			if c.wantUsage {
				wantExit(t, err, 2)
				if !strings.Contains(err.Error(), `unknown base "nope"`) {
					t.Errorf("err = %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if base != c.wantBase || from != c.wantFrom || skip != c.wantSkip {
				t.Errorf("checkBase = %q, %q, %q; want %q, %q, %q", base, from, skip, c.wantBase, c.wantFrom, c.wantSkip)
			}
		})
	}
}

func TestParseNameStatus(t *testing.T) {
	for _, c := range []struct {
		in   string
		want []check.Change
	}{
		{"", nil},
		{"A\x00.genesis/log/a/b.json\x00", []check.Change{{Status: "A", Path: ".genesis/log/a/b.json"}}},
		{"M\x00.genesis/log/a/b.json\x00D\x00.genesis/log/a/c d.json\x00T\x00.genesis/log/\"q\".json\x00",
			[]check.Change{{Status: "M", Path: ".genesis/log/a/b.json"}, {Status: "D", Path: ".genesis/log/a/c d.json"},
				{Status: "T", Path: `.genesis/log/"q".json`}}},
		// A trailing status with no path is dropped.
		{"M\x00", nil},
	} {
		got := parseNameStatus(c.in)
		if !slices.Equal(got, c.want) {
			t.Errorf("parseNameStatus(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestCheckInputChanges(t *testing.T) {
	ctx := initCtx(t)
	g := &fakeGit{out: map[string]string{
		gitWorkTree:                         "true\n",
		verifyArgs("main"):                  "m\n",
		"merge-base HEAD main":              "abc\n",
		diffArgs + "abc -- .genesis/log":    "M\x00.genesis/log/x/y.json\x00",
		diffArgs + "HEAD~1 -- .genesis/log": "",
		verifyArgs("HEAD~1"):                "h\n",
		diffArgs + "broken -- .genesis/log": "",
		verifyArgs("broken"):                "b\n",
	}}
	ctx.Git = g.run
	resp, err := callOp(t, ctx, "check", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(CheckResponse)
	if r.Base != "abc" || r.MergeBaseOf != "main" || r.Skipped != "" {
		t.Errorf("base %q from %q skipped %q", r.Base, r.MergeBaseOf, r.Skipped)
	}
	if !slices.ContainsFunc(r.Findings, func(f CheckFinding) bool {
		return f.Rule == check.RuleAppendOnly && f.Path == ".genesis/log/x/y.json" && f.Next == "git checkout abc -- .genesis/log/x/y.json"
	}) {
		t.Errorf("no append-only finding: %+v", r.Findings)
	}
	if r.OK || r.ExitCode() != 1 {
		t.Errorf("ok %v exit %d", r.OK, r.ExitCode())
	}
	out := human(t, r)
	for _, s := range []string{"failed\n", "base: abc (merge-base of HEAD and main)\n",
		"fail append-only .genesis/log/x/y.json: log file modified",
		"  next: git checkout abc -- .genesis/log/x/y.json\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("human lacks %q:\n%s", s, out)
		}
	}

	// A diff that fails is an error, not a silent pass.
	delete(g.out, diffArgs+"broken -- .genesis/log")
	_, err = callOp(t, ctx, "check", `{"base":"broken"}`)
	wantExit(t, err, 1)

	_, err = callOp(t, ctx, "check", `{"base":"nope"}`)
	wantExit(t, err, 2)
}

func TestCheckFreshRepo(t *testing.T) {
	ctx := initCtx(t)
	renderOK(t, ctx)
	resp, err := callOp(t, ctx, "check", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(CheckResponse)
	if !r.OK || r.ExitCode() != 0 || r.Skipped != "not a git work tree" || r.Base != "" {
		t.Fatalf("check = %+v", r)
	}
	want := []CheckFinding{{Severity: "warn", Rule: "append-only", Message: "append-only check skipped: not a git work tree"}}
	if !slices.Equal(r.Findings, want) {
		t.Errorf("findings = %+v", r.Findings)
	}
	want0 := "ok\nappend-only: skipped (not a git work tree)\n\n" +
		"warn append-only: append-only check skipped: not a git work tree\n"
	if out := human(t, r); out != want0 {
		t.Errorf("human:\n%s\nwant:\n%s", out, want0)
	}
}

func TestCheckFailsFirst(t *testing.T) {
	ctx := initCtx(t)
	// No render: the rendered folder is missing.
	resp, err := callOp(t, ctx, "check", `{}`)
	if err != nil {
		t.Fatal(err)
	}
	r := resp.(CheckResponse)
	if r.OK || r.ExitCode() != 1 {
		t.Fatalf("check = %+v", r)
	}
	want := "failed\nappend-only: skipped (not a git work tree)\n\n" +
		"fail render .genesis/rendered: the rendered folder is missing\n" +
		"  next: genesis render\n" +
		"warn append-only: append-only check skipped: not a git work tree\n"
	if out := human(t, r); out != want {
		t.Errorf("human:\n%s\nwant:\n%s", out, want)
	}

	// A file where the folder should be is missing too.
	if err := os.WriteFile(filepath.Join(ctx.Root, renderedDir), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	resp, _ = callOp(t, ctx, "check", `{}`)
	if resp.(CheckResponse).OK {
		t.Error("a rendered file passed")
	}
}

func TestCheckFindingWhere(t *testing.T) {
	for _, c := range []struct {
		f    CheckFinding
		want string
	}{
		{CheckFinding{Severity: "fail", Rule: "log", Path: ".genesis/log/a", Message: "bad"}, "fail log .genesis/log/a: bad"},
		{CheckFinding{Severity: "fail", Rule: "model", ChangeSet: "C", Op: "O", Message: "m"}, "fail model change set C, op O: m"},
		{CheckFinding{Severity: "fail", Rule: "link", ChangeSet: "C", Op: "O", Link: "L", Message: "m"}, "fail link change set C, op O, link L: m"},
		{CheckFinding{Severity: "warn", Rule: "fold", ChangeSet: "C", Entry: "E", Message: "m"}, "warn fold change set C, entry E: m"},
		{CheckFinding{Severity: "warn", Rule: "append-only", Message: "m"}, "warn append-only: m"},
	} {
		if got := c.f.line(); got != c.want {
			t.Errorf("line = %q, want %q", got, c.want)
		}
	}
}

func TestCheckErrors(t *testing.T) {
	ctx := testCtx(t.TempDir())
	_, err := callOp(t, ctx, "check", `{}`)
	wantExit(t, err, 4)
	wantNext(t, err, "genesis init")

	ctx = initCtx(t)
	ctx.DryRun = true
	_, err = callOp(t, ctx, "check", `{}`)
	wantExit(t, err, 2)
	wantNext(t, err, "genesis check")
}
