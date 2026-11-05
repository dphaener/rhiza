# queued-start-request-slot-kept-verified

evidence · `01M6MS0FF81P9Q13ZKS2N8VCDG` · actor [claude-code-51387005](../actor/claude-code-51387005.md)

## result

PASS

## statement

Verifier at 9d808d15. Read the diff: daemonPoll1 now passes daemonQueuedStarts startMax = max-1 when max > 1 and daemonRequestWaits (a request with no run and retry time passed), else max; the poll order (ship picks, queued continuations, queued starts, plan requests) is unchanged, and daemonPlanRequests still gets the full max. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-queued-start-slot|daemon-request-slot)$' -count=1 -v: both PASS. The new daemon-queued-start-slot.txtar uses a fake rex and no ready change sets, with two plan-workflow starts queued by rhiza run start on draft change sets: with --max-runs 2 and no request, daemon --once starts both and the queue empties; with --max-runs 1, one queued start and a queued request, it starts the queued run, prints no plan run, the request stays 'feature queued' and one start stays queued; with --max-runs 2 and a request, it prints one 'Started a run of' and one 'Started the plan run of request', one start stays queued and the request is 'planning'. Copied the script onto the parent commit (HEAD~1) in a scratch worktree, where it FAILS at line 52 (have 2 'Started a run of', want 1), so it covers the fix. With RHIZA_* unset, go test ./internal/api -run Daemon (12 tests) and the usage_hold tests PASS. The whole internal/api package hit the go test timeout under machine load while running TestInboxReadsChangeSetsOnce, which this change does not touch; PR CI runs the full suite

## Links

### Out

- verifies [queued-start-request-slot-kept](../capability/rhiza/orchestration/daemon/plan-requests/queued-start-request-slot-kept.md)

## History

- `01M6M3VQ85W2NVGR81R06V0XGS` Fix Plan-review runs of drafts bypass the kept slots: the daemon reserves a slot for a waiting plan request or queued continuation (#226, #243), but its plan-review starts (milestone 26's drafts 01M6J375, 01M6J37D on 2026-11-04) took freed slots while the request queue's head waited. Count plan-review starts as ship picks under the same shipMax, or give queued requests priority over them
