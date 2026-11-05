# queued-start-request-slot-kept

capability · `01M6M438EJP4VAN9P056CW2BNA` · parent [plan-requests](../plan-requests.md)

## acceptance

daemon-request-slot.txtar (or a new daemon-queued-start-slot.txtar) with a fake rex and no ready change sets: with max-runs 2, two runs queued with rhiza run start on draft change sets and no request, daemon --once starts both; with the same two queued and one queued request, daemon --once starts one queued run and the request's plan run, and the other queued start stays in the queue (it starts on a later poll once a slot frees); with max-runs 1, one queued start and a queued request, it starts the queued run and the request stays queued

## statement

While a queued plan request could start (daemonRequestWaits: no run, its retry time passed) and max-runs is above 1, the daemon starts owner-queued runs (rhiza run start, plan and plan-review runs of drafts included) only while fewer than max-runs minus 1 runs run, so one slot stays for the request's plan run; a queued start held back stays queued for a later poll. The poll order (ship picks, queued continuations, queued starts, plan requests) is unchanged, and with max-runs 1 or no request waiting, queued starts fill every free slot as before

## Links

### Out

- decided-by [queued-starts-keep-request-slot](../../../../../decision/queued-starts-keep-request-slot.md)

### In

- verifies [queued-start-request-slot-kept-verified](../../../../../evidence/queued-start-request-slot-kept-verified.md)
- verifies [queued-start-request-slot-kept-qa-1](../../../../../evidence/queued-start-request-slot-kept-qa-1.md)
- verifies [pr-277-fix-plan-review-runs-of-drafts](../../../../../evidence/pr-277-fix-plan-review-runs-of-drafts.md)

## History

- `01M6M3VQ85W2NVGR81R06V0XGS` Fix Plan-review runs of drafts bypass the kept slots: the daemon reserves a slot for a waiting plan request or queued continuation (#226, #243), but its plan-review starts (milestone 26's drafts 01M6J375, 01M6J37D on 2026-11-04) took freed slots while the request queue's head waited. Count plan-review starts as ship picks under the same shipMax, or give queued requests priority over them
