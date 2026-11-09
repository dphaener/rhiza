package api

// This file holds the check operation: it gathers check.Input from the
// repository, the rendered folder and git, runs check.Run and reports every
// finding. It reads the log through the store alone, never the workspace.
//
// The base for the append-only rule is --base when given (a revision that
// does not resolve to a commit is a usage error), else the merge-base of
// HEAD with the first of origin/HEAD, origin/main and main that resolves and
// shares history with HEAD. The default base is reported as that merge-base
// commit, with the ref it came from in MergeBaseOf. No git, no work tree or
// no base skips the rule with a warning.
//
// The append-only rule identifies a log entry by its ID and bytes, wherever
// it sits, so check gathers the entries on both sides rather than diffing
// paths. The base side is one `git ls-tree -r -z` of the base revision over
// the store's log folder and the legacy one (see legacyDir), giving each
// entry's path and blob ID; the current side is the work tree's log folder
// (as `git diff <base>` would compare, not HEAD), each file's blob ID
// computed in Go. Neither side uses git hash-object, which resolves paths
// against the git toplevel and so breaks a project below it. Both sides
// keep only paths shaped <log>/<change set>/<entry>.json where neither the
// change set folder nor the file name starts with ".", so log/.gitkeep is
// not an entry and deleting it passes the rule. Comparing blob IDs assumes
// no content filter on log files (none in .gitattributes) and a SHA-1
// object format.
//
// For the evidence-actor rule check reads HEAD's history once with git log,
// each commit's hash and its Change:, Actor: and Leaf: trailers (see
// trailerLogFormat), of only the commits that change a file outside the
// store folder: a commit of model writes alone, such as ship's land commit
// or a verifier's Evidence commit, is not a leaf's code. Outside a work
// tree, or with no git, the rule is skipped silently; a repository with no
// commits yet has no commits.

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/dispatch"
	"github.com/dphaener/rhiza/internal/docs"
	"github.com/dphaener/rhiza/internal/store"
	"github.com/dphaener/rhiza/internal/workflow"
)

// baseCandidates are the refs the default base is merged with, first that
// resolves wins.
var baseCandidates = []string{"origin/HEAD", "origin/main", "main"}

// legacyDir is where the log lived before the rename to Rhiza. Check reads
// it from the base revision only, as pre-rename history, never from the work
// tree; it is removable once no supported base predates the rename. Until
// the rename it is the store folder itself, and the two are read once.
const legacyDir = ".genesis"

// CheckRequest names the base of the append-only rule.
type CheckRequest struct {
	Base          string `json:"base,omitempty" desc:"The git revision the log may only have grown since; default the merge-base of HEAD with origin/HEAD, origin/main or main"`
	SkipReference bool   `json:"skip_reference,omitempty" flag:"skip-reference" desc:"Leave out the generated-reference comparison (docs/reference against rhiza docs) and the docs-command rule, which judge by this binary's registry; run code-check sets it, make self-check does not"`
}

// CheckResponse is the check report.
type CheckResponse struct {
	OK          bool           `json:"ok" desc:"True when no finding fails"`
	Base        string         `json:"base" desc:"The git revision the append-only rule compared against; empty when it was skipped"`
	MergeBaseOf string         `json:"merge_base_of,omitempty" desc:"The ref whose merge-base with HEAD is the default base; absent with --base"`
	Skipped     string         `json:"skipped,omitempty" desc:"Why the append-only rule was skipped, if it was"`
	Findings    []CheckFinding `json:"findings" desc:"Every finding, sorted by rule (log, fold, model, in-flight, link, depends-on, evidence-actor, unclaimed, anchor, text-ref, dispatch, render, append-only), then path and IDs"`
}

// CheckFinding is one thing check found.
type CheckFinding struct {
	Severity  string `json:"severity" enum:"fail,warn" desc:"fail makes the check fail; warn does not"`
	Rule      string `json:"rule" enum:"log,fold,model,in-flight,link,depends-on,evidence-actor,unclaimed,anchor,text-ref,dispatch,render,append-only" desc:"The rule that found it"`
	Path      string `json:"path,omitempty" desc:"The file, relative to the project root"`
	ChangeSet string `json:"change_set,omitempty" desc:"The change set's ID"`
	Entry     string `json:"entry,omitempty" desc:"The entry's ID"`
	Op        string `json:"op,omitempty" desc:"The operation's ID"`
	Link      string `json:"link,omitempty" desc:"The link's ID"`
	Message   string `json:"message" desc:"What is wrong"`
	Next      string `json:"next,omitempty" desc:"A command that would help"`
}

