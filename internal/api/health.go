package api

// This file holds health: the whole system's state on one screen, part by
// part. It reads what the other read operations read (daemon.status, the
// daemon settings, inbox, request.ls, roadmap and check) plus the free disk
// space and the main checkout, and writes nothing.

import (
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/machine"
	"github.com/dphaener/rhiza/internal/usage"
	"github.com/dphaener/rhiza/internal/workflow"
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
	HealthDaemon    = "daemon"
	HealthAutonomy  = "autonomy"
	HealthMachine   = "machine"
	HealthDisk      = "disk"
	HealthUsage     = "usage"
	HealthMain      = "main"
	HealthSchedules = "schedules"
	HealthInbox     = "inbox"
	HealthRequests  = "requests"
	HealthRoadmap   = "roadmap"
	HealthCheck     = "check"
)

// HealthRequest picks the projects the report covers.
type HealthRequest struct {
	All     bool   `json:"all,omitempty" flag:"all" desc:"Cover every registered project, from anywhere under a Rhiza root; the default at the root itself"`
	Project string `json:"project,omitempty" flag:"project" desc:"Only this registered project's parts; default the current project, or every project at the root"`
}

// HealthPart is one part of the system and its state.
type HealthPart struct {
	Project  string   `json:"project,omitempty" desc:"The registered project the part belongs to, on a report that covers every project; absent on the machine's parts and on a single project's report"`
	Part     string   `json:"part" enum:"daemon,autonomy,machine,disk,usage,main,schedules,inbox,requests,roadmap,check" desc:"The part"`
	State    string   `json:"state" enum:"ok,warn,bad" desc:"ok, warn (worth a look) or bad (rhiza health exits 1)"`
	Summary  string   `json:"summary" desc:"One line on the part's state"`
	Live     *int     `json:"live,omitempty" desc:"On the daemon part: the runs running now, summed across the daemon's projects; an interrupted run takes no slot"`
	Starting *int     `json:"starting,omitempty" desc:"On the daemon part: the starts and resumes not journaled yet, within the start grace"`
	Max      *int     `json:"max,omitempty" desc:"On the daemon part: rhiza.daemon.maxRuns as the daemon's last poll recorded it"`
	Detail   []string `json:"detail" desc:"More lines on the part, such as each setting or each count; empty when none"`
	Next     string   `json:"next,omitempty" desc:"A command that would help, when one would"`
}

