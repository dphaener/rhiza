# inbox-milestone-interview-verified

evidence · `01M6N2B64WDP5AS1CXE15QN1NP` · actor [claude-code-c479b398](../actor/claude-code-c479b398.md)

## result

PASS

## statement

At 31e81310, go test -trimpath ./internal/api -run 'Inbox|DaemonNotify' passes (also -count=1 -v: TestInboxMilestoneInterview and TestDaemonNotifyMilestoneInterview PASS). The tests show a live non-autonomous milestone plan run gives exactly one interview item with the interview session's tab (title, block, session) from the interview session's launch, with no question block; none once the run ended; none for a feature plan run or an autonomous one; and daemonNotify run twice sends one blocked Rex report 'Plan it: a milestone interview waits on you'. describe-golden.txtar diff updates the InboxItem kind list (adds interview) and tab/session descs; go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes. Diff read: inboxItems adds InboxInterview inside the live branch keyed on st.Session[interview], itemTab opens its tab, inboxNotice has the text; gofmt and go vet clean.

## Links

### Out

- verifies [inbox-milestone-interview](../capability/rhiza/orchestration/inbox/inbox-milestone-interview.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
