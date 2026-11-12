# plan-review-round-unrecorded-plan

source · `01M76F79565HYN8J26RM52EY8C`

## quotes

```json
[]
```

## raw_sha256

7876b1804c1e2e43b2bceb61e55421373c1fc655a616e0bcca69760302f4de64

## session

d97fabfd-9ac2-42dd-bb6c-01a1ccf27f7d

## summary

Autonomous plan-bug run, no owner. The bug request (plan run 01M76583, 2026-11-11) names the failure: a round-3 reviewer's rhiza run findings was refused because a finding named a depends-on link's ID, the reviewer had chained rhiza session done with ';', so the review step passed with no round journaled; review-check then judged round 2 again and revise exceeded max_visits. The planner took the request's four-part fix as the plan, serving m4-runs-hold-up (runs finish without hand repair) as earlier run-repair bug fixes do: (1) the review step passes only when its visit journaled a new findings round, and session done refuses for a reviewer without one; (2) review-check refuses to judge a round that predates the latest revise; (3) run findings accepts a link as a finding's target; (4) the plan-review guide tells every role not to chain session done after a command that can fail. No owner words; no question was raised.

## title

Plan of the bug fix for unrecorded plan review rounds (autonomous plan-bug run)

## Links

### In

- derived-from [review-step-needs-a-round](../decision/review-step-needs-a-round.md)
- derived-from [review-visit-records-round](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-visit-records-round.md)
- derived-from [stale-round-not-judged](../decision/stale-round-not-judged.md)
- derived-from [review-check-refuses-stale-round](../capability/rhiza/workflows/built-in-workflows/plan-workflow/review-check-refuses-stale-round.md)
- derived-from [link-findings-allowed](../decision/link-findings-allowed.md)
- derived-from [findings-name-links](../capability/rhiza/workflows/built-in-workflows/plan-workflow/findings-name-links.md)
- derived-from [session-done-not-chained](../decision/session-done-not-chained.md)
- derived-from [plan-review-done-unchained](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-done-unchained.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
