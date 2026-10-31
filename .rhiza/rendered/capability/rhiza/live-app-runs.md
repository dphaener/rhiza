# live-app-runs

capability · `01M686V969FDV2SVYM7P244E86` · parent [rhiza](../rhiza.md)

## acceptance

A Go test in internal/tui/liveapp over a fake runner with canned daemon.status, run.ls and run.status results and an injected clock asserts the view's title, step, leaves count, check outcomes and elapsed time, and that a tick with changed results updates the view keeping the cursor on the same run; a testscript checks rhiza tui app live --once prints run.ls's JSON, and the app is in the launcher registry

## statement

A live app, registered as a Rex app, shows the daemon's state and each live or recent run with its change set title, current step, leaves passed of total, checks and elapsed time, refreshing every 3 seconds

## Links

### Out

- decided-by [live-app-name-and-op](../../decision/live-app-name-and-op.md)

### In

- verifies [live-app-runs-verified](../../evidence/live-app-runs-verified.md)
- verifies [live-app-runs-qa-1](../../evidence/live-app-runs-qa-1.md)
- verifies [live-app-runs-verified-2](../../evidence/live-app-runs-verified-2.md)
- verifies [live-app-runs-qa-2](../../evidence/live-app-runs-qa-2.md)
- verifies [live-app-runs-verified-5](../../evidence/live-app-runs-verified-5.md)
- verifies [live-app-runs-qa-3](../../evidence/live-app-runs-qa-3.md)
- verifies [pr-134-live-runs-app](../../evidence/pr-134-live-runs-app.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
