package api

// This file holds the daemon's guide refresh: each time it starts (after a
// re-exec into a new binary as well as after a restart) it reinstalls the
// installed guides that are not current for the running binary, so no
// session starts with a command naming a skill the guide there lacks.

import (
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/rootcfg"
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
