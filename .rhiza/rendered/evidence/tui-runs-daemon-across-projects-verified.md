# tui-runs-daemon-across-projects-verified

evidence · `01M6SQ1FATQX53NTQKAT9PASFQ` · actor [claude-code-a13da17e](../actor/claude-code-a13da17e.md)

## result

PASS

## statement

At 196c50e2 ran go test -trimpath ./internal/tui/runsapp/... ./internal/tui/daemonapp/... ./internal/tui/autonomyapp/... (ok) and go test -trimpath ./cmd/rhiza -run TestTUI (TestTUIAppHeaders, TestTUISettingsScopes PASS); gofmt and go vet clean. TestMergedTagsRunsByProject runs the Runs app over projects alpha and beta and sees rows tagged 'alpha · Alpha work' and 'beta · Beta work', newest first, finished run dropped. TestTUISettingsScopes registers alpha and beta under a root, opens the Autonomy app in beta and toggles ship to owner: beta's git config gets rhiza.autonomy.ship=owner and alpha's stays unset; the Daemon app's run limit 7 lands in the root's rhiza.toml and not in git config. Read the diff: tui.go wires the merged runs model when projects are registered and daemonapp.NewIn with the current project; daemonapp adds project pause and run-cap fields writing rhiza.project.* keys. Checked by hand with a scratch build in a scratch root (two projects, RHIZA_* unset): daemon set of rhiza.daemon.maxRuns 7 wrote [daemon] in rhiza.toml, rhiza.project.maxRuns and rhiza.autonomy.ship wrote only beta's git config, and daemon settings lists each key once with its scope.

## Links

### Out

- verifies [tui-runs-daemon-across-projects](../capability/rhiza/tui-runs-daemon-across-projects.md)

## History

- `01M6RN441P2Z865ZVKAS7X7GMG` Project-aware TUI apps: the Rex apps know the current project, start milestones, change sets, plans and requests in it, switch project, and show a merged inbox
