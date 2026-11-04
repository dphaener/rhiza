package api

// This file holds health: the whole system's state on one screen, part by
// part. It reads what the other read operations read (daemon.status, the
// daemon settings, inbox, request.ls, roadmap and check) plus the free disk
// space and the main checkout, and writes nothing.

import (
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/machine"
	"github.com/dphaener/rhiza/internal/usage"
)

func init() {
	register(Operation{
		Name:     "health",
		Summary:  "Report the daemon, autonomy settings, the machine trust, free disk, the usage hold, the main checkout, inbox, queued requests, roadmap and rhiza check, each part ok, warn or bad; exit 1 when any part is bad",
		Request:  reflect.TypeFor[HealthRequest](),
		Response: reflect.TypeFor[HealthResponse](),
		Run:      runHealth,
	})
}

// The states of a health part, best first.
const (
	HealthOK   = "ok"
	HealthWarn = "warn"
	HealthBad  = "bad"
)

// The parts of the health report, in the order it lists them.
const (
	HealthDaemon   = "daemon"
	HealthAutonomy = "autonomy"
	HealthMachine  = "machine"
	HealthDisk     = "disk"
	HealthUsage    = "usage"
	HealthMain     = "main"
	HealthInbox    = "inbox"
	HealthRequests = "requests"
	HealthRoadmap  = "roadmap"
	HealthCheck    = "check"
)

// HealthRequest takes nothing.
type HealthRequest struct{}

// HealthPart is one part of the system and its state.
type HealthPart struct {
	Part    string   `json:"part" enum:"daemon,autonomy,machine,disk,usage,main,inbox,requests,roadmap,check" desc:"The part"`
	State   string   `json:"state" enum:"ok,warn,bad" desc:"ok, warn (worth a look) or bad (rhiza health exits 1)"`
	Summary string   `json:"summary" desc:"One line on the part's state"`
	Detail  []string `json:"detail" desc:"More lines on the part, such as each setting or each count; empty when none"`
	Next    string   `json:"next,omitempty" desc:"A command that would help, when one would"`
}

// HealthResponse is the report.
type HealthResponse struct {
	State string       `json:"state" enum:"ok,warn,bad" desc:"The worst state of any part"`
	Parts []HealthPart `json:"parts" desc:"Every part, in a fixed order: daemon, autonomy, machine, disk, usage, main, inbox, requests, roadmap, check"`
}

// ExitCode is 1 when a part is bad, else 0.
func (r HealthResponse) ExitCode() int {
	if r.State == HealthBad {
		return 1
	}
	return 0
}

// Human writes one line per part, its state first, then its detail and
// next command indented.
func (r HealthResponse) Human(w io.Writer) {
	for _, p := range r.Parts {
		fmt.Fprintf(w, "%-4s %-8s %s\n", p.State, p.Part, p.Summary)
		for _, d := range p.Detail {
			fmt.Fprintf(w, "     %s\n", d)
		}
		if p.Next != "" {
			fmt.Fprintf(w, "     next: %s\n", p.Next)
		}
	}
}

// healthRank orders the states, worst highest.
func healthRank(s string) int {
	switch s {
	case HealthBad:
		return 2
	case HealthWarn:
		return 1
	}
	return 0
}

func runHealth(ctx *Ctx, _ any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("health does not take --dry-run: it writes nothing").WithNext("rhiza health")
	}
	settings, serr := settingsOf(ctx, nil)
	parts := []HealthPart{
		healthDaemon(ctx),
		healthAutonomy(settings, serr),
		healthMachine(ctx),
		healthDisk(ctx, settings, serr),
		healthUsage(ctx, settings, serr),
		healthMain(ctx),
		healthInbox(ctx),
		healthRequests(ctx),
		healthRoadmap(ctx),
		healthCheck(ctx),
	}
	resp := HealthResponse{State: HealthOK, Parts: parts}
	for i := range resp.Parts {
		if resp.Parts[i].Detail == nil {
			resp.Parts[i].Detail = []string{}
		}
		if healthRank(resp.Parts[i].State) > healthRank(resp.State) {
			resp.State = resp.Parts[i].State
		}
	}
	return resp, nil
}

// healthFailed is a part that could not be read: warn, saying why.
func healthFailed(part string, err error) HealthPart {
	return HealthPart{Part: part, State: HealthWarn, Summary: "could not be read: " + err.Error()}
}

