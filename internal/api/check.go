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

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/dphaener/genesis/internal/check"
	"github.com/dphaener/genesis/internal/store"
)

// baseCandidates are the refs the default base is merged with, first that
// resolves wins.
var baseCandidates = []string{"origin/HEAD", "origin/main", "main"}

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
		return nil, Usagef("check does not take --dry-run: it writes nothing").WithNext("genesis check")
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
	in.Base, in.Skipped = base, skipped
	if base != "" {
		out, err := ctx.Git("diff", "--name-status", "--no-renames", "-z", "--relative", base, "--", store.Dir+"/"+store.LogDir)
		if err != nil {
			return nil, Errorf("diff the log against %s: %v", base, err)
		}
		in.Changes = parseNameStatus(out)
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
				WithNext("genesis check --base <rev>")
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

// parseNameStatus parses `git diff --name-status -z`: a status, then a
// path, each ended by a NUL. With -z git writes paths verbatim, unquoted.
func parseNameStatus(out string) []check.Change {
	var cs []check.Change
	fields := strings.Split(out, "\x00")
	for i := 0; i+1 < len(fields); i += 2 {
		if fields[i] == "" || fields[i+1] == "" {
			break
		}
		cs = append(cs, check.Change{Status: fields[i], Path: fields[i+1]})
	}
	return cs
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
