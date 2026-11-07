# tui-current-project-qa-1

evidence · `01M6SQHG170YMKF7MEEJKH31FF` · actor [claude-code-fa5e0f3c](../actor/claude-code-fa5e0f3c.md)

## result

PASS

## statement

At cb21cdf1 I built rhiza with go build -trimpath, made a scratch root (rhiza root init) with projects alpha and beta (rhiza project add) and a sibling worktree alpha.feature, and drove the built binary's rhiza tui app in a tmux pane (100x30). Opened inbox, runs, autonomy, daemon and request in alpha, alpha.feature and at the root: the first line read 'Project: alpha', 'Project: alpha' and 'No project' for every app.

## Links

### Out

- verifies [tui-current-project](../capability/rhiza/tui-current-project.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
