package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Soak is how long the root daemon must have run a commit before it ships.
const Soak = 24 * time.Hour

// PatchSpacing is the least time between two patch releases.
const PatchSpacing = 24 * time.Hour

// Run history file: JSON lines, one per commit the root daemon ran, as
// {"commit":"<sha>","at":"<RFC3339 of the first run>"}. A commit may appear
// more than once; the earliest line counts.

// RunHistory maps a commit to when the root daemon first ran it.
type RunHistory map[string]time.Time

// ReadRunHistory reads the history file at path. A missing file is an empty
// history, so every commit then fails the soak gate.
func ReadRunHistory(path string) (RunHistory, error) {
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		return RunHistory{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	h := RunHistory{}
	sc := bufio.NewScanner(f)
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var e struct {
			Commit string `json:"commit"`
			At     string `json:"at"`
		}
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("%s:%d: %w", path, n, err)
		}
		at, err := time.Parse(time.RFC3339, e.At)
		if err != nil || e.Commit == "" {
			return nil, fmt.Errorf("%s:%d: want commit and an RFC3339 at", path, n)
		}
		if old, ok := h[e.Commit]; !ok || at.Before(old) {
			h[e.Commit] = at
		}
	}
	return h, sc.Err()
}

// GateInputs is everything the gates read about the candidate commit.
type GateInputs struct {
	Now         time.Time
	Commit      string
	Kind        Kind
	Severe      bool        // a severe: fix is in the release
	CI          string      // CI conclusion on Commit: "success" passes
	Runs        RunHistory  // when the root daemon first ran each commit
	BugRequests []time.Time // when bug requests were filed
	InFlight    []string    // change sets serving a finished milestone not landed in Commit
	Blocking    []string    // open questions blocking one of them
	LastPatch   time.Time   // zero when no patch was released
}

// Gate is one release gate's verdict and why.
type Gate struct {
	Name   string `json:"name"`
	Pass   bool   `json:"pass"`
	Reason string `json:"reason"`
}

// Gates evaluates every gate for the candidate commit. A severe fix skips the
// soak and the patch batching, never CI.
func Gates(in GateInputs) []Gate {
	return []Gate{gateCI(in), gateSoak(in), gateMilestone(in), gatePatch(in)}
}

// AllPass reports whether every gate passes.
func AllPass(gs []Gate) bool {
	for _, g := range gs {
		if !g.Pass {
			return false
		}
	}
	return true
}

func gateCI(in GateInputs) Gate {
	g := Gate{Name: "ci"}
	switch in.CI {
	case "success":
		g.Pass, g.Reason = true, "CI is green on "+in.Commit
	case "":
		g.Reason = "no CI status for " + in.Commit
	default:
		g.Reason = fmt.Sprintf("CI on %s is %s", in.Commit, in.CI)
	}
	return g
}

func gateSoak(in GateInputs) Gate {
	g := Gate{Name: "soak"}
	if in.Severe {
		g.Pass, g.Reason = true, "skipped: a severe fix is in the release"
		return g
	}
	first, ok := in.Runs[in.Commit]
	if !ok {
		g.Reason = in.Commit + " is not in the daemon run history"
		return g
	}
	if age := in.Now.Sub(first); age < Soak {
		g.Reason = fmt.Sprintf("the daemon first ran %s %s ago, under %s", in.Commit, age.Round(time.Minute), Soak)
		return g
	}
	for _, b := range in.BugRequests {
		if b.After(first) {
			g.Reason = "a bug request was filed at " + b.UTC().Format(time.RFC3339) + ", after the soak began"
			return g
		}
	}
	g.Pass, g.Reason = true, "the daemon has run "+in.Commit+" for over 24h with no bug request since"
	return g
}

func gateMilestone(in GateInputs) Gate {
	g := Gate{Name: "milestone"}
	switch {
	case len(in.InFlight) > 0:
		g.Reason = "change sets not landed in " + in.Commit + ": " + strings.Join(in.InFlight, ", ")
	case len(in.Blocking) > 0:
		g.Reason = "questions blocking the milestone: " + strings.Join(in.Blocking, ", ")
	default:
		g.Pass, g.Reason = true, "no change set in flight and no question blocking the milestone"
	}
	return g
}

func gatePatch(in GateInputs) Gate {
	g := Gate{Name: "patch", Pass: true}
	switch {
	case in.Kind != KindPatch:
		g.Reason = "not a patch release"
	case in.Severe:
		g.Reason = "skipped: a severe fix is released at once"
	case !in.LastPatch.IsZero() && in.Now.Sub(in.LastPatch) < PatchSpacing:
		g.Pass = false
		g.Reason = fmt.Sprintf("the last patch was %s ago, under %s", in.Now.Sub(in.LastPatch).Round(time.Minute), PatchSpacing)
	default:
		g.Reason = "no patch in the last 24h"
	}
	return g
}

