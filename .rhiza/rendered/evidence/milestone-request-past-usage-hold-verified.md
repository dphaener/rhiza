# milestone-request-past-usage-hold-verified

evidence · `01M6N1CMY65VR6XSFSBPERZ222` · actor [claude-code-8ff91e3c](../actor/claude-code-8ff91e3c.md)

## result

PASS

## statement

At 01f5053a: go test -trimpath ./cmd/rhiza -run TestScripts/daemon-milestone-request -count=1 -v passes; the new script queues a feature then a milestone request with a Claude transcript filling the 1000-token session limit: with rhiza.daemon.paused true and with RHIZA_TEST_FREE_BYTES at 1GiB (disk-low) daemon --once starts nothing (no detach.log), and under the usage hold it starts only the milestone request's plan run while 'Add exports' stays queued (oldest-first order makes this discriminating). go test -trimpath ./internal/api -run Usage passes (incl. TestUsageHoldStopsLaunches and TestUsageHoldStopsEveryLaunch, which still holds a feature request), as do TestScripts/usage-hold and TestDaemonPlanRequests*. Diff read: daemonPoll1 now returns early only on paused||low; under a hold it skips resumes, ship picks, queued continues and queued starts, and calls daemonPlanRequests with onlyKind=milestone, which skips other kinds.

## Links

### Out

- verifies [milestone-request-past-usage-hold](../capability/rhiza/orchestration/usage-governor/milestone-request-past-usage-hold.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
