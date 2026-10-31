# roadmap-detail-run-tick

capability · `01M686VBX481BWSM5W1YT6N1PQ` · parent [rhiza](../rhiza.md)

## acceptance

A test in internal/tui/roadmapapp opens the detail with a fake run.status saying running, changes it to succeeded, sends the tick and asserts the run section updated without calling change.show again; a tick in the list view calls no run.status

## statement

While the roadmap app's detail view is open it reloads its latest-run section every 3 seconds

## Links

### Out

- decided-by [live-app-name-and-op](../../decision/live-app-name-and-op.md)

### In

- verifies [roadmap-detail-run-tick-verified](../../evidence/roadmap-detail-run-tick-verified.md)
- verifies [roadmap-detail-run-tick-qa-1](../../evidence/roadmap-detail-run-tick-qa-1.md)
- verifies [pr-134-live-runs-app](../../evidence/pr-134-live-runs-app.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
