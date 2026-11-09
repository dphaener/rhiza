package api

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/docs"
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
	lsTreeArgs  = "ls-tree -r -z "
	logSpec     = " -- .rhiza/log .genesis/log"
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

func TestGitBlob(t *testing.T) {
	// The IDs `git hash-object` prints for these bytes.
	for in, want := range map[string]string{
		"":        "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391",
		"hello\n": "ce013625030ba8dba906f756967f9e9ca394464a",
	} {
		if got := gitBlob([]byte(in)); got != want {
			t.Errorf("gitBlob(%q) = %s, want %s", in, got, want)
		}
	}
}

func TestParseLsTree(t *testing.T) {
	line := func(typ, obj, p string) string { return "100644 " + typ + " " + obj + "\t" + p + "\x00" }
	for _, c := range []struct {
		name string
		in   string
		want []check.Entry
	}{
		{name: "Empty"},
		{name: "One", in: line("blob", "b1", ".rhiza/log/cs/e.json"),
			want: []check.Entry{{ID: "e", Path: ".rhiza/log/cs/e.json", Blob: "b1"}}},
		{name: "SortedOnceVerbatim",
			in: line("blob", "b2", ".rhiza/log/cs/z \"q\".json") + line("blob", "b1", ".rhiza/log/cs/a.json") +
				line("blob", "b1", ".rhiza/log/cs/a.json"),
			want: []check.Entry{{ID: "a", Path: ".rhiza/log/cs/a.json", Blob: "b1"},
				{ID: `z "q"`, Path: `.rhiza/log/cs/z "q".json`, Blob: "b2"}}},
		// Not entries: the .gitkeep, dot folders and files, other names and
		// depths, other folders, a submodule, a truncated line.
		{name: "Skips", in: line("blob", "k", ".rhiza/log/.gitkeep") +
			line("blob", "k", ".rhiza/log/.x/e.json") + line("blob", "k", ".rhiza/log/cs/.e.json") +
			line("blob", "k", ".rhiza/log/cs/e.txt") + line("blob", "k", ".rhiza/log/e.json") +
			line("blob", "k", ".rhiza/log/cs/d/e.json") + line("blob", "k", ".rhiza/repo.json") +
			line("blob", "k", "other/log/cs/e.json") + line("commit", "k", ".rhiza/log/cs/m.json") +
			"100644 blob k"},
	} {
		t.Run(c.name, func(t *testing.T) {
			if got := parseLsTree(c.in); !slices.Equal(got, c.want) {
				t.Errorf("parseLsTree = %+v, want %+v", got, c.want)
			}
		})
	}
}

