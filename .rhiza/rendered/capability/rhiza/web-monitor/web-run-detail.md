# web-run-detail

capability · `01M6PFKAE3P8MJB7E7CJ4FCWHB` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test renders /runs/<id> for a fixture plan run and a fixture review run and asserts the steps in order with outcomes, the sessions and the review summary's findings, and a 404 for an unknown run.

## statement

A run's page at /runs/<id> shows its steps from the journal as a timeline with each step's outcome and time, the agent sessions of each step with their status, and for a review run its review summary with findings and severities; it re-renders live on that run's journal moves.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-run-detail-verified](../../../evidence/web-run-detail-verified.md)
- verifies [web-run-detail-qa-1](../../../evidence/web-run-detail-qa-1.md)
- verifies [web-run-detail-verified-2](../../../evidence/web-run-detail-verified-2.md)
- verifies [web-run-detail-qa-2](../../../evidence/web-run-detail-qa-2.md)
- verifies [pr-320-web-ui-runs-pages](../../../evidence/pr-320-web-ui-runs-pages.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
