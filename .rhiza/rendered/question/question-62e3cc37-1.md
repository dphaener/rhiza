# question-62e3cc37-1

question · `01M6M3ZS1WHXPTCRBWMWD168EH`

## session

62e3cc37-a9a9-440e-a34d-01a1a83df551

## statement

The plan-review runs on 2026-11-04 were owner-queued starts (rhiza run start --workflow plan-review.toml, started by daemonQueuedStarts), which run after ship picks but before daemonPlanRequests and ignore the slot kept for a waiting plan request. How should the daemon fix it? / This change set serves no intent yet. Recent request-planned fixes (e.g. 01M6M1M1 'Fix Test temp folders leak') serve milestone-25 ('The backlog is empty…'). Which intent should it serve?

## Links

### In

- answers [queued-starts-keep-request-slot](../decision/queued-starts-keep-request-slot.md)

## History

- `01M6M3VQ85W2NVGR81R06V0XGS` Fix Plan-review runs of drafts bypass the kept slots: the daemon reserves a slot for a waiting plan request or queued continuation (#226, #243), but its plan-review starts (milestone 26's drafts 01M6J375, 01M6J37D on 2026-11-04) took freed slots while the request queue's head waited. Count plan-review starts as ship picks under the same shipMax, or give queued requests priority over them
