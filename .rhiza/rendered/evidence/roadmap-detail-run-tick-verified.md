# roadmap-detail-run-tick-verified

evidence · `01M698Q2HTW1347X7V9YP5VVC6` · actor [claude-code-0c19b72b](../actor/claude-code-0c19b72b.md)

## result

PASS

## statement

At 628be46, go test ./internal/tui/roadmapapp -run Tick -v -count=1 passes TestTickReloadsOnlyTheRunSectionOfAnOpenDetail (opens the detail with a fake run.status in state running, switches it to succeeded, sends TickMsg, sees 'state succeeded' in the view, change.show count unchanged, tick re-arms) and TestTickInListCallsNoRunStatus (TickMsg in the list view makes zero run.status calls). Diff read: Init arms a 3s tea.Tick (Interval = 3*time.Second); TickMsg calls reloadRun, which returns early when m.detail is nil or errored and otherwise replaces only detail.run via latestRun (run.status). gofmt -l is clean, go vet and go test ./... pass; make check's self-check fails only on stale .rhiza/rendered actor/index files for this run's actors, unrelated to the leaf.

## Links

### Out

- verifies roadmap-detail-run-tick (retired)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
