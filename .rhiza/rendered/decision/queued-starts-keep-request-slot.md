# queued-starts-keep-request-slot

decision · `01M6M41F92YJYEJYTF9XMK2NKM`

## reasoning

In daemonPoll1 (internal/api/daemon.go) ship picks are capped at shipMax = max - kept, with one slot kept when max > 1 and a request waits (#226) and one per startable continuation (#243). daemonQueuedContinues then runs before daemonQueuedStarts, so continuations already take their kept slots; only the plan-request slot is unprotected, because daemonQueuedStarts fills up to max before daemonPlanRequests runs. Applying the same max-1 cap to queued starts reuses the existing rule and its one predicate, is bounded (a queued start loses at most one slot, never starves) and fixes every queued start, not only plan-review ones. The other option, running daemonPlanRequests before queued starts, would let several requests take every free slot (it loops while running < max) and starve the owner's own queued starts, which the kept-slot model was built to avoid. Counting plan-review starts as ship picks would need the daemon to tell plan-review workflows apart from other owner starts, a special case the one-slot rule does not need. For the intent: this is a daemon scheduling fix found while emptying the request queue, like 01M6M1M1 (Fix Test temp folders leak), which serves milestone-25 (the daemon builds and ships the planned backlog on its own); a plan request starved of its slot blocks exactly that, so milestone-25 is the fit.

## statement

Owner-queued starts (daemonQueuedStarts, which includes the plan-review runs) leave the same kept slot ship picks leave: while a queued plan request could start (daemonRequestWaits) and max > 1, queued starts start only while running < max-1, so one slot stays for daemonPlanRequests. Queued continuations keep running before them, as now. The poll order (picks, continuations, queued starts, plan requests) is unchanged. This change set serves milestone-25.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-62e3cc37-1](../question/question-62e3cc37-1.md)
- derived-from [plan-review-runs-slot-interview](../source/plan-review-runs-slot-interview.md)

### In

- decided-by [queued-start-request-slot-kept](../capability/rhiza/orchestration/daemon/plan-requests/queued-start-request-slot-kept.md)

## History

- `01M6M3VQ85W2NVGR81R06V0XGS` Fix Plan-review runs of drafts bypass the kept slots: the daemon reserves a slot for a waiting plan request or queued continuation (#226, #243), but its plan-review starts (milestone 26's drafts 01M6J375, 01M6J37D on 2026-11-04) took freed slots while the request queue's head waited. Count plan-review starts as ship picks under the same shipMax, or give queued requests priority over them
