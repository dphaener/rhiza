# runs-app-glance-verified

evidence · `01M6NAE91BFPRTJ0MVVC171E9N` · actor [claude-code-a5dbd61d](../actor/claude-code-a5dbd61d.md)

## result

PASS

## statement

At ff652bac, go test -trimpath -count=1 -v ./internal/tui/runsapp/ passed 7 tests with a fake runner: TestRowsAreLiveRunsOnly (live run shows title, current step review and elapsed 1m05s; finished run not listed), TestTickRefreshes (TickMsg reloads rows and schedules the next tick), TestEmptyState (exactly one 'No live runs' line), TestOnlyReadOperations (after tick and keypresses every runner op is run.ls), TestResultIsRunLs, TestClosesInShell, TestRegistered. go test ./cmd/rhiza -run 'TestScripts/(tui-runs-app|rex-install)$' passed: tui-runs-app cmps rhiza tui app runs --once against rhiza run ls --json byte for byte, empty and with a run. healthapp tests pass; go build ./... and go vet clean. Diff read: runsapp is a rexuikit InfoPanel in the modal shell calling only run.ls; no liveapp references remain.

## Links

### Out

- verifies [runs-app-glance](../capability/rhiza/runs-app-glance.md)

## History

- `01M6MYNJGY6Q735KB2ERT3V4R5` Runs app: a read-only glance at live runs and their current steps, rebuilt on the modal shell over run.ls, replacing internal/tui/liveapp and retiring the capabilities of the old live app it replaces