// HealthResponse is the report.
type HealthResponse struct {
	State string       `json:"state" enum:"ok,warn,bad" desc:"The worst state of any part"`
	Parts []HealthPart `json:"parts" desc:"Every part, in a fixed order: daemon, autonomy, machine, disk, usage, then main, schedules (only for a project with schedules), inbox, requests, roadmap and check, once for the project or, over every project, once for each in name order tagged with its project"`
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
		sum := p.Summary
		if p.Project != "" {
			sum = p.Project + ": " + sum
		}
		fmt.Fprintf(w, "%-4s %-8s %s\n", p.State, p.Part, sum)
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

// healthProjectParts are the parts that belong to one project.
func healthProjectParts(ctx *Ctx, daemonRunning bool) []HealthPart {
	parts := []HealthPart{healthMain(ctx)}
	if p, ok := healthSchedules(ctx, daemonRunning); ok {
		parts = append(parts, p)
	}
	return append(parts,
		healthInbox(ctx),
		healthRequests(ctx),
		healthRoadmap(ctx),
		healthCheck(ctx),
	)
}

// runHealth answers health over the one project ctx works in, or, at the
// root or with --all, over every registered project: the machine's parts
// (daemon, autonomy, machine trust, disk, usage hold) once, the project's
// parts once for each project, tagged with its name. A caller passing no
// request (the web UI) reads the one project.
func runHealth(ctx *Ctx, req any) (any, error) {
	if ctx.DryRun {
		return nil, Usagef("health does not take --dry-run: it writes nothing").WithNext("rhiza health")
	}
	r, typed := req.(HealthRequest)
	scopes, merged := []inboxScope{{"", ctx}}, false
	// A root folder that holds a model of its own is read as that one
	// project, unless --project names a registered one.
	own := false
	if root, cfg, ok := rootOf(ctx); ok && r.Project == "" && projectInitialised(ctx.Wd) {
		_, inProject := cfg.ProjectOf(root, ctx.Wd)
		own = !inProject
	}
	if typed && !own {
		_, _, haveRoot := rootOf(ctx)
		var err error
		if scopes, merged, err = inboxScopes(ctx, r.Project, r.All && haveRoot); err != nil {
			return nil, err
		}
	}
	// The machine's parts read through a project's checkout (its git
	// config, the daemon's state in its git directory): over every
	// project, the root daemon's home project, whose git directory holds
	// the root daemon's state; with no project the root's own context is
	// all there is.
	mctx := ctx
	if merged && len(scopes) > 0 {
		mctx = scopes[0].ctx
		if root, cfg, ok := rootOf(ctx); ok {
			if home, ok := rootDaemonHome(ctx, root, cfg); ok {
				mctx = home
			}
		}
	}
	settings, serr := settingsOf(mctx, nil)
	parts := []HealthPart{
		healthDaemon(ctx, mctx, merged),
		healthAutonomy(mctx, settings, serr),
		healthMachine(mctx),
		healthDisk(mctx, settings, serr),
		healthUsage(mctx, settings, serr),
	}
	daemonUp := false
	if out, err := runDaemonStatus(ctx, DaemonStatusRequest{All: merged}); err == nil {
		daemonUp = out.(DaemonStatusResponse).Running
	}
	for _, sc := range scopes {
		for _, p := range healthProjectParts(sc.ctx, daemonUp) {
			if merged {
				p.Project = sc.name
			}
			parts = append(parts, p)
		}
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

// healthDaemon is ok while a daemon runs, warn when none does. Over every
// project (all) its live runs are every project's; ctx asks the daemon's
// status, home is the checkout its build is compared with origin/main in.
func healthDaemon(ctx, home *Ctx, all bool) HealthPart {
	out, err := runDaemonStatus(ctx, DaemonStatusRequest{All: all})
	if err != nil {
		return healthFailed(HealthDaemon, err)
	}
	s := out.(DaemonStatusResponse)
	byProject := s
	for _, p := range s.Projects {
		s.Runs = append(s.Runs, p.Runs...)
	}
	if !s.Running {
		p := HealthPart{Part: HealthDaemon, State: HealthWarn,
			Summary: "not running: nothing starts or resumes runs on its own; start it with rhiza daemon",
			Next:    "rhiza daemon --autonomous"}
		if s.Stopped != "" {
			p.Detail = append(p.Detail, "last daemon "+s.Process+" stopped "+s.Stopped)
		}
		return healthWatchdogActions(p, s.WatchdogActions)
	}
	p := HealthPart{Part: HealthDaemon, State: HealthOK,
		Summary: fmt.Sprintf("running, process %s, %d live runs", s.Process, len(s.Runs))}
	p = healthRunsAgainstCap(ctx, p, byProject)
	if s.Commit != "" {
		p.Detail = append(p.Detail, "built from "+s.Commit)
		if origin := healthOriginAhead(home, s.Commit); origin != "" {
			p.State = HealthWarn
			p.Summary += fmt.Sprintf("; built from %s, origin/main is at %s", shortSHA(s.Commit), shortSHA(origin))
			p.Next = "make update"
		}
	}
	if s.LastPoll != "" {
		p.Detail = append(p.Detail, "last poll "+s.LastPoll)
	}
	return healthWatchdogActions(healthPollBehind(p, s.PollBehind), s.WatchdogActions)
}

// healthRunsAgainstCap adds the runs against rhiza.daemon.maxRuns to the
// daemon part: live runs (an interrupted one takes no slot) plus starts not
// journaled yet, summed over the daemon's projects, as a detail line. The
// part is bad when they exceed the cap, with a detail line for each live
// run: its project, run and change set title.
func healthRunsAgainstCap(ctx *Ctx, p HealthPart, s DaemonStatusResponse) HealthPart {
	if s.MaxRuns <= 0 {
		return p
	}
	type scoped struct {
		project string
		runs    []DaemonRun
		starts  int
	}
	groups := []scoped{{"", s.Runs, s.Starting}}
	titles := map[string]func(string) string{"": healthTitler(ctx)}
	if len(s.Projects) > 0 {
		scopes, _, err := inboxScopes(ctx, "", true)
		if err != nil {
			scopes = nil
		}
		for _, sc := range scopes {
			titles[sc.name] = healthTitler(sc.ctx)
		}
	}
	for _, pr := range s.Projects {
		groups = append(groups, scoped{pr.Project, pr.Runs, pr.Starting})
	}
	live, starting := 0, 0
	var lines []string
	for _, g := range groups {
		starting += g.starts
		for _, r := range g.runs {
			if r.State == "interrupted" {
				continue
			}
			live++
			name := r.Change
			if t := titles[g.project]; t != nil {
				if title := t(r.Change); title != "" {
					name = title
				}
			}
			where := ""
			if g.project != "" {
				where = g.project + ", "
			}
			lines = append(lines, fmt.Sprintf("live run: %s%s, %s", where, r.ID, name))
		}
	}
	p.Live, p.Starting, p.Max = &live, &starting, &s.MaxRuns
	p.Detail = append(p.Detail, fmt.Sprintf("runs: %d/%d live, %d starting", live, s.MaxRuns, starting))
	if live+starting > s.MaxRuns {
		p.State = HealthBad
		p.Summary += fmt.Sprintf("; over rhiza.daemon.maxRuns %d", s.MaxRuns)
		p.Detail = append(p.Detail, lines...)
	}
	return p
}

// healthTitler gives a change set's display title in ctx's model, empty
// when it cannot be read.
func healthTitler(ctx *Ctx) func(string) string {
	fold, err := daemonFoldOf(ctx)
	return func(cs string) string {
		if err != nil {
			return ""
		}
		c, perr := id.Parse(cs)
		if perr != nil {
			return ""
		}
		if set, ok := fold.changeSet(c); ok {
			return set.Open.DisplayTitle()
		}
		return ""
	}
}

// healthAutonomy lists the autonomy settings and whether the daemon is
// paused, which is warn. It is also warn, never bad, while any workflow is
// on the request-exposed list: outside text reaches its agents with no
// permission checks. Each such workflow gets a detail line, and Next points
// at the workflows guide unless the daemon is paused, which keeps its Next.
// Each workflow's own autonomy declarations are detail lines too.
func healthAutonomy(ctx *Ctx, settings []Setting, err error) HealthPart {
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
	p.Detail = append(p.Detail, healthWorkflowAutonomy(ctx)...)
	if exposed := requestExposedWorkflows(ctx, settings); len(exposed) > 0 {
		p.State = HealthWarn
		p.Summary += "; outside text reaches agents with no permission checks"
		for _, e := range exposed {
			p.Detail = append(p.Detail, ExposureDetail(e))
		}
		p.Next = "rhiza guide workflows"
	}
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
// with each window's use and limit. The use is what the daemon's last poll
// counted and saved, not a fresh read of every transcript, which takes
// seconds; a saved window whose reset has passed is used 0.
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
	used := map[usage.Kind]int64{}
	now := ctx.Now()
	for _, w := range st.UsageWindows {
		if now.Before(w.Reset) {
			used[usage.Kind(w.Kind)] = w.Used
		}
	}
	live := usageLiveLimits(ctx, common)
	var parts []string
	for _, kind := range usage.Kinds {
		used := used[kind]
		if limit, ok := usageLimit(settings, st.Learned, live, kind); ok {
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
	out, err := runRoadmap(ctx, RoadmapRequest{All: true, local: true})
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

// healthWorkflowAutonomy lists, one line each in name order, every workflow
// the repo resolves (personal, repo or built-in) that declares autonomy, with
// each stage it declares and the file it comes from: the owner sees where a
// run's autonomy will come from. A workflow that cannot be read is skipped.
func healthWorkflowAutonomy(ctx *Ctx) []string {
	seen := map[string]bool{}
	for _, d := range []string{personalWorkflowsDir(), filepath.Join(ctx.Root, filepath.FromSlash(workflowsDir))} {
		for _, n := range workflowNames(d) {
			seen[n] = true
		}
	}
	for _, n := range workflow.BuiltinNames() {
		seen[n] = true
	}
	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	var out []string
	for _, n := range names {
		file, def, _, err := resolveWorkflow(ctx.Root, n)
		if err != nil {
			continue
		}
		var stages []string
		for _, st := range workflow.AutonomyStages {
			if m, ok := def.Autonomy(st); ok {
				stages = append(stages, st+" "+m)
			}
		}
		if len(stages) == 0 {
			continue
		}
		if rel, rerr := filepath.Rel(ctx.Root, file); rerr == nil && !strings.HasPrefix(rel, "..") {
			file = filepath.ToSlash(rel)
		}
		out = append(out, fmt.Sprintf("workflow %s (%s): %s", n, file, strings.Join(stages, ", ")))
	}
	return out
}
