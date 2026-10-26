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

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/check"
	"github.com/dphaener/rhiza/internal/store"
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
	Base string `json:"base,omitempty" desc:"The git revision the log may only have grown since; default the merge-base of HEAD with origin/HEAD, origin/main or main"`
}

// CheckResponse is the check report.
type CheckResponse struct {
	OK          bool           `json:"ok" desc:"True when no finding fails"`
	Base        string         `json:"base" desc:"The git revision the append-only rule compared against; empty when it was skipped"`
	MergeBaseOf string         `json:"merge_base_of,omitempty" desc:"The ref whose merge-base with HEAD is the default base; absent with --base"`
	Skipped     string         `json:"skipped,omitempty" desc:"Why the append-only rule was skipped, if it was"`
	Findings    []CheckFinding `json:"findings" desc:"Every finding, sorted by rule (log, fold, model, in-flight, link, render, append-only), then path and IDs"`
}

// CheckFinding is one thing check found.
type CheckFinding struct {
	Severity  string `json:"severity" enum:"fail,warn" desc:"fail makes the check fail; warn does not"`
	Rule      string `json:"rule" enum:"log,fold,model,in-flight,link,render,append-only" desc:"The rule that found it"`
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
	rep := check.Run(in)
	resp := CheckResponse{OK: rep.OK, Base: rep.Base, MergeBaseOf: from, Skipped: skipped,
		Findings: make([]CheckFinding, 0, len(rep.Findings))}
	for _, f := range rep.Findings {
		resp.Findings = append(resp.Findings, CheckFinding(f))
	}
	return resp, nil
}

// checkBase picks the base revision with git (nil for none): base as given
// when it resolves to a commit, else the default (see the file comment).
// from is the ref the default came from; skipped says why there is no
// base. An explicit base that does not resolve is a usage error.
func checkBase(git func(...string) (string, error), base string) (rev, from, skipped string, err error) {
	if git == nil {
		return "", "", "not a git work tree", nil
	}
	if out, err := git("rev-parse", "--is-inside-work-tree"); err != nil || strings.TrimSpace(out) != "true" {
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
