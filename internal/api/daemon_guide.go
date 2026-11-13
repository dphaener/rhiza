package api

// This file holds the daemon's guide refresh: each time it starts (after a
// re-exec into a new binary as well as after a restart) it reinstalls the
// installed guides that are not current for the running binary, so no
// session starts with a command naming a skill the guide there lacks.

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/rootcfg"
	"github.com/dphaener/rhiza/internal/workflow"
)

// guideRefresh is one guide scope the daemon looked at when it started.
type guideRefresh struct {
	Scope     string // "user", or "project <name>"
	Dir       string // the plugin folder
	Project   bool   // a project scope
	Refreshed bool   // rewritten for the running binary
	Err       string // why it was not refreshed; empty when it was or needed none
	Tracked   bool   // Err is that the guide is tracked by git in the project
}

// Next is the command that installs the guide for this scope.
func (g guideRefresh) Next() string {
	return guideInstallNext + projectFlag(g.Project)
}

// refreshGuides reinstalls, for the running version, every guide installed
// and not current: the user's, then each registered project's, in name
// order. It never installs where no guide was installed, never touches a
// folder guide install did not write, and never rewrites a project guide
// with a file tracked by git (it reports that as a failed refresh). It
// returns the scopes it refreshed or failed to, and logs each.
func refreshGuides(ctx *Ctx, root string, reg rootcfg.Config, log func(string)) []guideRefresh {
	var out []guideRefresh
	check := func(scope, base string, project bool) {
		dir := filepath.Join(base, ".claude", "skills", pluginName)
		v, ok := installedVersion(dir)
		if !ok || (v == ctx.Version && pluginCurrent(dir, ctx.Version)) {
			return
		}
		g := guideRefresh{Scope: scope, Dir: dir, Project: project}
		defer func() {
			out = append(out, g)
			switch {
			case g.Refreshed:
				log(fmt.Sprintf("refreshed the %s guide in %s to %s", scope, dir, ctx.Version))
			default:
				log(fmt.Sprintf("could not refresh the %s guide in %s: %s; run %s", scope, dir, g.Err, g.Next()))
			}
		}()
		pctx := ctx
		if project {
			pctx = projectCtx(ctx, base)
			tracked, err := gitAt(pctx, base, "ls-files", "--", ".claude/skills/"+pluginName)
			if err != nil {
				g.Err = "git ls-files: " + err.Error()
				return
			}
			if strings.TrimSpace(tracked) != "" {
				g.Tracked = true
				g.Err = "the guide is committed in the repo: refresh it with " + g.Next() + " and commit it through a change set"
				return
			}
		}
		if _, err := runGuideInstall(pctx, GuideInstallRequest{Agent: "claude", Project: project}); err != nil {
			g.Err = err.Error()
			return
		}
		g.Refreshed = true
	}
	if ctx.Home != "" {
		check("user", ctx.Home, false)
	}
	names := slices.Clone(reg.Names())
	slices.Sort(names)
	for _, name := range names {
		dir := reg.Dir(root, name)
		if !projectInitialised(dir) {
			continue
		}
		check("project "+name, dir, true)
	}
	return out
}

// guideGate holds the daemon's starts while a guide it could not refresh
// still lacks a skill a step's start command names. It keeps the scopes the
// start-up refresh failed or skipped; each poll it drops those whose guide
// has since become current, so the owner's rhiza guide install clears the
// hold on the next poll.
type guideGate struct {
	version string
	home    string // the user's home, "" when the daemon has none
	since   time.Time
	failed  []guideRefresh
}

// newGuideGate is the gate for the scopes refreshGuides returned: nil when
// none of them failed.
func newGuideGate(ctx *Ctx, scopes []guideRefresh) *guideGate {
	g := &guideGate{version: ctx.Version, home: ctx.Home, since: ctx.Now().UTC()}
	for _, sc := range scopes {
		if !sc.Refreshed {
			g.failed = append(g.failed, sc)
		}
	}
	if len(g.failed) == 0 {
		return nil
	}
	return g
}

// holds is the failed scopes that still hold for the project at dir: the
// user's, and the project's own, whose guide is still not current.
func (g *guideGate) holds(dir string) []guideRefresh {
	if g == nil {
		return nil
	}
	var out []guideRefresh
	for _, sc := range g.failed {
		if sc.Project && sc.Dir != filepath.Join(dir, ".claude", "skills", pluginName) {
			continue
		}
		if pluginCurrent(sc.Dir, g.version) {
			continue
		}
		out = append(out, sc)
	}
	return out
}

// state is the holds as the daemon's state records them.
func (g *guideGate) state(dir string) []daemon.GuideHold {
	if g == nil {
		return nil
	}
	var out []daemon.GuideHold
	for _, sc := range g.holds(dir) {
		out = append(out, daemon.GuideHold{Since: g.since, Scope: sc.Scope, Dir: sc.Dir, Error: sc.Err, Next: sc.Next()})
	}
	return out
}

