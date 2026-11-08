package api

// This file holds what ship.check --main does when the merge commit's tree
// differs from the tested head's: main moved between the checks and the
// merge, so nothing tested what is on main. CI no longer runs on a push to
// main, so Ship dispatches it on the merge commit: every workflow at
// origin/main's tip that runs on pull_request and can be dispatched,
// pinned to the commit by a sha input its run-name shows.

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// shipWorkflowsDir is where GitHub reads workflow files.
const shipWorkflowsDir = ".github/workflows"

// dispatchWorkflow is what ship needs to know about one workflow file.
type dispatchWorkflow struct {
	File   string // base name, as gh workflow run takes it
	Lacks  bool   // dispatchable, but cannot be pinned to or found by a sha
	Reason string
}

// shipCheckMainDispatch runs CI on the merge commit and checks its runs
// the way a pull request's checks are: settled once every run finished
// (settled asks only that), green once they all passed, a failed run
// rerun once. It stops, never waits, when nothing can be dispatched.
func shipCheckMainDispatch(ctx *Ctx, t shipTarget, st shipState, resp ShipCheckResponse, path string, settled bool, budget int) (any, error) {
	mc := st.MergeCommit
	dir := t.dir()
	stop := func(format string, a ...any) (any, error) {
		resp.Stop = fmt.Sprintf("merge commit %s: %s", shortSHA(mc), fmt.Sprintf(format, a...))
		resp.Detail = resp.Stop
		return resp, nil
	}
	wfs, err := shipDispatchWorkflows(ctx, dir)
	if err != nil {
		return nil, err
	}
	if len(wfs) == 0 {
		return stop("its tree differs from the tested head %s, and no dispatchable workflow was found at origin/%s (one declaring pull_request and workflow_dispatch under %s)", shortSHA(st.Head), shipMain, shipWorkflowsDir)
	}
	for _, w := range wfs {
		if w.Lacks {
			return stop("workflow %s %s; it cannot be dispatched on the merge commit", w.File, w.Reason)
		}
	}
	if st.DispatchCommit != mc {
		st.DispatchCommit, st.Dispatched = mc, nil
	}
	for _, w := range wfs {
		if slices.Contains(st.Dispatched, w.File) {
			continue
		}
		// Record before dispatching: a crash between the two then costs a
		// missed dispatch, never a second one.
		st.Dispatched = append(st.Dispatched, w.File)
		if err := writeShipState(path, st); err != nil {
			return nil, err
		}
		if _, err := gh(ctx, dir, "workflow", "run", w.File, "--ref", shipMain, "-f", "sha="+mc); err != nil {
			return nil, Errorf("%v", err)
		}
	}
	var checks []ShipCheck
	var waiting []string
	for _, w := range wfs {
		out, err := gh(ctx, dir, "run", "list", "--workflow", w.File, "--event", "workflow_dispatch", "--json", "databaseId,name,displayTitle,status,conclusion", "--limit", "100")
		if err != nil {
			return nil, Errorf("%v", err)
		}
		found, err := parseDispatchedRuns(out, mc)
		if err != nil {
			return nil, Errorf("%v", err)
		}
		if len(found) == 0 {
			waiting = append(waiting, w.File)
		}
		checks = append(checks, found...)
	}
	resp.Checks = checks
	if len(waiting) > 0 {
		resp.Detail = fmt.Sprintf("not settled: no run of %s on %s yet", strings.Join(waiting, ", "), shortSHA(mc))
		return resp, nil
	}
	resp.Settled, resp.Green, resp.Failed = checksVerdict(checks)
	if resp.Settled && !resp.Green {
		if err := shipRerunFailed(ctx, t, &st, &resp, path, budget); err != nil {
			return nil, err
		}
	}
	switch {
	case !resp.Settled:
		if resp.Detail == "" {
			resp.Detail = notSettled(checks)
		}
		return resp, nil
	case !resp.Green && settled:
		// A wait passes once everything finished; main-green stops.
		resp.Passed = true
		resp.Detail = "failed: " + strings.Join(resp.Failed, ", ")
		return resp, nil
	case !resp.Green:
		return stop("CI failed on main: %s", strings.Join(resp.Failed, ", "))
	}
	resp.Passed = true
	if settled {
		return resp, nil
	}
	st.Main = checks
	if err := writeShipState(path, st); err != nil {
		return nil, err
	}
	resp.Written = true
	return resp, nil
}

// shipDispatchWorkflows reads the workflows at origin/main's tip that
// declare both pull_request and workflow_dispatch.
func shipDispatchWorkflows(ctx *Ctx, dir string) ([]dispatchWorkflow, error) {
	ref := "refs/remotes/origin/" + shipMain
	out, err := gitAt(ctx, dir, "ls-tree", "--name-only", ref, shipWorkflowsDir+"/")
	if err != nil {
		return nil, Errorf("%v", err)
	}
	var wfs []dispatchWorkflow
	for _, p := range strings.Fields(out) {
		if !strings.HasSuffix(p, ".yml") && !strings.HasSuffix(p, ".yaml") {
			continue
		}
		text, err := gitAt(ctx, dir, "show", ref+":"+p)
		if err != nil {
			return nil, Errorf("%v", err)
		}
		info := scanWorkflow(text)
		if !info.pullRequest || !info.dispatch {
			continue
		}
		w := dispatchWorkflow{File: p[strings.LastIndex(p, "/")+1:]}
		switch {
		case !info.shaInput:
			w.Lacks, w.Reason = true, "lacks the sha input on its workflow_dispatch"
		case !strings.Contains(info.runName, "inputs.sha"):
			w.Lacks, w.Reason = true, "lacks a run-name referencing inputs.sha"
		}
		wfs = append(wfs, w)
	}
	return wfs, nil
}

