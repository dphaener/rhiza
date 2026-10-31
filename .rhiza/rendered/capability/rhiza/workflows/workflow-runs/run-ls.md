# run-ls

capability · `01M686V7WES41E5KRCMQJ9MGCH` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript with seeded run journals and RHIZA_TEST_NOW checks rhiza run ls --json's order, fields and --limit, and its human output; describe-golden.txtar is regenerated

## statement

A run.ls operation lists the repository's runs newest first with each run's ID, change set, title, state, current step, start time and whether it is live, limited by --limit

## Links

### Out

- decided-by [live-app-name-and-op](../../../../decision/live-app-name-and-op.md)

### In

- verifies [run-ls-verified](../../../../evidence/run-ls-verified.md)
- verifies [run-ls-qa-1](../../../../evidence/run-ls-qa-1.md)
- verifies [run-ls-qa-2](../../../../evidence/run-ls-qa-2.md)
- verifies [pr-134-live-runs-app](../../../../evidence/pr-134-live-runs-app.md)

## History

- `01M686RW277XSQ20MYCE580GG8` A live app shows the daemon and each live or recent run's progress, backed by a run.ls operation, and the roadmap detail refreshes its latest run