// CheckEnv is what tools/release check reads from outside the repo: the
// clock, the CI conclusion on a commit and the daemon run history file.
type CheckEnv struct {
	Now    func() time.Time
	CI     func(commit string) (string, error)
	Runs   string // daemon run history file; "" is an empty history
	Commit string // candidate commit; "" is HEAD
	// Breaking detects the breaking differences since the last tag; nil is
	// breakingAgainstLastTag, which builds both versions.
	Breaking func(dir string) (string, []Break, error)
	// Release is the state of a version's GitHub release (ReleaseNone,
	// ReleaseDraft or ReleasePublished) and Draft prepares a draft one: how a
	// major release is signed off. Only a major release uses them.
	Release func(tag string) (string, error)
	Draft   func(tag, commit, notes string) error
}

// parseCheck reads check's flags. The run history defaults to
// $RHIZA_RELEASE_RUNS.
func parseCheck(args []string) (CheckEnv, error) {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	env := CheckEnv{Now: time.Now, CI: ghCI, Release: ghRelease, Draft: ghDraft}
	fs.StringVar(&env.Runs, "runs", os.Getenv("RHIZA_RELEASE_RUNS"), "daemon run history file (JSON lines of commit and at)")
	fs.StringVar(&env.Commit, "commit", "", "candidate commit (default: HEAD)")
	return env, fs.Parse(args)
}

// ghCI is the conclusion of the latest CI workflow run on commit, through
// gh: "success", "failure" and the like, "pending" while it runs, "" for
// no run.
func ghCI(commit string) (string, error) {
	out, err := exec.Command("gh", "run", "list", "--workflow", "ci.yml", "--commit", commit,
		"--limit", "1", "--json", "status,conclusion").Output()
	if err != nil {
		return "", fmt.Errorf("gh run list: %w", err)
	}
	var runs []struct{ Status, Conclusion string }
	if err := json.Unmarshal(out, &runs); err != nil {
		return "", err
	}
	if len(runs) == 0 {
		return "", nil
	}
	if runs[0].Status != "completed" {
		return "pending", nil
	}
	return runs[0].Conclusion, nil
}

// gateInputs gathers the gates' inputs for the candidate commit of the repo
// at dir: the release's kind, its severe fixes, the milestones it finishes
// and what is in flight or blocking them, bug requests, the last patch tag,
// the CI conclusion and the daemon run history.
func gateInputs(dir, logDir string, env CheckEnv, h History, in Inputs, kind Kind, marks map[string]Marks) (GateInputs, error) {
	commit := env.Commit
	if commit == "" {
		c, err := gitOut(dir, "rev-parse", "HEAD")
		if err != nil {
			return GateInputs{}, fmt.Errorf("rev-parse HEAD: %w", err)
		}
		commit = c
	}
	g := GateInputs{Now: env.Now(), Commit: commit, Kind: kind, BugRequests: h.Bugs, Runs: RunHistory{}}
	for _, f := range in.Fixes {
		if marks[f].IsSevere() {
			g.Severe = true
		}
	}
	if env.Runs != "" {
		runs, err := ReadRunHistory(env.Runs)
		if err != nil {
			return g, err
		}
		g.Runs = runs
	}
	ci, err := env.CI(commit)
	if err != nil {
		return g, err
	}
	g.CI = ci
	q, err := ReadQuestions(logDir)
	if err != nil {
		return g, err
	}
	var serving []string
	for _, m := range in.Milestone {
		serving = append(serving, h.Milestones[m]...)
	}
	sort.Strings(serving)
	inCommit, err := logFiles(dir, commit)
	if err != nil {
		return g, err
	}
	for _, cs := range serving {
		if l, ok := h.Landed[cs]; !ok || !inCommit[l.Entry] {
			g.InFlight = append(g.InFlight, cs)
		}
	}
	g.Blocking = q.Blocking(serving)
	g.LastPatch = lastPatch(dir)
	return g, nil
}

// logFiles lists the log entries in commit's tree, as paths under
// .rhiza/log. A milestone finished in the log read from disk is still in
// flight for a candidate commit that lacks one of its land entries: an
// uncommitted land, or an older --commit.
func logFiles(dir, commit string) (map[string]bool, error) {
	out, err := gitOut(dir, "ls-tree", "-r", "--name-only", "--full-tree", commit, "--", ".rhiza/log")
	if err != nil {
		return nil, fmt.Errorf("ls-tree %s: %w", commit, err)
	}
	files := map[string]bool{}
	for _, l := range strings.Split(out, "\n") {
		if rel, ok := strings.CutPrefix(l, ".rhiza/log/"); ok {
			files[rel] = true
		}
	}
	return files, nil
}