// skillPresent says the guide folder dir holds the skill name.
func skillPresent(dir, name string) bool {
	p := filepath.Join(dir, "skills", name, "SKILL.md")
	if name == pluginName {
		p = filepath.Join(dir, pluginSkillPath)
	}
	_, err := os.Stat(p)
	return err == nil
}

// missingSkill is a skill of the guide that a start command of workflow
// names and no installed guide for the project at dir holds, while a hold
// stands; "" when none is missing or no hold stands.
func (g *guideGate) missingSkill(dir string, def *workflow.Definition) string {
	if len(g.holds(dir)) == 0 || def == nil {
		return ""
	}
	files, err := pluginFiles(g.version)
	if err != nil {
		return ""
	}
	var dirs []string
	if g.home != "" {
		dirs = append(dirs, filepath.Join(g.home, ".claude", "skills", pluginName))
	}
	dirs = append(dirs, filepath.Join(dir, ".claude", "skills", pluginName))
	for _, name := range workflowStartSkills(def) {
		path := "skills/" + name + "/SKILL.md"
		if name == pluginName {
			path = pluginSkillPath
		}
		if _, ours := files[path]; !ours {
			continue
		}
		if !slices.ContainsFunc(dirs, func(d string) bool { return skillPresent(d, name) }) {
			return name
		}
	}
	return ""
}

// workflowStartSkills are the skill names, in order of first use, that the
// start commands of def's steps and roles name, each read by the adapter of
// its agent (the default adapter for none); an agent whose adapter reads no
// slash commands names none.
func workflowStartSkills(def *workflow.Definition) []string {
	var out []string
	add := func(agent, prompt string) {
		a := adapter.Default()
		if agent != "" {
			a, _ = adapter.Lookup(agent)
		}
		sc, ok := a.(adapter.SkillChecker)
		if !ok {
			return
		}
		name, ok := sc.StartSkill(strings.TrimSpace(prompt))
		if !ok {
			return
		}
		name = strings.TrimPrefix(name, pluginName+":")
		if name != "" && !slices.Contains(out, name) {
			out = append(out, name)
		}
	}
	for _, n := range def.Order {
		st := def.Steps[n]
		add(st.Agent, st.Prompt)
		for _, r := range []*workflow.Role{st.Builder, st.Verifier, st.Critic, st.Fixer} {
			if r != nil {
				add(r.Agent, r.Prompt)
			}
		}
		for _, c := range st.Checks {
			add(c.Agent, c.Prompt)
		}
	}
	return out
}

// guideSkipOf is the skip for change set cs when its workflow names a skill
// the held guide lacks: the skip, and true.
func guideSkipOf(ctx *Ctx, fold *daemonFolded, cs string) (daemon.Skip, bool) {
	g := ctx.guideGate
	if g == nil || len(g.holds(ctx.Root)) == 0 {
		return daemon.Skip{}, false
	}
	cid, err := id.Parse(cs)
	if err != nil {
		return daemon.Skip{}, false
	}
	c, found := fold.changeSet(cid)
	if !found {
		return daemon.Skip{}, false
	}
	name, err := readyWorkflow(ctx.Root, c.Open)
	if err != nil {
		return daemon.Skip{}, false
	}
	_, def, _, err := resolveWorkflow(ctx.Root, name)
	if err != nil {
		return daemon.Skip{}, false
	}
	skill := g.missingSkill(ctx.Root, def)
	if skill == "" {
		return daemon.Skip{}, false
	}
	hs := g.holds(ctx.Root)
	return daemon.Skip{Change: cs, Reason: daemon.SkipGuideStale, Detail: []string{
		"workflow " + name + " starts /" + skill + ", missing from the installed guide", hs[0].Err, hs[0].Next()}}, true
}

// guideHoldItems is one guide-stale item per hold.
func guideHoldItems(hs []daemon.GuideHold) []InboxItem {
	var out []InboxItem
	for _, h := range hs {
		at := h.Since.UTC().Format(time.RFC3339)
		sum := sha256.Sum256([]byte(strings.Join([]string{InboxGuideStale, h.Scope, at}, "\x00")))
		out = append(out, InboxItem{
			ID: hex.EncodeToString(sum[:])[:7], Kind: InboxGuideStale, At: at, Next: h.Next,
			Detail: fmt.Sprintf("the %s guide in %s is out of date and was not refreshed: %s; the daemon starts no step whose start command names a skill it lacks until it is current; run %s", h.Scope, h.Dir, h.Error, h.Next),
		})
	}
	return out
}

// inboxGuideHold is the guide-stale items of the daemon's last poll.
func inboxGuideHold(ctx *Ctx) []InboxItem {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil
	}
	st, ok := daemon.ReadState(common)
	if !ok {
		return nil
	}
	return guideHoldItems(st.GuideHolds)
}