func TestCurrentEntries(t *testing.T) {
	root := t.TempDir()
	if es, err := currentEntries(root); err != nil || es != nil {
		t.Fatalf("no log: %v, %v", es, err)
	}
	logDir := filepath.Join(root, ".rhiza", "log")
	for p, b := range map[string]string{
		"cs2/e2.json": "hello\n", "cs1/e1.json": "", "cs1/bad.json": "{not json",
		".gitkeep": "", "e0.json": "", ".hidden/e3.json": "", "cs1/.e4.json": "", "cs1/e5.txt": "",
	} {
		full := filepath.Join(logDir, filepath.FromSlash(p))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(b), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Mkdir(filepath.Join(logDir, "cs1", "dir.json"), 0o755); err != nil {
		t.Fatal(err)
	}
	es, err := currentEntries(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []check.Entry{
		{ID: "bad", Path: ".rhiza/log/cs1/bad.json", Blob: gitBlob([]byte("{not json"))},
		{ID: "e1", Path: ".rhiza/log/cs1/e1.json", Blob: "e69de29bb2d1d6434b8b29ae775ad8c2e48c5391"},
		{ID: "e2", Path: ".rhiza/log/cs2/e2.json", Blob: "ce013625030ba8dba906f756967f9e9ca394464a"},
	}
	if !slices.Equal(es, want) {
		t.Errorf("currentEntries = %+v, want %+v", es, want)
	}
}

// writeEntry writes bytes b at the entry path p under ctx's root.
func writeEntry(t *testing.T, ctx *Ctx, p, b string) {
	t.Helper()
	full := filepath.Join(ctx.Root, filepath.FromSlash(p))
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(b), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCheckInputChanges(t *testing.T) {
	ctx := initCtx(t)
	const (
		edited = ".rhiza/log/x/y.json"
		gone   = ".rhiza/log/x/z.json"
		moved  = ".rhiza/log/w/m.json"
	)
	writeEntry(t, ctx, edited, "now\n")
	writeEntry(t, ctx, ".rhiza/log/v/m.json", "same\n")
	base := func(p, b string) string { return "100644 blob " + gitBlob([]byte(b)) + "\t" + p + "\x00" }
	g := &fakeGit{out: map[string]string{
		gitWorkTree:            "true\n",
		verifyArgs("main"):     "m\n",
		"merge-base HEAD main": "abc\n",
		lsTreeArgs + "abc" + logSpec: base(edited, "then\n") + base(gone, "z\n") + base(moved, "same\n") +
			base(".rhiza/log/.gitkeep", ""),
		lsTreeArgs + "HEAD~1" + logSpec: "",
		verifyArgs("HEAD~1"):            "h\n",
		lsTreeArgs + "broken" + logSpec: "",
		verifyArgs("broken"):            "b\n",
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
	var got []CheckFinding
	for _, f := range r.Findings {
		if f.Rule == check.RuleAppendOnly {
			got = append(got, f)
		}
	}
	// The moved entry kept its bytes, so only the edit and the deletion
	// fail; the missing .gitkeep is not an entry.
	want := []CheckFinding{
		{Severity: "fail", Rule: "append-only", Path: edited, Entry: "y",
			Message: "log entry y modified since the base; entries may only be added",
			Next:    "git checkout abc -- " + edited},
		{Severity: "fail", Rule: "append-only", Path: gone, Entry: "z",
			Message: "log entry z deleted since the base; entries may only be added",
			Next:    "git checkout abc -- " + gone},
	}
	if !slices.Equal(got, want) {
		t.Errorf("append-only findings\n got %+v\nwant %+v", got, want)
	}
	if r.OK || r.ExitCode() != 1 {
		t.Errorf("ok %v exit %d", r.OK, r.ExitCode())
	}
	out := human(t, r)
	for _, s := range []string{"failed\n", "base: abc (merge-base of HEAD and main)\n",
		"fail append-only .rhiza/log/x/y.json: log entry y modified since the base",
		"  next: git checkout abc -- .rhiza/log/x/y.json\n"} {
		if !strings.Contains(out, s) {
			t.Errorf("human lacks %q:\n%s", s, out)
		}
	}

	// A listing that fails is an error, not a silent pass.
	delete(g.out, lsTreeArgs+"broken"+logSpec)
	_, err = callOp(t, ctx, "check", `{"base":"broken"}`)
	wantExit(t, err, 1)

	_, err = callOp(t, ctx, "check", `{"base":"nope"}`)
	wantExit(t, err, 2)
}

// A base from before the rename keeps its log at .genesis/log: an entry
// moved to .rhiza/log with its bytes passes, an edited or deleted one fails
// with a Next that restores it and moves it under .rhiza/.
func TestCheckInputPreRename(t *testing.T) {
	ctx := initCtx(t)
	writeEntry(t, ctx, ".rhiza/log/x/kept.json", "same\n")
	writeEntry(t, ctx, ".rhiza/log/x/edited.json", "now\n")
	base := func(p, b string) string { return "100644 blob " + gitBlob([]byte(b)) + "\t" + p + "\x00" }
	ctx.Git = (&fakeGit{out: map[string]string{
		gitWorkTree:       "true\n",
		verifyArgs("old"): "o\n",
		lsTreeArgs + "old" + logSpec: base(".genesis/log/x/kept.json", "same\n") +
			base(".genesis/log/x/edited.json", "then\n") + base(".genesis/log/x/gone.json", "z\n"),
	}}).run
	resp, err := callOp(t, ctx, "check", `{"base":"old"}`)
	if err != nil {
		t.Fatal(err)
	}
	var got []CheckFinding
	for _, f := range resp.(CheckResponse).Findings {
		if f.Rule == check.RuleAppendOnly {
			got = append(got, f)
		}
	}
	want := []CheckFinding{
		{Severity: "fail", Rule: "append-only", Path: ".genesis/log/x/gone.json", Entry: "gone",
			Message: "log entry gone deleted since the base; entries may only be added",
			Next:    "git checkout old -- .genesis/log/x/gone.json && git mv .genesis/log/x/gone.json .rhiza/log/x/gone.json"},
		{Severity: "fail", Rule: "append-only", Path: ".rhiza/log/x/edited.json", Entry: "edited",
			Message: "log entry edited modified since the base, where it was .genesis/log/x/edited.json; entries may only be added",
			Next:    "git checkout old -- .genesis/log/x/edited.json && git mv -f .genesis/log/x/edited.json .rhiza/log/x/edited.json"},
	}
	if !slices.Equal(got, want) {
		t.Errorf("append-only findings\n got %+v\nwant %+v", got, want)
	}
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
		"fail render .rhiza/rendered: the rendered folder is missing\n" +
		"  next: rhiza render\n" +
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
		{CheckFinding{Severity: "fail", Rule: "log", Path: ".rhiza/log/a", Message: "bad"}, "fail log .rhiza/log/a: bad"},
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
	wantNext(t, err, "rhiza init")

	ctx = initCtx(t)
	ctx.DryRun = true
	_, err = callOp(t, ctx, "check", `{}`)
	wantExit(t, err, 2)
	wantNext(t, err, "rhiza check")
}

func TestParseTrailerLog(t *testing.T) {
	out := "aaa\x1fSubject\n\nChange: CS1\nActor: builder\nLeaf: leaf\nLeaf: other\x1e\n" +
		"bbb\x1fSubject only\x1e\n" +
		"ccc\x1fSubject\n\nprose\n\nchange: CS2\nCHANGE:  CS3 \nactor: 01M5ACTOR\x1e\n" +
		// Trailers in a paragraph before Co-Authored-By's.
		"ddd\x1fSubject\n\nChange: CS4\nActor: a\nLeaf: split\nReview: R1\n\nCo-Authored-By: Claude <noreply@anthropic.com>\x1e\n" +
		// Trailers sharing the final paragraph.
		"eee\x1fSubject\n\nChange: CS5\nLeaf: one\nCo-Authored-By: Claude <noreply@anthropic.com>\x1e\n" +
		// The subject is not a trailer.
		"fff\x1fLeaf: subject\n\nbody\x1e\n" +
		// A paragraph with prose in it is not trailers, even its last.
		"ggg\x1fSubject\n\nChange: CS6\nnot a trailer\x1e\n" +
		// Prose wrapped onto a "Leaf: ..." line, then real trailers.
		"hhh\x1fSubject\n\nThe builder's\nLeaf: trailer names it, and the check reads it.\n\nChange: CS7\nLeaf: real\x1e\n" +
		// A squash merge: only the trailing run, the last commit's, counts.
		"iii\x1fSquash (#53)\n\n* First\n\nWhy.\n\nChange: CS8\nActor: first\nLeaf: one\n\n" +
		"* Second\n\nChange: CS8\nActor: second\nLeaf: two\n\n---------\n\nCo-authored-by: Claude <noreply@anthropic.com>\x1e\n" +
		// A folded trailer is unfolded; blank lines may hold spaces.
		"jjj\x1fSubject\n\nprose\n \t\nChange: CS9\nLeaf: long\n  name\x1e\n"
	got := parseTrailerLog(out)
	claims := map[string][]string{}
	for i := range got {
		claims[got[i].Hash] = got[i].Claims
		got[i].Subject, got[i].Claims = "", nil
	}
	// Claims are the Change: values of every wholly-trailer paragraph after
	// the subject, where Changes are those of the trailing run only.
	if c := claims["iii"]; !reflect.DeepEqual(c, []string{"CS8", "CS8"}) {
		t.Errorf("squash claims = %v", c)
	}
	want := []check.Commit{
		{Hash: "aaa", Changes: []string{"CS1"}, Actors: []string{"builder"}, Leaves: []string{"leaf", "other"}},
		{Hash: "bbb"},
		{Hash: "ccc", Changes: []string{"CS2", "CS3"}, Actors: []string{"01M5ACTOR"}},
		{Hash: "ddd", Changes: []string{"CS4"}, Actors: []string{"a"}, Leaves: []string{"split"}, Reviews: []string{"R1"}},
		{Hash: "eee", Changes: []string{"CS5"}, Leaves: []string{"one"}},
		{Hash: "fff"},
		{Hash: "ggg"},
		{Hash: "hhh", Changes: []string{"CS7"}, Leaves: []string{"real"}},
		{Hash: "iii"},
		{Hash: "jjj", Changes: []string{"CS9"}, Leaves: []string{"long name"}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("parseTrailerLog = %+v, want %+v", got, want)
	}
	if got := parseTrailerLog(""); got == nil || len(got) != 0 {
		t.Errorf("parseTrailerLog(\"\") = %#v, want empty non-nil", got)
	}
}

func TestGitCommits(t *testing.T) {
	// Only commits that change a file outside the store folder, anywhere
	// in the repository, and on every branch merged in.
	logArgs := "log --full-history --format=" + trailerLogFormat + " HEAD -- :/ :(exclude).rhiza"
	if got := gitCommits(nil); got != nil {
		t.Errorf("no git: %+v", got)
	}
	if got := gitCommits((&fakeGit{out: map[string]string{}}).run); got != nil {
		t.Errorf("not a work tree: %+v", got)
	}
	// No commit yet: git log fails, and there are no commits, not nil.
	if got := gitCommits((&fakeGit{out: map[string]string{gitWorkTree: "true\n"}}).run); got == nil || len(got) != 0 {
		t.Errorf("no commits: %#v", got)
	}
	g := &fakeGit{out: map[string]string{gitWorkTree: "true\n", logArgs: "abc\x1fSubject\n\nActor: x\x1e\n"}}
	want := []check.Commit{{Hash: "abc", Actors: []string{"x"}, Subject: "Subject"}}
	if got := gitCommits(g.run); !reflect.DeepEqual(got, want) {
		t.Errorf("gitCommits = %+v, want %+v (calls %q)", got, want, g.calls)
	}
}

// skip_reference leaves out the generated-reference comparison and the
// docs-command rule; without it both fail.
func TestCheckSkipReference(t *testing.T) {
	ctx := initCtx(t)
	renderOK(t, ctx)
	if _, err := callOp(t, ctx, "docs", `{}`); err != nil {
		t.Fatal(err)
	}
	stale := filepath.Join(ctx.Root, filepath.FromSlash(docsDir), "cli.md")
	if err := os.WriteFile(stale, []byte(docs.Marker+"\nstale\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(ctx.Root, "README.md"), []byte("Try `rhiza no-such-command`.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	reference := func(req string) (n int) {
		resp, err := callOp(t, ctx, "check", req)
		if err != nil && resp == nil {
			t.Fatal(err)
		}
		for _, f := range resp.(CheckResponse).Findings {
			if f.Severity == check.SeverityFail && (f.Path == "docs/reference/cli.md" || f.Path == "README.md") {
				n++
			}
		}
		return n
	}
	if n := reference(`{}`); n != 2 {
		t.Errorf("default check has %d reference findings, want 2", n)
	}
	if n := reference(`{"skip_reference":true}`); n != 0 {
		t.Errorf("skip_reference check has %d reference findings, want 0", n)
	}
}