// healthDaemon is ok while a daemon runs, warn when none does.
func healthDaemon(ctx *Ctx) HealthPart {
	out, err := runDaemonStatus(ctx, DaemonStatusRequest{})
	if err != nil {
		return healthFailed(HealthDaemon, err)
	}
	s := out.(DaemonStatusResponse)
	if !s.Running {
		p := HealthPart{Part: HealthDaemon, State: HealthWarn,
			Summary: "not running: nothing starts or resumes runs on its own; start it with rhiza daemon",
			Next:    "rhiza daemon --autonomous"}
		if s.Stopped != "" {
			p.Detail = append(p.Detail, "last daemon "+s.Process+" stopped "+s.Stopped)
		}
		return p
	}
	p := HealthPart{Part: HealthDaemon, State: HealthOK,
		Summary: fmt.Sprintf("running, process %s, %d live runs", s.Process, len(s.Runs))}
	if s.Commit != "" {
		p.Detail = append(p.Detail, "built from "+s.Commit)
		if origin := healthOriginAhead(ctx, s.Commit); origin != "" {
			p.State = HealthWarn
			p.Summary += fmt.Sprintf("; built from %s, origin/main is at %s", shortSHA(s.Commit), shortSHA(origin))
			p.Next = "make update"
		}
	}
	if s.LastPoll != "" {
		p.Detail = append(p.Detail, "last poll "+s.LastPoll)
	}
	return p
}

// healthAutonomy lists the autonomy settings and whether the daemon is
// paused, which is warn.
func healthAutonomy(settings []Setting, err error) HealthPart {
	if err != nil {
		return HealthPart{Part: HealthAutonomy, State: HealthBad, Summary: err.Error(), Next: "rhiza daemon settings"}
	}
	p := HealthPart{Part: HealthAutonomy, State: HealthOK}
	var stages []string
	for _, s := range settings {
		if stage, ok := strings.CutPrefix(s.Key, "rhiza.autonomy."); ok {
			stages = append(stages, stage+" "+s.Value)
			p.Detail = append(p.Detail, fmt.Sprintf("%s = %s (%s)", s.Key, s.Value, s.Source))
		}
	}
	p.Summary = strings.Join(stages, ", ")
	if settingValue(settings, "rhiza.daemon.paused") == "true" {
		p.State = HealthWarn
		p.Summary = "daemon paused; " + p.Summary
		p.Next = "rhiza daemon set rhiza.daemon.paused false"
	}
	return p
}

// healthMachine is ok with the machine's trust and its file, daily when
// there is no machine config, and warn naming the file when it cannot be
// read.
func healthMachine(ctx *Ctx) HealthPart {
	path := machine.Path(ctxGetenv(ctx), ctx.Home)
	c, err := machine.Read(path)
	if err != nil {
		p := HealthPart{Part: HealthMachine, State: HealthWarn, Summary: "could not be read: " + err.Error()}
		if path != "" {
			p.Detail = append(p.Detail, "machine config "+path)
		}
		return p
	}
	if !machineFileExists(path) {
		return HealthPart{Part: HealthMachine, State: HealthOK, Summary: "trust daily (no machine config)",
			Detail: []string{"it would be " + path}}
	}
	return HealthPart{Part: HealthMachine, State: HealthOK, Summary: "trust " + c.Trust + " (from " + path + ")"}
}

// healthDisk is bad when free space is under the daemon's minimum, so the
// daemon holds new starts and resumes.
func healthDisk(ctx *Ctx, settings []Setting, err error) HealthPart {
	if err != nil {
		return healthFailed(HealthDisk, err)
	}
	free, ferr := repoFreeBytes(ctx)
	if ferr != nil {
		return healthFailed(HealthDisk, ferr)
	}
	min, _ := daemon.ParseSize(settingValue(settings, "rhiza.daemon.minFree"))
	if min == 0 {
		return HealthPart{Part: HealthDisk, State: HealthOK, Summary: daemon.FormatSize(free) + " free; the minimum is off"}
	}
	if free < min {
		return HealthPart{Part: HealthDisk, State: HealthBad,
			Summary: fmt.Sprintf("%s free, under the minimum of %s: the daemon holds new starts and resumes", daemon.FormatSize(free), daemon.FormatSize(min)),
			Next:    "rhiza daemon set rhiza.daemon.minFree <size>"}
	}
	return HealthPart{Part: HealthDisk, State: HealthOK,
		Summary: fmt.Sprintf("%s free, minimum %s", daemon.FormatSize(free), daemon.FormatSize(min))}
}