// workflowInfo is what scanWorkflow reads out of a workflow file.
type workflowInfo struct {
	pullRequest, dispatch, shaInput bool
	runName                         string
}

var (
	yamlKeyRe  = regexp.MustCompile(`^(\s*)["']?([A-Za-z0-9_-]+)["']?\s*:\s*(.*)$`)
	yamlWordRe = regexp.MustCompile(`[A-Za-z0-9_-]+`)
	flowShaRe  = regexp.MustCompile(`\bsha\s*:`)
)

// scanWorkflow reads the triggers, the dispatch inputs and the run-name of
// a workflow file with no YAML parser: top-level keys by indentation, the
// on: block as a scalar, a flow list or map, or a block map.
func scanWorkflow(text string) workflowInfo {
	var info workflowInfo
	lines := strings.Split(strings.ReplaceAll(text, "\r", ""), "\n")
	indent := func(l string) int { return len(l) - len(strings.TrimLeft(l, " ")) }
	blank := func(l string) bool {
		s := strings.TrimSpace(l)
		return s == "" || strings.HasPrefix(s, "#")
	}
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		if blank(l) || indent(l) != 0 {
			continue
		}
		m := yamlKeyRe.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		val := strings.TrimSpace(stripYAMLComment(m[3]))
		switch m[2] {
		case "run-name":
			info.runName = val
		case "on", "true":
			// The block is the following deeper lines.
			end := i + 1
			for end < len(lines) && (blank(lines[end]) || indent(lines[end]) > 0) {
				end++
			}
			block := lines[i+1 : end]
			if val != "" {
				for _, w := range yamlWordRe.FindAllString(val, -1) {
					info.mark(w)
				}
				info.shaInput = info.shaInput || (info.dispatch && flowShaRe.MatchString(val))
				continue
			}
			scanOnBlock(block, indent, blank, &info)
		}
	}
	return info
}

// mark records a trigger name.
func (w *workflowInfo) mark(name string) {
	switch name {
	case "pull_request":
		w.pullRequest = true
	case "workflow_dispatch":
		w.dispatch = true
	}
}

// scanOnBlock reads the children of an on: block map.
func scanOnBlock(block []string, indent func(string) int, blank func(string) bool, info *workflowInfo) {
	child := -1
	for i, l := range block {
		if blank(l) {
			continue
		}
		if child < 0 {
			child = indent(l)
		}
		if indent(l) != child {
			continue
		}
		s := strings.TrimSpace(l)
		if strings.HasPrefix(s, "- ") { // on: a list of names
			info.mark(strings.TrimSpace(strings.TrimPrefix(s, "- ")))
			continue
		}
		m := yamlKeyRe.FindStringSubmatch(l)
		if m == nil {
			continue
		}
		info.mark(m[2])
		if m[2] != "workflow_dispatch" {
			continue
		}
		if flowShaRe.MatchString(stripYAMLComment(m[3])) {
			info.shaInput = true
		}
		// Its own block: the lines deeper than the key, looking for an
		// inputs: map with a sha key.
		inInputs, inputsIndent, keyIndent := false, 0, -1
		for _, d := range block[i+1:] {
			if blank(d) {
				continue
			}
			if indent(d) <= child {
				break
			}
			dm := yamlKeyRe.FindStringSubmatch(d)
			if dm == nil {
				continue
			}
			if inInputs && indent(d) <= inputsIndent {
				inInputs = false
			}
			if dm[2] == "inputs" && !inInputs {
				inInputs, inputsIndent, keyIndent = true, indent(d), -1
				if flowShaRe.MatchString(stripYAMLComment(dm[3])) {
					info.shaInput = true
				}
				continue
			}
			if inInputs {
				if keyIndent < 0 {
					keyIndent = indent(d)
				}
				if indent(d) == keyIndent && dm[2] == "sha" {
					info.shaInput = true
				}
			}
		}
	}
}

// stripYAMLComment drops a trailing " # comment".
func stripYAMLComment(s string) string {
	if i := strings.Index(s, " #"); i >= 0 {
		return s[:i]
	}
	return s
}

// parseDispatchedRuns parses `gh run list --json
// databaseId,name,displayTitle,status,conclusion` and keeps the runs
// whose display title holds the commit: the dispatched runs on it. A
// run's check is named by its workflow's name.
func parseDispatchedRuns(b []byte, commit string) ([]ShipCheck, error) {
	var raw []struct {
		DatabaseID   int64  `json:"databaseId"`
		Name         string `json:"name"`
		DisplayTitle string `json:"displayTitle"`
		Status       string `json:"status"`
		Conclusion   string `json:"conclusion"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return nil, fmt.Errorf("gh run list: %v", err)
	}
	var out []ShipCheck
	for _, r := range raw {
		if !strings.Contains(r.DisplayTitle, commit) {
			continue
		}
		c := ShipCheck{Name: r.Name, Status: r.Status, Conclusion: r.Conclusion}
		if r.DatabaseID != 0 {
			c.RunID = fmt.Sprint(r.DatabaseID)
		}
		out = append(out, c)
	}
	return out, nil
}
