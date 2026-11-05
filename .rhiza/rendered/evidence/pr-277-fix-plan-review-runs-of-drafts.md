# pr-277-fix-plan-review-runs-of-drafts

evidence · `01M6MTTV1YP9KGCYDV6P4DTVTC` · actor [dphaener](../actor/dphaener.md)

## head

fbf42a0fce883f016b083e62b58422ed8cebbd3f

## pr

https://github.com/dphaener/rhiza/pull/277

## result

1 check passed: check SUCCESS

## statement

Pull request #277 (https://github.com/dphaener/rhiza/pull/277) passed CI at fbf42a0

## status

verified

## Links

### Out

- verifies [queued-start-request-slot-kept](../capability/rhiza/orchestration/daemon/plan-requests/queued-start-request-slot-kept.md)

## History

- `01M6M3VQ85W2NVGR81R06V0XGS` Fix Plan-review runs of drafts bypass the kept slots: the daemon reserves a slot for a waiting plan request or queued continuation (#226, #243), but its plan-review starts (milestone 26's drafts 01M6J375, 01M6J37D on 2026-11-04) took freed slots while the request queue's head waited. Count plan-review starts as ship picks under the same shipMax, or give queued requests priority over them
