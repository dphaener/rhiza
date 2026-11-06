# web-overview-daemon

capability · `01M6PFJX06P8P1JQS8RKY327J2` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test renders / over a fixture repo with daemon status and health data and asserts the paused state, the slots and each health part with its state.

## statement

The Overview page at / shows the daemon: running or paused, run slots used of its limit, and each health part with its ok, warn or bad state.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-overview-daemon-verified](../../../evidence/web-overview-daemon-verified.md)
- verifies [web-overview-daemon-qa-1](../../../evidence/web-overview-daemon-qa-1.md)
- verifies [web-overview-daemon-qa-2](../../../evidence/web-overview-daemon-qa-2.md)
- verifies [pr-323-web-ui-overview-page](../../../evidence/pr-323-web-ui-overview-page.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
