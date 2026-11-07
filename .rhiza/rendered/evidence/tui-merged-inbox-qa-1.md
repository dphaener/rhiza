# tui-merged-inbox-qa-1

evidence · `01M6SQHG6M8ZARVKKYBVG8CBQE` · actor [claude-code-fa5e0f3c](../actor/claude-code-fa5e0f3c.md)

## result

PASS

## statement

At cb21cdf1 I built rhiza with go build -trimpath, made a scratch root (rhiza root init) with projects alpha and beta (rhiza project add) and a sibling worktree alpha.feature, and drove the built binary's rhiza tui app in a tmux pane (100x30). With a blocker question landed in beta and none in alpha, the Inbox app opened in alpha listed 'beta · blocker · Blocked  Which disk holds beta's store?'. Answering it from alpha wrote Decision answer-store-q into beta, beta's inbox emptied, and alpha got no change set.

## Links

### Out

- verifies [tui-merged-inbox](../capability/rhiza/tui-merged-inbox.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
