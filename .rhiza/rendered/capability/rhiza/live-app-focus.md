# live-app-focus

capability · `01M686VAPN2NBZCGV4T968DYTD` · parent [rhiza](../rhiza.md)

## acceptance

A Go test with a fake runner: enter calls run.open with the row's change set, opened true shows focused and tab_gone true shows tab gone; enter on a finished run calls nothing

## statement

Enter on a live run in the live app focuses its Rex tab, and the app says so when the run has no live tab

## Links

### Out

- decided-by [live-app-name-and-op](../../decision/live-app-name-and-op.md)

### In

- verifies [live-app-focus-verified](../../evidence/live-app-focus-verified.md)
- verifies [live-app-focus-qa-1](../../evidence/live-app-focus-qa-1.md)
- verifies [live-app-focus-qa-3](../../evidence/live-app-focus-qa-3.md)
- verifies [pr-134-live-runs-app](../../evidence/pr-134-live-runs-app.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
