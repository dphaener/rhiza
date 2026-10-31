# live-app-runs-verified-2

evidence · `01M69A15WZATFGTEGKEM21K4G5` · actor [claude-code-fb01944c](../actor/claude-code-fb01944c.md)

## result

PASS

## statement

Attempt 2 at 2d2519d (fix for review findings quality f1, model f1, qa f1). Read the diff: leavesPassed counts each leaf once (dedupes the checks step's fixer entry) and takes the total from change.show's leaves, falling back to distinct listed leaves; View scrolls via top/follow sized from tea.WindowSizeMsg, keeping the daemon line and key help. Ran go test ./internal/tui/... (liveapp passes: TestViewShowsRunProgress asserts title, step, leaves 2/4, check outcomes and elapsed 5m07s on an injected clock; TestTickUpdatesViewKeepingCursorOnRun; new TestLeavesCountChangeSetLeavesOnceEach and TestViewScrollsToCursorKeepingDaemonLine). Ran go test ./cmd/rhiza -run TestScripts/(tui-live-app|describe-golden): pass, tui-live-app cmps rhiza tui app live --once with rhiza run ls --json byte for byte. liveapp.go init registers Name with launcher.Register. make fmt-check vet test pass; make self-check fails only on stale .rhiza/rendered actor files from the review run's chore commits, not this leaf. Real run: a scratch build of rhiza tui app live in an 80x24 Rex pane over this repo showed the daemon line with 2 live runs at the top and the cursor on the first; this change set read leaves 4/4 although run.status lists 5 entries (a review/review fixer entry for live-app-runs), matching change.show's 4 leaves; after 25 j presses the cursor sat on the oldest run, still on screen, with the daemon line and key help kept, 24 lines. Not part of the acceptance and left open: a finished run's elapsed still grows (quality f2, qa f2, minor).

## Links

### Out

- verifies [live-app-runs](../capability/rhiza/live-app-runs.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
