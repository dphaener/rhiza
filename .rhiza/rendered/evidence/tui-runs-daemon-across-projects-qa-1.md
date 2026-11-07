# tui-runs-daemon-across-projects-qa-1

evidence · `01M6SQHG95TF25AADVBMAEK527` · actor [claude-code-fa5e0f3c](../actor/claude-code-fa5e0f3c.md)

## result

PASS

## statement

At cb21cdf1 I built rhiza with go build -trimpath, made a scratch root (rhiza root init) with projects alpha and beta (rhiza project add) and a sibling worktree alpha.feature, and drove the built binary's rhiza tui app in a tmux pane (100x30). With a live foreground run (sleep step) in each project, the Runs app in alpha and at the root listed 'alpha · Slow run in alpha' and 'beta · Slow run in beta'. The Daemon app in beta showed Pause the daemon and Run limit once, then Pause project beta and Project run limit; saving Run limit 5 and Project run limit 3 wrote rhiza.daemon.maxRuns=5 to the root's rhiza.toml and rhiza.project.maxruns=3 to beta's git config only. Setting Plan to auto in the Autonomy app (switched to beta) wrote rhiza.autonomy.plan=auto to beta's git config and nothing to alpha's.

## Links

### Out

- verifies [tui-runs-daemon-across-projects](../capability/rhiza/tui-runs-daemon-across-projects.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
