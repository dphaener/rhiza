# runs-app-glance-verified-2

evidence · `01M6NBTH4B9MFTDC7Y7WZE8KZC` · actor [claude-code-8cb3cee5](../actor/claude-code-8cb3cee5.md)

## result

PASS

## statement

Verifier claude-code-8cb3cee5 at 4619ed9c (attempt 2, after fixer commit retiring live-app-runs and live-app-focus): go test -trimpath -count=1 ./internal/tui/runsapp/... passed. With a fake runner the tests cover one row per live run with title, current step and elapsed (TestRowsAreLiveRunsOnly), a finished run not listed, a TickMsg reloading the rows and scheduling the next tick (TestTickRefreshes), exactly one 'No live runs' line (TestEmptyState), and every runner call being run.ls across a tick, j and enter (TestOnlyReadOperations). go test -trimpath -count=1 ./cmd/rhiza -run TestScripts/tui-runs passed; tui-runs-app.txtar cmps rhiza tui app runs --once against rhiza run ls --json byte for byte, both empty and with one run. I read runsapp.go: an InfoPanel over run.ls only, reloading every 3s, writing nothing. internal/tui/liveapp is gone with no references left, and both retired nodes show retired=true. make self-check fails only on stale .rhiza/rendered files from the review actors' records. That failure already happened at bfd33acd, and the run's render step regenerates those files.

## Links

### Out

- verifies [runs-app-glance](../capability/rhiza/runs-app-glance.md)

## History

- `01M6MYNJGY6Q735KB2ERT3V4R5` Runs app: a read-only glance at live runs and their current steps, rebuilt on the modal shell over run.ls, replacing internal/tui/liveapp and retiring the capabilities of the old live app it replaces
