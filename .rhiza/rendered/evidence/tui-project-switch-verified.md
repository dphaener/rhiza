# tui-project-switch-verified

evidence · `01M6SPJAHAF4KWEMMQZ9SCHXBF` · actor [claude-code-177617ec](../actor/claude-code-177617ec.md)

## result

PASS

## statement

At a93e818b I ran go test -trimpath -count=1 -run TestSwitch ./internal/tui/modalshell/: TestSwitchProjectReloadsHeaderAndData opens a Switcher-hosted list app in project alpha (header 'Project: alpha', alpha-item rows), presses ctrl+o, sees the picker list alpha and beta, picks beta, and sees header 'Project: beta' with beta-item-1 and no alpha rows; Esc in the picker keeps the project and app open, and no registered projects ignores the key. All PASS, as do go test ./internal/tui/..., go vet, and the tui-* testscripts in cmd/rhiza. In the diff, cmd/rhiza/tui.go routes every non --once app (inbox, runs, request, autonomy, daemon) through modalshell.NewSwitcher with the registered projects from projectctx.Folders, a header from projectctx.Header and a Build that rebuilds the app after the chdir. Note for review: an app rebuilt after a switch gets no fresh WindowSizeMsg (the Switcher does not replay the last one), so e.g. the inbox's later panels skip their size-on-open path.

## Links

### Out

- verifies [tui-project-switch](../capability/rhiza/tui-project-switch.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
