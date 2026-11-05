# plan-review-runs-slot-interview

source · `01M6M42RFTNGGRFNT2CF6Y1DCC`

## quotes

```json
[]
```

## raw_sha256

b9c68f956a3e96593b2d7934892637187b51d18294d767ee7f2dc65778a76ac2

## session

62e3cc37-a9a9-440e-a34d-01a1a83df551

## summary

Autonomous plan run, no owner. The planner traced the bug: the plan-review runs of milestone 26's drafts (runs 01M6KYKJ on 01M6J375 and 01M6KZ9E on 01M6J37D, 2026-11-04) were owner-queued starts (rhiza run start with a plan-review workflow file), which daemonQueuedStarts starts up to max-runs after ship picks and continuations but before daemonPlanRequests; only ship picks honoured the slot kept for a waiting plan request (#226), so request 01M6KTAZ (queued 17:40) waited until 19:19. Two questions went to the arbiter: the fix (queued starts keep the slot, or plan requests before queued starts) and the intent. The arbiter decided queued-starts-keep-request-slot: while a queued plan request could start and max > 1, queued starts start only while running < max-1; poll order unchanged; the change set serves milestone-25. The planner wrote one leaf under plan-requests.

## title

Plan of Fix Plan-review runs of drafts (autonomous)

## Links

### In

- derived-from [queued-starts-keep-request-slot](../decision/queued-starts-keep-request-slot.md)

## History

- `01M6M3VQ85W2NVGR81R06V0XGS` Fix Plan-review runs of drafts bypass the kept slots: the daemon reserves a slot for a waiting plan request or queued continuation (#226, #243), but its plan-review starts (milestone 26's drafts 01M6J375, 01M6J37D on 2026-11-04) took freed slots while the request queue's head waited. Count plan-review starts as ship picks under the same shipMax, or give queued requests priority over them