// lastPatch is when the latest vX.Y.Z tag with Z above 0 was made; zero for
// none.
func lastPatch(dir string) time.Time {
	out, err := gitOut(dir, "tag", "--list", "v[0-9]*.[0-9]*.[0-9]*", "--format=%(refname:short) %(creatordate:iso-strict)")
	if err != nil {
		return time.Time{}
	}
	var last time.Time
	for _, line := range strings.Split(out, "\n") {
		name, date, ok := strings.Cut(line, " ")
		m := tagRE.FindStringSubmatch(name)
		if !ok || m == nil || m[3] == "0" {
			continue
		}
		if at, err := time.Parse(time.RFC3339, date); err == nil && at.After(last) {
			last = at
		}
	}
	return last
}

func printGates(w io.Writer, commit string, gs []Gate) {
	verdict := "blocked"
	if AllPass(gs) {
		verdict = "ready"
	}
	fmt.Fprintf(w, "\ngates for %s: %s\n", commit, verdict)
	for _, g := range gs {
		mark := "pass"
		if !g.Pass {
			mark = "FAIL"
		}
		fmt.Fprintf(w, "  %-9s %s  %s\n", g.Name, mark, g.Reason)
	}
}

// Questions is what the log says about open questions: a question is open
// until something answers it or it is retired.
type Questions struct {
	Open    map[string]string   // open question node -> handle
	Blocks  map[string][]string // open question node -> nodes or change sets it blocks
	Created map[string]string   // node -> the change set that created it
}

// ReadQuestions reads the open questions and what they block from the log
// under logDir. Retired links count for nothing, and neither do ops a
// body.drop retracts.
func ReadQuestions(logDir string) (Questions, error) {
	type rawOp struct {
		ID, Op, Kind, Node, Handle, Link, From, To, Type string
	}
	files, err := filepath.Glob(filepath.Join(logDir, "*", "*.json"))
	if err != nil {
		return Questions{}, err
	}
	sort.Strings(files)
	type opIn struct {
		cs string
		o  rawOp
	}
	var ops []opIn
	dropped := map[string]bool{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return Questions{}, err
		}
		var e struct {
			ChangeSet string `json:"change_set"`
			Body      struct {
				Ops  []rawOp  `json:"ops"`
				Drop []string `json:"drop"`
			} `json:"body"`
		}
		if err := json.Unmarshal(b, &e); err != nil {
			return Questions{}, fmt.Errorf("%s: %w", f, err)
		}
		for _, o := range e.Body.Ops {
			ops = append(ops, opIn{e.ChangeSet, o})
		}
		for _, d := range e.Body.Drop {
			dropped[d] = true
		}
	}
	q := Questions{Open: map[string]string{}, Blocks: map[string][]string{}, Created: map[string]string{}}
	retired := map[string]bool{}
	var links []rawOp
	for _, x := range ops {
		o := x.o
		if o.ID != "" && dropped[o.ID] {
			continue
		}
		switch o.Op {
		case "node.create":
			q.Created[o.Node] = x.cs
			if o.Kind == "question" {
				q.Open[o.Node] = o.Handle
			}
		case "link.create":
			links = append(links, o)
		case "node.retire":
			retired[o.Node] = true
		case "link.retire":
			retired[o.Link] = true
		}
	}
	for _, l := range links {
		if retired[l.Link] || retired[l.From] {
			continue
		}
		if l.Type == "answers" {
			delete(q.Open, l.To)
		}
	}
	for n := range q.Open {
		if retired[n] {
			delete(q.Open, n)
		}
	}
	for _, l := range links {
		if l.Type == "blocks" && !retired[l.Link] && !retired[l.To] {
			if _, ok := q.Open[l.From]; ok {
				q.Blocks[l.From] = append(q.Blocks[l.From], l.To)
			}
		}
	}
	return q, nil
}

// Blocking lists the handles of the open questions that block one of the
// change sets, or a node one of them created.
func (q Questions) Blocking(sets []string) []string {
	in := map[string]bool{}
	for _, s := range sets {
		in[s] = true
	}
	var out []string
	for n, targets := range q.Blocks {
		for _, t := range targets {
			if in[t] || in[q.Created[t]] {
				out = append(out, q.Open[n])
				break
			}
		}
	}
	sort.Strings(out)
	return out
}
