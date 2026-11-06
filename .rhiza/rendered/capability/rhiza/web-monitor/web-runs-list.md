# web-runs-list

capability · `01M6PFKAD5JGK5R4PD07FT82DE` · parent [web-monitor](../web-monitor.md)

## acceptance

A web package test renders /runs over fixture journals with live, done and failed runs and asserts their order, fields and links, and that the list section declares the run events.

## statement

The Runs page at /runs lists live runs first, then recent ones newest first, each named by its change set title, with its workflow, state, current step and start time, linking to its run page; the list re-renders live on run journal moves.

## Links

### Out

- decided-by [web-ui-in-daemon-read-only](../../../decision/web-ui-in-daemon-read-only.md)

### In

- verifies [web-runs-list-verified](../../../evidence/web-runs-list-verified.md)
- verifies [web-runs-list-qa-1](../../../evidence/web-runs-list-qa-1.md)
- verifies [web-runs-list-verified-2](../../../evidence/web-runs-list-verified-2.md)
- verifies [web-runs-list-qa-2](../../../evidence/web-runs-list-qa-2.md)
- verifies [pr-320-web-ui-runs-pages](../../../evidence/pr-320-web-ui-runs-pages.md)

## History

- `01M6PFH1X52KNK0C5ZCDYEEF64` Web Runs pages: runs live first then recent, and a run's detail with its steps, sessions, review summary and a live log tail
