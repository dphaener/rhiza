# web-overview-live-verified

evidence · `01M6Q8CC4F68RP44P0P9B7N2A1` · actor [claude-code-01131481](../actor/claude-code-01131481.md)

## result

PASS

## statement

At e7511d72, go test -trimpath ./internal/web/... -count=1 passes. TestOverviewSectionsAreLive (internal/web/overview_live_test.go) renders / and parses each section's data-rz-fragment and data-rz-events attributes, checks runs/inbox/milestones/daemon declare run/inbox/change/daemon as event tokens, then GETs each fragment URL and asserts it returns 200 with only its own distinctive marker (incl. 'Paused, 1 of 3 slots in use' for daemon), no <html> or rz-section wrapper, and no other section's marker. Mutation in a scratch copy: dropping the required event from runs.go, inbox.go, milestones.go and daemon.go each makes the test FAIL with '<section> declares events [...], lacks "<event>"', so attempt 1's findings f1 and f2 are fixed.

## Links

### Out

- verifies [web-overview-live](../capability/rhiza/web-monitor/web-overview-live.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
