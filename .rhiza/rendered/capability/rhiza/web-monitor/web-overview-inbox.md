# web-overview-inbox

capability · `01M6PFJX24XAF6CBGRW7W837AV` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test renders / over a fixture inbox with a question, a sign-off and a failed run and asserts all three, oldest first, and an empty state when the inbox is empty.

## statement

The Overview page lists what waits on the owner, as rhiza inbox does (questions, sign-offs, failed or given-up runs, stalled sessions), oldest first, marked in the warning accent.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-overview-inbox-verified](../../../evidence/web-overview-inbox-verified.md)
- verifies [web-overview-inbox-qa-1](../../../evidence/web-overview-inbox-qa-1.md)
- verifies [web-overview-inbox-verified-2](../../../evidence/web-overview-inbox-verified-2.md)
- verifies [web-overview-inbox-qa-2](../../../evidence/web-overview-inbox-qa-2.md)
- verifies [pr-323-web-ui-overview-page](../../../evidence/pr-323-web-ui-overview-page.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
