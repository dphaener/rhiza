# review-01m6xv4tr07qx5x8xpxh8z9m6z-r3-f1

question · `01M6XW7W0TT8W64R75667JDHZS`

## statement

Plan review round 3 left finding f1 (missing-acceptance on daemon-starts-within-max-runs): Only created capabilities are leaves (internal/api/change.go changeLeavesIn skips updates), so the new acceptance checks written on the updated milestone-request-starts-at-once (milestone request waits at --max-runs 1 with one live run; starts ahead of a resume) and request-slot-kept (milestone request plus one ship run, never three) are never run by any builder or verifier. This leaf's statement covers the milestone request's plan run waiting at the cap, but its acceptance tests only feature-request, ship and schedule starts; it should name daemon-milestone-request.txtar and daemon-request-slot.txtar with those milestone-at-cap cases (or the milestone change needs its own created leaf). The critic confirmed it: internal/api/change.go:455 changeLeavesIn takes only node-create ops of capabilities, and change_indexed.go:287 builds the change set's leaves from it, so the updated milestone-request-starts-at-once and request-slot-kept are not leaves and their new acceptance checks are never given to a builder or verifier. daemon-starts-within-max-runs's statement lists 'a queued milestone request's plan run' among the starts held at the cap, but its acceptance tests only a feature request, ship picks and a schedule, and its -run pattern leaves out daemon-milestone-request. The milestone-at-cap behaviour has no check a build would run. How should it be settled?

## Links

### Out

- blocks [daemon-starts-within-max-runs](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-starts-within-max-runs.md)

### In

- answers [milestone-cap-checks-on-max-runs-leaf](../decision/milestone-cap-checks-on-max-runs-leaf.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
