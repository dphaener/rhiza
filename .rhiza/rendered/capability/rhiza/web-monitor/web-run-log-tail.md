# web-run-log-tail

capability · `01M6PFKAF52MNXBPR8STKSW15M` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test renders a fixture run's log tail with its last lines, appends lines to its driver log and reads them through the stream's run log event and the tail fragment, and asserts the no-log state for a run without one.

## statement

A headless run's page shows the last lines of its driver log in an always-dark terminal block that follows new lines live, and says so when the run has no log.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-run-log-tail-verified](../../../evidence/web-run-log-tail-verified.md)
- verifies [web-run-log-tail-qa-1](../../../evidence/web-run-log-tail-qa-1.md)
- verifies [web-run-log-tail-verified-2](../../../evidence/web-run-log-tail-verified-2.md)
- verifies [web-run-log-tail-qa-2](../../../evidence/web-run-log-tail-qa-2.md)
- verifies [pr-320-web-ui-runs-pages](../../../evidence/pr-320-web-ui-runs-pages.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
