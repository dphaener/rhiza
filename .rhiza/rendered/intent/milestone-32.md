# milestone-32

intent · `01M7640BA83BXBMSEMZNC4K97G`

## milestone

32

## statement

Runs clean up after themselves and Rhiza's shipped defaults are language-neutral: no Go-, make- or Rhiza-specific behaviour ships to users, Rhiza keeps its own in its own workflows, and auto-update is a documented setting with channels.

## History

- `01M75W99SF6Q1T4QJM5GKA4PM1` Workflow cleanup and language-neutral defaults. Owner, 2026-11-11: "I'd like to propose a feature for rhiza. I think the Go cache is causing a lot of disk bloat, so the workflows should manage this when they are done. THere should be a cleanup phase. This is where the default Rhiza workflow (what is shipped and installed for all users as the default) from what we use to build Rhiza itself. This may actually be a milestone because I think it's going to cover a lot of ground. So we want to add a cleanup phase for THIS Rhiza only. But we may want to add some guidance in the docs, some workflow examples for different languages, etc. THis will require some milestone planning for sure. I also noticed when you were talking about the existing daemon features that the daemon was specifically looking for things like make check, which won't be universal to all programming languages, so things like that need to be workflow dependent, project dependent, or a global setting. Let's get to it!" Orchestrator's findings to start from (verify, don't trust): on the VM disk is 80% used (73G of 97G); ~/.cache/go-build 1.3G, ~/go/pkg/mod 3.1G, run worktrees 55-210M each. Go- and Rhiza-specific behaviour shipped in the default product today: the daemon setting rhiza.daemon.trimGoCache (internal/api/daemon.go:218, internal/gosweep) trims the Go build cache only when disk runs low; internal/tmpsweep sweeps go test temp folders; the watchdog's transient signatures include Go's go-build $WORK error (internal/api/daemon_watchdog.go:425); ship's head build runs go build ./cmd/rhiza only in Rhiza's own repo (internal/api/ship_build.go:85); daemon self-update runs make (internal/api/daemon.go:336); rhiza workflow new scaffolds run = "make test" (internal/api/workflow_new.go:75); ci-triage parses Go, RSpec, Minitest and pytest failures (internal/api/ship.go:134,1049); the local gate is already a project setting (rhiza.ship.localGate). The built-in workflows are internal/workflow/builtin/{plan,plan-milestone,build,review,ship-change}.toml; Rhiza's repo adds .rhiza/workflows/docs*.toml and the bugfix workflow (#435). Scope to settle with the owner: a cleanup step for Rhiza's own workflows (Go cache, worktrees, temp folders) versus the shipped defaults; which language- or tool-specific behaviours move into workflows, project settings or global settings; docs guidance and per-language workflow examples.
