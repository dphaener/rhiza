# milestone-cap-checks-on-max-runs-leaf

decision · `01M6XWA2TPB8VVVZW9FNFJ0XX7`

## reasoning

The finding is confirmed: changeLeavesIn (internal/api/change.go) takes only created capabilities, so the acceptance updates on milestone-request-starts-at-once and request-slot-kept reach no builder or verifier. daemon-starts-within-max-runs is the created leaf whose statement already lists 'a queued milestone request's plan run' among the starts held at the cap, and plan-runs-share-max-runs (3) makes the milestone wait at the cap part of the same one machine-wide count this leaf implements; the code change is one slot count, so splitting the milestone checks into a second created leaf would duplicate a capability and split one implementation across two builders. Both testscripts already exist (daemon-milestone-request.txtar, daemon-request-slot.txtar), so naming them with the exact cases from the updated capabilities' acceptance, and adding daemon-milestone-request to the -run pattern (daemon-request-slot is already in it), is the smallest edit that gives every new check a build that runs it.

## statement

Settle it on the existing leaf: @{01M6XVG9DK25PFK9BHJCBDQHQ6}'s acceptance takes over the milestone-at-cap checks written on @milestone-request-starts-at-once and @request-slot-kept. It names cmd/rhiza/testdata/script/daemon-milestone-request.txtar (with --max-runs 1 and one run live, daemon --once starts nothing and the milestone and feature requests stay queued; with one interrupted run and one free slot the milestone request starts and the run is not resumed) and daemon-request-slot.txtar (with max-runs 2, two ready change sets and only a queued milestone request, daemon --once starts the milestone request's plan run and one ship run, never three, and daemon status shows the other as no-slot), and its go test -run pattern adds daemon-milestone-request. No new created leaf is added for the milestone change.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6xv4tr07qx5x8xpxh8z9m6z-r3-f1](../question/review-01m6xv4tr07qx5x8xpxh8z9m6z-r3-f1.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
