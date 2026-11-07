# tui-starts-in-project-qa-1

evidence · `01M6SQHG561T92WHCP7EDPXYGF` · actor [claude-code-fa5e0f3c](../actor/claude-code-fa5e0f3c.md)

## result

PASS

## statement

At cb21cdf1 I built rhiza with go build -trimpath, made a scratch root (rhiza root init) with projects alpha and beta (rhiza project add) and a sibling worktree alpha.feature, and drove the built binary's rhiza tui app in a tmux pane (100x30). Submitting the Request app opened in beta queued the request in beta only (rhiza request ls per project). Opened at the root, the app showed a Project picker first; Esc exited with nothing queued in either project; picking alpha then submitting queued the request in alpha only.

## Links

### Out

- verifies [tui-starts-in-project](../capability/rhiza/tui-starts-in-project.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
