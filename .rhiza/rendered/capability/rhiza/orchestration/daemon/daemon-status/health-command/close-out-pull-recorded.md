# close-out-pull-recorded

capability · `01M6CS2CVV2KZ1JBQ89ANTZW1J` · parent [health-command](../health-command.md)

## acceptance

A close-out test with a dirty main checkout records the pull failure; health shows the main part warn naming it and the inbox lists a main-stale item; a later successful pull clears both

## statement

When ship close-out cannot fast-forward the main checkout, it records the failure where rhiza health and the inbox show it, instead of only in its step detail

## Links

### Out

- decided-by [health-one-screen](../../../../../../decision/health-one-screen.md)

### In

- blocks [close-out-pull-needs-health-report](../../../../../../question/close-out-pull-needs-health-report.md)
- verifies [close-out-pull-recorded-verified](../../../../../../evidence/close-out-pull-recorded-verified.md)
- verifies [close-out-pull-recorded-qa-1](../../../../../../evidence/close-out-pull-recorded-qa-1.md)
- verifies [close-out-pull-recorded-verified-2](../../../../../../evidence/close-out-pull-recorded-verified-2.md)
- verifies [pr-194-health-app](../../../../../../evidence/pr-194-health-app.md)

## History

- `01M6CRZSRFQ94MAJ8CRTEE44P2` rhiza health and a health app show the whole system's state on one screen
