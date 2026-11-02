# health-report

capability · `01M6D4RPTDYGEG48RYPD0557BA` · parent [health-command](../health-command.md)

## acceptance

health.txtar: with no daemon the daemon part is warn and the report says how to start it; with RHIZA_TEST_FREE_BYTES under the minimum the disk part is bad and health exits 1; a dirty main checkout is warn; an inbox failed item is warn; a model problem makes check bad; health --json and api call health print the same bytes; describe-golden.txtar is regenerated

## statement

rhiza health reports the daemon, autonomy settings, free disk, the main checkout (branch, clean, behind origin), inbox counts by kind, queued requests, roadmap counts and rhiza check's result, each part ok, warn or bad, and exits 1 when any part is bad

## Links

### Out

- decided-by [health-one-screen](../../../../../../decision/health-one-screen.md)

### In

- verifies [health-report-verified](../../../../../../evidence/health-report-verified.md)
- verifies [health-report-qa-1](../../../../../../evidence/health-report-qa-1.md)
- verifies [pr-194-health-app](../../../../../../evidence/pr-194-health-app.md)

## History

- `01M6CRZSRFQ94MAJ8CRTEE44P2` rhiza health and a health app show the whole system's state on one screen
