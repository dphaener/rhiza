# tui-current-project-verified

evidence · `01M6SPR05826WAWHJ5SM48X00C` · actor [claude-code-dcc9e971](../actor/claude-code-dcc9e971.md)

## result

PASS

## statement

At cbe88beb (verifier, attempt 2) I read the diff and ran `go test -trimpath ./cmd/rhiza -run TestTUIAppHeaders -v`. It opens all five apps (inbox, runs, autonomy, daemon, request) through appSwitcher, the path `rhiza tui app` uses, in a registered project's folder (alpha), its sibling worktree (alpha.feature) and the root. All 15 subtests PASS. The header and the first view line are 'Project: alpha', 'Project: alpha' and 'No project', and the view fits a 30-line window. `go test -trimpath ./internal/tui/modalshell/... ./internal/tui/projectctx/...` PASS. Within those, TestHeaderFitsWindow checks that a headed shell, a headed model and the switcher each render exactly the WindowSizeMsg height, because underHeader takes the header row off the content's height. go vet and gofmt are clean. Both attempt-1 findings (header overflowing by one row, no app-level test) are fixed.

## Links

### Out

- verifies [tui-current-project](../capability/rhiza/tui-current-project.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
