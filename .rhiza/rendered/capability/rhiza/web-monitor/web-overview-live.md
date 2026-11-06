# web-overview-live

capability · `01M6PFJX4SRG58AJJWHZYKN4DK` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test asserts each Overview section declares its event names and that each section's fragment URL renders that section alone.

## statement

Each Overview section is a live section that re-renders on its events: the runs section on run journal moves, the inbox on inbox changes, the milestone section on change set moves, the daemon section on the stream's daemon status event.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-overview-live-verified](../../../evidence/web-overview-live-verified.md)
- verifies [web-overview-live-qa-1](../../../evidence/web-overview-live-qa-1.md)
- verifies [web-overview-live-qa-2](../../../evidence/web-overview-live-qa-2.md)
- verifies [pr-323-web-ui-overview-page](../../../evidence/pr-323-web-ui-overview-page.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