// healthUsage is warn while the daemon holds its launches for usage, else ok
// with each window's use and limit.
func healthUsage(ctx *Ctx, settings []Setting, err error) HealthPart {
	if err != nil {
		return healthFailed(HealthUsage, err)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return healthFailed(HealthUsage, err)
	}
	st, _ := daemon.ReadState(common)
	if h := st.UsageHold; h != nil {
		return HealthPart{Part: HealthUsage, State: HealthWarn,
			Summary: "the daemon holds new launches for usage until the window resets",
			Detail:  []string{usageHoldDetail(*h)},
			Next:    "rhiza daemon set rhiza.daemon.usage.reserve <percent>"}
	}
	windows, _ := usageWindowsNow(ctx, common, st.UsageWindows)
	var parts []string
	for _, kind := range usage.Kinds {
		used := windows[kind].Used
		if limit, ok := usageLimit(settings, st.Learned, kind); ok {
			parts = append(parts, fmt.Sprintf("%s %d of %d", kind, used, limit))
		} else {
			parts = append(parts, fmt.Sprintf("%s %d, no limit known", kind, used))
		}
	}
	return HealthPart{Part: HealthUsage, State: HealthOK, Summary: strings.Join(parts, "; ")}
}

// healthMain reads the main checkout: its branch, whether tracked files
// changed, how far it is behind origin/main as last fetched (no fetch is
// made), and the last ship close-out that could not pull. Anything but a
// clean, current main is warn.
func healthMain(ctx *Ctx) HealthPart {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return healthFailed(HealthMain, err)
	}
	out, err := ctx.Git("worktree", "list", "--porcelain")
	if err != nil {
		return healthFailed(HealthMain, err)
	}
	main, ok := firstWorktree(out)
	if !ok {
		return healthFailed(HealthMain, fmt.Errorf("git worktree list names no main checkout"))
	}
	git := func(args ...string) (string, error) {
		s, err := ctx.Git(append([]string{"-C", main}, args...)...)
		return strings.TrimSpace(s), err
	}
	p := HealthPart{Part: HealthMain, State: HealthOK}
	var facts []string
	branch, err := git("symbolic-ref", "--short", "HEAD")
	switch {
	case err != nil:
		facts = append(facts, "not on a branch")
		p.State = HealthWarn
	case branch != "main":
		facts = append(facts, "on "+branch+", not main")
		p.State = HealthWarn
	default:
		facts = append(facts, "on main")
	}
	if st, err := git("status", "--porcelain", "--untracked-files=no"); err == nil && st != "" {
		facts = append(facts, "uncommitted changes")
		p.State = HealthWarn
	} else if err == nil {
		facts = append(facts, "clean")
	}
	if _, err := git("rev-parse", "--verify", "--quiet", "refs/remotes/origin/main"); err == nil {
		if n, err := git("rev-list", "--count", "HEAD..origin/main"); err == nil {
			if behind, _ := strconv.Atoi(n); behind > 0 {
				facts = append(facts, fmt.Sprintf("%d behind origin/main", behind))
				p.State = HealthWarn
			} else {
				facts = append(facts, "up to date with origin/main")
			}
		}
	}
	p.Summary = strings.Join(facts, ", ")
	p.Detail = append(p.Detail, "checkout "+main)
	if ms, ok := daemon.ReadMainStale(common); ok {
		p.State = HealthWarn
		p.Summary += "; the last close-out could not pull"
		p.Detail = append(p.Detail, fmt.Sprintf("change set %s, since %s: %s", ms.Change, ms.Since.UTC().Format("2006-01-02T15:04:05Z"), ms.Detail))
	}
	if p.State == HealthWarn {
		p.Next = "git -C " + main + " status"
	}
	return p
}

// healthInbox counts what waits on the owner by kind; anything waiting is
// warn.
func healthInbox(ctx *Ctx) HealthPart {
	out, err := runInbox(ctx, InboxRequest{})
	if err != nil {
		return healthFailed(HealthInbox, err)
	}
	items := out.(InboxResponse).Items
	if len(items) == 0 {
		return HealthPart{Part: HealthInbox, State: HealthOK, Summary: "nothing waits on you"}
	}
	return HealthPart{Part: HealthInbox, State: HealthWarn,
		Summary: fmt.Sprintf("%d waiting: %s", len(items), healthCounts(items, func(it InboxItem) string { return it.Kind })),
		Next:    "rhiza inbox"}
}

