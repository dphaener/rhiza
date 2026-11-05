# milestone-request-starts-at-once-verified

evidence · `01M6N360MDXBQCD3P0V73WT42M` · actor [claude-code-96dc1730](../actor/claude-code-96dc1730.md)

## result

PASS

## statement

At c9abeb84 ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-milestone-request|plan-requests)$' -count=1: ok. daemon-milestone-request.txtar (fake rex, --max-runs 1, one ship run already started) queues a feature then a milestone request; daemon --once starts only the milestone request (detach log shows run start --foreground --workflow plan --new First milestone ... --milestone), request ls shows the feature still queued, and a second milestone request queued while the first's run is in flight stays queued on the next --once. Mutations in a scratch copy: an extra --once with only the feature queued starts nothing (the started run really holds the one slot); restoring HEAD~1's internal/api/daemon.go and request.go makes the script fail at line 79 (no plan run started). The diff starts milestone requests first via daemonPlanRequests(onlyKind=milestone) ignoring slots, stops at one, skips while a milestone request's run is live or unjournaled, and counts it in kept/running for later picks. Also ran TestScripts/(daemon|usage|plan-req|request): ok; go test -trimpath ./internal/api/: ok twice, after one earlier FAIL in a run that shared the CPU with the testscripts, not reproduced.

## Links

### Out

- verifies [milestone-request-starts-at-once](../capability/rhiza/orchestration/daemon/plan-requests/milestone-request-starts-at-once.md)

## History

- `01M6MYFH4QYGSKDATBH509SEJF` A milestone request automatically launches an interactive Claude planning session for the owner: when a request of kind milestone is queued (from the CLI or the Request app), the daemon opens a Rex tab running a Claude session in the plan step with the owner as interviewer, not an autonomous plan, and notifies the owner that it is waiting. Builds on change set 01M6MY8BMPFK9ZWD3TDQHKT64J (milestone planning always needs a human). Owner, 2026-11-05: "a milestone request should automatically launch a Claude session to begin the planning"
