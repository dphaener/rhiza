# web-overview-inbox-verified

evidence · `01M6Q7DBH4HQGEX1M3ZRGGS1ZY` · actor [claude-code-5f74ee0d](../actor/claude-code-5f74ee0d.md)

## result

PASS

## statement

At 44e59aa3 ran go test -trimpath -count=1 ./internal/web/... (ok) and -run TestOverviewListsInbox -v (PASS, plus TestLayoutGolden and TestPartialsGolden PASS); go test -run 'Web|Inbox|Daemon' ./internal/api/... ok. Read the diff: TestOverviewListsInbox serves / and /fragments/inbox via HandlerWithInbox over a fixture inbox fed out of order (failed 09:00, question 07:00, signoff 08:00) and asserts all three appear oldest first (question < sign-off < failed), with the count heading, rz-panel--waiting and rz-pill--waiting warning accent, run links, and no IDs shown; with an empty inbox it asserts the 'Nothing waits on you' empty state and no waiting accent. The daemon wires WebInboxSource, which reads runInbox (the rhiza inbox read) filtered to question, open, signoff, failed, gave-up and stalled.

## Links

### Out

- verifies [web-overview-inbox](../capability/rhiza/web-monitor/web-overview-inbox.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