func init() {
	register(Operation{
		Name:     "check",
		Summary:  "Check the log, the links and the rendered markdown",
		Request:  reflect.TypeFor[CheckRequest](),
		Response: reflect.TypeFor[CheckResponse](),
		Run:      runCheck,
	})
}

// runCheck checks the repository in ctx.Root. Check writes nothing, so
// --dry-run is a usage error, as for init.
func runCheck(ctx *Ctx, req any) (any, error) {
	r := req.(CheckRequest)
	if ctx.DryRun {
		return nil, Usagef("check does not take --dry-run: it writes nothing").WithNext("rhiza check")
	}
	st, err := openStore(ctx.Root)
	if err != nil {
		return nil, err
	}
	in := check.Input{Log: st.ReadLog()}
	rendered := filepath.Join(ctx.Root, filepath.FromSlash(renderedDir))
	if fi, err := os.Stat(rendered); err == nil && fi.IsDir() {
		in.Rendered = os.DirFS(rendered)
	}
	if !r.SkipReference {
		in.DocsWant = docsFiles()
		in.Commands = []string{}
		for _, o := range Operations() {
			in.Commands = append(in.Commands, CommandPath(o.Name))
		}
		for _, f := range FrameCommands {
			in.Commands = append(in.Commands, f.Name)
		}
		// Only a docs/reference rhiza docs already wrote into is checked: another
		// project's own docs/reference is not Rhiza's reference.
		if fsys := os.DirFS(filepath.Join(ctx.Root, filepath.FromSlash(docsDir))); docs.Generated(fsys) {
			in.Docs = fsys
		}
	}
	base, from, skipped, err := checkBase(ctx.Git, r.Base)
	if err != nil {
		return nil, err
	}
	in.Rev, in.Skipped = base, skipped
	if base != "" {
		if in.Base, err = logAt(ctx, base); err != nil {
			return nil, err
		}
		if in.Current, err = currentEntries(ctx.Root); err != nil {
			return nil, Errorf("read the log: %v", err)
		}
	}
	in.Tree = os.DirFS(ctx.Root)
	in.Commits = gitCommits(ctx.Git)
	if in.UnclaimedSince, err = unclaimedCutoff(ctx); err != nil {
		return nil, err
	}
	if in.UnclaimedSince != "" && inWorkTree(ctx.Git) {
		in.Unclaimed = gitUnclaimed(ctx.Git, in.UnclaimedSince)
	}
	in.Dispatch = dispatchProblems(ctx.Root)
	rep := check.Run(in)
	resp := CheckResponse{OK: rep.OK, Base: rep.Base, MergeBaseOf: from, Skipped: skipped,
		Findings: make([]CheckFinding, 0, len(rep.Findings))}
	for _, f := range rep.Findings {
		resp.Findings = append(resp.Findings, CheckFinding(f))
	}
	return resp, nil
}

// checkFile is the optional committed check configuration inside the store
// folder.
const checkFile = "check.json"

// unclaimedCutoff reads the unclaimed-changes cutoff from the store's
// check.json, {"unclaimed_since": "<full commit hash>"}. No file or no key
// is "" (the rule is off); unknown keys are ignored. Unreadable JSON, or a
// hash that does not resolve to a commit, is an error naming the file.
// Outside a git work tree the hash cannot be resolved and is taken as is.
func unclaimedCutoff(ctx *Ctx) (string, error) {
	path := store.Dir + "/" + checkFile
	b, err := os.ReadFile(filepath.Join(ctx.Root, filepath.FromSlash(path)))
	if errors.Is(err, fs.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", Errorf("read %s: %v", path, err)
	}
	var cfg struct {
		UnclaimedSince *string `json:"unclaimed_since"`
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return "", Errorf("%s: %v", path, err)
	}
	if cfg.UnclaimedSince == nil {
		return "", nil
	}
	h := *cfg.UnclaimedSince
	if !inWorkTree(ctx.Git) {
		return h, nil
	}
	if h == "" || strings.HasPrefix(h, "-") {
		return "", Errorf("%s: unclaimed_since %q is not a commit in this git repository", path, h)
	}
	if _, err := ctx.Git("rev-parse", "--verify", "--quiet", h+"^{commit}"); err != nil {
		return "", Errorf("%s: unclaimed_since %q is not a commit in this git repository", path, h)
	}
	return h, nil
}

