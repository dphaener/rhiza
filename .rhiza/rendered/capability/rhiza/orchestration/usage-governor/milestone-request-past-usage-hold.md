# milestone-request-past-usage-hold

capability · `01M6MYSMGQFQDM84ANE2GRY4M9` · parent [usage-governor](../usage-governor.md)

## acceptance

daemon-milestone-request.txtar also covers: with a usage-reserve hold in effect daemon --once starts the queued milestone request and not the queued feature; with rhiza.daemon.paused true, and with the disk-low guard tripped (as the existing disk-low script sets it up), it starts neither; go test -trimpath ./cmd/rhiza -run TestScripts/daemon-milestone-request and go test -trimpath ./internal/api -run Usage pass

## statement

A usage-reserve hold does not hold a milestone request's plan run start, while it still holds feature and bug requests, build, review and other launches; daemon pause and the disk-low guard still hold it

## Links

### Out

- decided-by [milestone-interview-starts-at-once](../../../../decision/milestone-interview-starts-at-once.md)

### In

- verifies [milestone-request-past-usage-hold-verified](../../../../evidence/milestone-request-past-usage-hold-verified.md)
- verifies [milestone-request-past-usage-hold-qa-1](../../../../evidence/milestone-request-past-usage-hold-qa-1.md)
- verifies [pr-291-a-milestone-request](../../../../evidence/pr-291-a-milestone-request.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