// healthRequests counts the plan requests by state; a failed one is warn.
func healthRequests(ctx *Ctx) HealthPart {
	out, err := runRequestLs(ctx, RequestLsRequest{})
	if err != nil {
		return healthFailed(HealthRequests, err)
	}
	qs := out.(RequestLsResponse).Requests
	if len(qs) == 0 {
		return HealthPart{Part: HealthRequests, State: HealthOK, Summary: "no plan requests"}
	}
	p := HealthPart{Part: HealthRequests, State: HealthOK,
		Summary: fmt.Sprintf("%d requests: %s", len(qs), healthCounts(qs, func(q RequestInfo) string { return q.State }))}
	if slices.ContainsFunc(qs, func(q RequestInfo) bool { return q.State == "failed" }) {
		p.State, p.Next = HealthWarn, "rhiza request ls"
	}
	return p
}

// healthRoadmap counts the milestones and the in-flight change sets by
// state, and how many are ready.
func healthRoadmap(ctx *Ctx) HealthPart {
	out, err := runRoadmap(ctx, RoadmapRequest{All: true})
	if err != nil {
		return healthFailed(HealthRoadmap, err)
	}
	var milestones, ready int
	var sets []RoadmapChangeSet
	for _, g := range out.(RoadmapResponse).Groups {
		if g.Kind == "milestone" {
			milestones++
		}
		for _, c := range g.ChangeSets {
			sets = append(sets, c)
			if c.Ready {
				ready++
			}
		}
	}
	p := HealthPart{Part: HealthRoadmap, State: HealthOK,
		Summary: fmt.Sprintf("%d milestones, %d change sets, %d ready", milestones, len(sets), ready)}
	if len(sets) > 0 {
		p.Detail = append(p.Detail, healthCounts(sets, func(c RoadmapChangeSet) string { return string(c.State) }))
	}
	return p
}

// healthCheck is rhiza check's result: bad when a finding fails, warn when
// one warns.
func healthCheck(ctx *Ctx) HealthPart {
	out, err := runCheck(ctx, CheckRequest{})
	if err != nil {
		return HealthPart{Part: HealthCheck, State: HealthBad, Summary: "check failed to run: " + err.Error(), Next: "rhiza check"}
	}
	r := out.(CheckResponse)
	var fails, warns int
	p := HealthPart{Part: HealthCheck}
	for _, f := range r.Findings {
		if f.Severity == "fail" {
			fails++
		} else {
			warns++
		}
		p.Detail = append(p.Detail, fmt.Sprintf("%s %s: %s", f.Severity, f.Rule, f.Message))
	}
	switch {
	case fails > 0:
		p.State, p.Summary, p.Next = HealthBad, fmt.Sprintf("%d failing, %d warnings", fails, warns), "rhiza check"
	case warns > 0:
		p.State, p.Summary, p.Next = HealthWarn, fmt.Sprintf("ok, %d warnings", warns), "rhiza check"
	default:
		p.State, p.Summary = HealthOK, "ok"
	}
	return p
}

// healthCounts is "2 failed, 1 question": each key's count, by count
// descending then key.
func healthCounts[T any](xs []T, key func(T) string) string {
	n := map[string]int{}
	var keys []string
	for _, x := range xs {
		k := key(x)
		if n[k] == 0 {
			keys = append(keys, k)
		}
		n[k]++
	}
	slices.SortFunc(keys, func(a, b string) int {
		if n[a] != n[b] {
			return n[b] - n[a]
		}
		return strings.Compare(a, b)
	})
	parts := make([]string, len(keys))
	for i, k := range keys {
		parts[i] = fmt.Sprintf("%d %s", n[k], k)
	}
	return strings.Join(parts, ", ")
}

// healthOriginAhead is origin/main's commit (as last fetched, nothing is
// fetched) when it is a descendant of commit and not commit itself, else "".
// A commit this repository does not know, or no origin/main, is "".
func healthOriginAhead(ctx *Ctx, commit string) string {
	out, err := ctx.Git("rev-parse", "--verify", "--quiet", "refs/remotes/origin/main")
	origin := strings.TrimSpace(out)
	if err != nil || origin == "" || origin == commit {
		return ""
	}
	if _, err := ctx.Git("merge-base", "--is-ancestor", commit, origin); err != nil {
		return ""
	}
	return origin
}
