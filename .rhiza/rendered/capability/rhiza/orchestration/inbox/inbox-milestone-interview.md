# inbox-milestone-interview

capability · `01M6MYT17VS6C6MX8YRDZ6PMKD` · parent [inbox](../inbox.md)

## acceptance

go test -trimpath ./internal/api -run 'Inbox|DaemonNotify' passes with tests that a live non-autonomous milestone plan run gives one interview item with its tab from launch, before any question block, that it is gone once the run ended, that a feature plan run gives none, and that daemonNotify notifies it once with that text; describe-golden.txtar is regenerated for the InboxItem kind list and TestScripts/describe-golden passes

## statement

rhiza inbox lists a live milestone plan run as an item of a new kind, interview, naming its change set and the interview's Rex tab, from the interview's launch until the run ends, and the daemon sends one Rex notification for it (<title>: a milestone interview waits on you)

## Links

### Out

- decided-by [milestone-interview-starts-at-once](../../../../decision/milestone-interview-starts-at-once.md)

### In

- verifies [inbox-milestone-interview-verified](../../../../evidence/inbox-milestone-interview-verified.md)
- verifies [inbox-milestone-interview-qa-1](../../../../evidence/inbox-milestone-interview-qa-1.md)
- verifies [pr-291-a-milestone-request](../../../../evidence/pr-291-a-milestone-request.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
