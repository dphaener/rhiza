# web-overview-inbox-verified-2

evidence · `01M6QB6VNNFCAKZCMKSCPN0KYA` · actor [claude-code-53dff1fd](../actor/claude-code-53dff1fd.md)

## result

PASS

## statement

Verifier, attempt 2 at 2c8d7ef9: go test -trimpath -count=1 ./internal/web -run TestOverview passes. TestOverviewListsInbox renders / and /fragments/inbox over a fixture inbox with a question, a sign-off and a failed run, asserts all three oldest first (question 07:00, sign-off 08:00, failed 09:00) in the waiting accent with no ID shown, and asserts the empty state 'Nothing waits on you' without the waiting accent. The fix commit's TestWebInboxViewsOwnerWording and TestWebInboxViewsOpenWithoutQuestions pass (go test ./internal/api -run TestWebInbox). They fold run journals into open, sign-off, failed and stalled inbox items and assert that no row's title or detail carries a run or session ID, a question handle, a rhiza command or a Rex window, and that an open item becomes one row per question, by its statement. That resolves the quality finding about passing the CLI detail straight through. Read the diff: webInboxViews maps each kind to owner wording, and Blocks is the blocked node's name (inbox.go:125), not a handle. go vet and gofmt are clean, and internal/web passes in full. The full internal/api run failed only TestReadBudget (roadmap what-if median over 50ms at a load average of about 17), a timing budget on roadmap reads that this diff does not touch.

## Links

### Out

- verifies [web-overview-inbox](../capability/rhiza/web-monitor/web-overview-inbox.md)

## History

- `01M6PFGXV205YBYX2MZ9CC590B` Web Overview page: the daemon and its health, live runs and their current step, the inbox and the active milestone's progress on one live page