// checkBase picks the base revision with git (nil for none): base as given
// when it resolves to a commit, else the default (see the file comment).
// from is the ref the default came from; skipped says why there is no
// base. An explicit base that does not resolve is a usage error.
func checkBase(git func(...string) (string, error), base string) (rev, from, skipped string, err error) {
	if !inWorkTree(git) {
		return "", "", "not a git work tree", nil
	}
	resolves := func(r string) bool {
		_, err := git("rev-parse", "--verify", "--quiet", r+"^{commit}")
		return err == nil
	}
	if base != "" {
		if !resolves(base) {
			return "", "", "", Usagef("unknown base %q: not a commit in this git repository", base).
				WithNext("rhiza check --base <rev>")
		}
		return base, "", "", nil
	}
	for _, c := range baseCandidates {
		if !resolves(c) {
			continue
		}
		out, err := git("merge-base", "HEAD", c)
		if mb := strings.TrimSpace(out); err == nil && mb != "" {
			return mb, c, "", nil
		}
	}
	return "", "", "no base: none of " + strings.Join(baseCandidates, ", ") + " resolves", nil
}

// inWorkTree reports whether git (nil for none) runs inside a work tree.
func inWorkTree(git func(...string) (string, error)) bool {
	if git == nil {
		return false
	}
	out, err := git("rev-parse", "--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

// trailerLogFormat has git log write, per commit, its hash, US (0x1f) and
// the raw message, the commit ended by RS (0x1e). git's own %(trailers)
// reads only the last paragraph, so the message is parsed by parseTrailerLog.
const trailerLogFormat = "%H%x1f%B%x1e"

// gitCommits is HEAD's history with its trailers, newest first, for the
// evidence-actor rule: the commits that change a file outside the store
// folder, anywhere in the repository (":/" is its top; the exclude is
// relative to Root, where git runs). --full-history keeps such commits on a
// merged branch that history simplification would drop. It is nil when git
// (nil for none) is not in a work tree, so the rule is skipped. When git
// log fails, as on a repository with no commits yet, there are no commits,
// and check goes on.
func gitCommits(git func(...string) (string, error)) []check.Commit {
	if !inWorkTree(git) {
		return nil
	}
	out, err := git("log", "--full-history", "--format="+trailerLogFormat, "HEAD",
		"--", ":/", ":(exclude)"+store.Dir)
	if err != nil {
		return []check.Commit{}
	}
	return parseTrailerLog(out)
}

// gitUnclaimed lists the non-merge commits of since..HEAD that change a file
// outside the store folder, with their claims, for the unclaimed rule.
// --full-history keeps such commits on a merged branch whose net change
// outside the store folder is none, which history simplification would
// drop. When git log fails there are none.
func gitUnclaimed(git func(...string) (string, error), since string) []check.Commit {
	out, err := git("log", "--no-merges", "--full-history", "--format="+trailerLogFormat, since+"..HEAD",
		"--", ":/", ":(exclude)"+store.Dir)
	if err != nil {
		return nil
	}
	return parseTrailerLog(out)
}

// parseTrailerLog parses git log output in trailerLogFormat. The Change:,
// Actor:, Leaf: and Review: trailers are read from the message's trailing
// run of trailer paragraphs (see trailerRun), so a Co-Authored-By in a later
// paragraph does not hide them, while prose and the earlier commits' trailers
// in a squash merge's body are not read. Trailer keys match
// case-insensitively, as git's do; values are trimmed and empty ones
// dropped. The result is never nil.
func parseTrailerLog(out string) []check.Commit {
	cs := []check.Commit{}
	for _, rec := range strings.Split(out, "\x1e") {
		hash, msg, _ := strings.Cut(strings.TrimSpace(rec), "\x1f")
		if hash == "" {
			continue
		}
		c := check.Commit{Hash: hash}
		c.Subject, _, _ = strings.Cut(strings.TrimSpace(msg), "\n")
		c.Subject = strings.TrimSpace(c.Subject)
		for _, p := range blankLines.Split(strings.TrimSpace(msg), -1)[1:] {
			ts, ok := trailerParagraph(p)
			if !ok {
				continue
			}
			for _, t := range ts {
				if k, v, _ := strings.Cut(t, ":"); strings.EqualFold(k, "change") && strings.TrimSpace(v) != "" {
					c.Claims = append(c.Claims, strings.TrimSpace(v))
				}
			}
		}
		for _, t := range trailerRun(msg) {
			k, v, _ := strings.Cut(t, ":")
			if v = strings.TrimSpace(v); v == "" {
				continue
			}
			switch strings.ToLower(k) {
			case "change":
				c.Changes = append(c.Changes, v)
			case "actor":
				c.Actors = append(c.Actors, v)
			case "leaf":
				c.Leaves = append(c.Leaves, v)
			case "review":
				c.Reviews = append(c.Reviews, v)
			}
		}
		cs = append(cs, c)
	}
	return cs
}

var (
	// trailerKey is a trailer line's key and colon, as git spells one.
	trailerKey = regexp.MustCompile(`^[A-Za-z0-9-]+:`)
	// blankLines separates a message's paragraphs.
	blankLines = regexp.MustCompile(`\r?\n\s*\n`)
)

// trailerRun is the trailers, unfolded, of msg's trailing run of trailer
// paragraphs: the paragraphs after the subject's, from the last one back to
// the first that is not wholly trailers. A trailer paragraph is git's: every
// line a "Key: value" trailer or, indented, the previous one's continuation.
func trailerRun(msg string) []string {
	paras := blankLines.Split(strings.TrimSpace(msg), -1)
	var run []string
	for i := len(paras) - 1; i > 0; i-- {
		ts, ok := trailerParagraph(paras[i])
		if !ok {
			break
		}
		run = append(ts, run...)
	}
	return run
}

// trailerParagraph is paragraph p's trailers, unfolded, and whether p is
// wholly trailers.
func trailerParagraph(p string) ([]string, bool) {
	var ts []string
	for _, line := range strings.Split(strings.TrimRight(p, "\r\n"), "\n") {
		switch {
		case trailerKey.MatchString(line):
			ts = append(ts, line)
		case len(ts) > 0 && strings.TrimLeft(line, " \t") != line:
			ts[len(ts)-1] += " " + strings.TrimSpace(line)
		default:
			return nil, false
		}
	}
	return ts, true
}

// logAt lists the log entries at git revision rev: one ls-tree of the
// store's log folder and the legacy one, parsed by parseLsTree.
func logAt(ctx *Ctx, rev string) ([]check.Entry, error) {
	args := []string{"ls-tree", "-r", "-z", rev, "--", store.Dir + "/" + store.LogDir}
	if legacyDir != store.Dir {
		args = append(args, legacyDir+"/"+store.LogDir)
	}
	out, err := ctx.Git(args...)
	if err != nil {
		return nil, Errorf("list the log at %s: %v", rev, err)
	}
	return parseLsTree(out), nil
}

// parseLsTree parses `git ls-tree -r -z`: per line a mode, a type and an
// object ID separated by spaces, a tab, then the path, each line ended by a
// NUL; with -z git writes paths verbatim, unquoted. It keeps the blobs whose
// path is an entry path (see entryID), sorted by path, each path once.
func parseLsTree(out string) []check.Entry {
	var es []check.Entry
	seen := map[string]bool{}
	for _, line := range strings.Split(out, "\x00") {
		meta, p, ok := strings.Cut(line, "\t")
		f := strings.Fields(meta)
		if !ok || len(f) != 3 || f[1] != "blob" || seen[p] {
			continue
		}
		id, ok := entryID(p)
		if !ok {
			continue
		}
		seen[p] = true
		es = append(es, check.Entry{ID: id, Path: p, Blob: f[2]})
	}
	slices.SortFunc(es, func(a, b check.Entry) int { return strings.Compare(a.Path, b.Path) })
	return es
}

// entryID is the entry ID of path p, relative to the project root, when p is
// <store or legacy folder>/log/<change set>/<entry>.json with neither the
// change set folder nor the file name starting with ".".
func entryID(p string) (string, bool) {
	for _, dir := range []string{store.Dir, legacyDir} {
		rest, ok := strings.CutPrefix(p, dir+"/"+store.LogDir+"/")
		if !ok {
			continue
		}
		cs, name, ok := strings.Cut(rest, "/")
		if !ok || cs == "" || strings.HasPrefix(cs, ".") || strings.HasPrefix(name, ".") ||
			strings.Contains(name, "/") || !strings.HasSuffix(name, ".json") {
			return "", false
		}
		return strings.TrimSuffix(name, ".json"), true
	}
	return "", false
}

// currentEntries reads the work tree's log under root: every regular file
// <store folder>/log/*/*.json, skipping names that start with "." as
// countEntryFiles does, with its blob ID, sorted by path. It reads the
// files, not store.ReadLog's valid entries, so an entry edited into invalid
// JSON is modified, not deleted. No log folder is no entries.
func currentEntries(root string) ([]check.Entry, error) {
	logDir := filepath.Join(root, store.Dir, store.LogDir)
	sets, err := os.ReadDir(logDir)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var es []check.Entry
	for _, cs := range sets {
		if !cs.IsDir() || strings.HasPrefix(cs.Name(), ".") {
			continue
		}
		files, err := os.ReadDir(filepath.Join(logDir, cs.Name()))
		if err != nil {
			return nil, err
		}
		for _, f := range files {
			name := f.Name()
			if !f.Type().IsRegular() || strings.HasPrefix(name, ".") || !strings.HasSuffix(name, ".json") {
				continue
			}
			b, err := os.ReadFile(filepath.Join(logDir, cs.Name(), name))
			if err != nil {
				return nil, err
			}
			es = append(es, check.Entry{ID: strings.TrimSuffix(name, ".json"),
				Path: store.Dir + "/" + store.LogDir + "/" + cs.Name() + "/" + name, Blob: gitBlob(b)})
		}
	}
	slices.SortFunc(es, func(a, b check.Entry) int { return strings.Compare(a.Path, b.Path) })
	return es, nil
}

// gitBlob is the git blob ID of b in a SHA-1 repository:
// sha1("blob " + len + NUL + bytes).
func gitBlob(b []byte) string {
	h := sha1.New()
	fmt.Fprintf(h, "blob %d\x00", len(b))
	h.Write(b)
	return hex.EncodeToString(h.Sum(nil))
}

// ExitCode is 1 when a finding fails, else 0.
func (r CheckResponse) ExitCode() int {
	if r.OK {
		return 0
	}
	return 1
}

// Human writes ok or failed, the base, then the findings, those that fail
// first, each with its next command.
func (r CheckResponse) Human(w io.Writer) {
	if r.OK {
		fmt.Fprintln(w, "ok")
	} else {
		fmt.Fprintln(w, "failed")
	}
	switch {
	case r.Skipped != "":
		fmt.Fprintf(w, "append-only: skipped (%s)\n", r.Skipped)
	case r.MergeBaseOf != "":
		fmt.Fprintf(w, "base: %s (merge-base of HEAD and %s)\n", r.Base, r.MergeBaseOf)
	default:
		fmt.Fprintf(w, "base: %s\n", r.Base)
	}
	if len(r.Findings) == 0 {
		return
	}
	fmt.Fprintln(w)
	for _, sev := range []string{check.SeverityFail, check.SeverityWarn} {
		for _, f := range r.Findings {
			if f.Severity != sev {
				continue
			}
			fmt.Fprintln(w, f.line())
			if f.Next != "" {
				fmt.Fprintf(w, "  next: %s\n", f.Next)
			}
		}
	}
}

// line is the finding on one line: severity, rule, where, message.
func (f CheckFinding) line() string {
	where := f.Path
	if where == "" {
		var ids []string
		for _, p := range []struct{ label, id string }{
			{"change set", f.ChangeSet}, {"entry", f.Entry}, {"op", f.Op}, {"link", f.Link},
		} {
			if p.id != "" {
				ids = append(ids, p.label+" "+p.id)
			}
		}
		where = strings.Join(ids, ", ")
	}
	if where == "" {
		return fmt.Sprintf("%s %s: %s", f.Severity, f.Rule, f.Message)
	}
	return fmt.Sprintf("%s %s %s: %s", f.Severity, f.Rule, where, f.Message)
}

// dispatchProblems checks root's dispatch.toml: that it parses and that each
// rule's workflow resolves in the repo layer or the built-ins. The personal
// layer is never read, so the result is the same on every machine.
func dispatchProblems(root string) []string {
	rules, err := dispatch.Load(root)
	if err != nil {
		return []string{err.Error()}
	}
	var out []string
	for _, t := range []struct {
		table string
		m     map[string]string
	}{{dispatch.TableRequest, rules.Request}, {dispatch.TableReady, rules.Ready}} {
		for _, k := range dispatch.TableKinds(t.table) {
			name := t.m[k]
			if repoOrBuiltinWorkflow(root, name) {
				continue
			}
			out = append(out, fmt.Sprintf("%s: %s.%s names workflow %q, which is neither in %s nor built in (built-in: %s)",
				dispatch.File, t.table, k, name, workflowsDir, strings.Join(workflow.BuiltinNames(), ", ")))
		}
	}
	for _, sc := range rules.Schedules {
		if sc.Workflow == "" || repoOrBuiltinWorkflow(root, sc.Workflow) {
			continue
		}
		out = append(out, fmt.Sprintf("%s: schedule %q: workflow: names workflow %q, which is neither in %s nor built in (built-in: %s)",
			dispatch.File, sc.Name, sc.Workflow, workflowsDir, strings.Join(workflow.BuiltinNames(), ", ")))
	}
	return out
}

func repoOrBuiltinWorkflow(root, name string) bool {
	if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(workflowsDir), name+".toml")); err == nil {
		return true
	}
	_, ok := workflow.Builtin(name)
	return ok
}
