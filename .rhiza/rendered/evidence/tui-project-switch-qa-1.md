# tui-project-switch-qa-1

evidence · `01M6SQHG33PQA70VCK8NBESV3F` · actor [claude-code-fa5e0f3c](../actor/claude-code-fa5e0f3c.md)

## result

PASS

## statement

At cb21cdf1 I built rhiza with go build -trimpath, made a scratch root (rhiza root init) with projects alpha and beta (rhiza project add) and a sibling worktree alpha.feature, and drove the built binary's rhiza tui app in a tmux pane (100x30). In the Autonomy app opened in alpha (Plan: owner), ctrl+o showed a picker of alpha and beta; picking beta changed the header to 'Project: beta' and the form to beta's settings (Plan: auto, set earlier from the app). The Runs and Request apps switched the same way, and a request submitted after switching to beta from alpha was queued in beta (rhiza request ls in beta), not in alpha.

## Links

### Out

- verifies [tui-project-switch](../capability/rhiza/tui-project-switch.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
