# live-app-runs-verified

evidence · `01M698HT9437SVEDTF5SS572TJ` · actor [claude-code-fc402633](../actor/claude-code-fc402633.md)

## result

PASS

## statement

Verifier at e5429a5: go test -count=1 -v ./internal/tui/liveapp passes TestViewShowsRunProgress (fake runner with canned daemon.status, run.ls and run.status and an injected clock; asserts title, step, leaves 2/3, checks 1 passed, 1 failed, 1 running, elapsed 5m07s/1h05m07s and the daemon line), TestTickUpdatesViewKeepingCursorOnRun (a TickMsg with changed results updates the daemon line, a new run and leaves 3/3 while the cursor stays on R1), TestEmptyAndRegistered (live is in launcher.Apps()) and TestResultIsRunLsBytes. go test ./cmd/rhiza -run TestScripts/tui-live-app passes: rhiza tui app live --once cmp-matches rhiza run ls --json with no runs and with one run. Reading the diff: Interval is 3s via tea.Tick; launcher.Register in init; the app is wired in cmd/rhiza/tui.go. A scratch build's rhiza tui app live --once printed this repo's live run.ls JSON. make check exits 0 (one earlier run failed only on read-budget timing under load and passed on rerun). Note: the leaves total is len(run.status leaves), meaning the leaves conducted so far, not every leaf in the change set.

## Links

### Out

- verifies [live-app-runs](../capability/rhiza/live-app-runs.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
