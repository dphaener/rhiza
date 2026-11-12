# review-visit-records-round

capability · `01M76F8RPF7AF6C8Y5Y9D11DRW` · parent [plan-workflow](../plan-workflow.md)

## acceptance

New tests in internal/run and/or internal/api, plus a testscript under cmd/rhiza/testdata/script: a plan run's review-step session that says done with no rhiza run findings since the visit began does not pass to the critic and session done exits non-zero naming run findings; after rhiza run findings records {"findings":[]} in the same visit, session done exits 0 and the step passes to the critic; a round journaled in an earlier visit does not count. Run: go test -trimpath ./internal/run/... ./internal/api/... -run <the new tests>, and go test -trimpath ./cmd/rhiza -run TestScripts/<the new script>.

## statement

In the plan workflows (plan and plan-milestone), the review step passes only when the reviewer's visit journaled a new findings round (empty for a clean review), and rhiza session done in a review-step session that has journaled no round since the visit began exits non-zero naming rhiza run findings, so the session stays open to record one.

## Links

### Out

- decided-by [review-step-needs-a-round](../../../../../decision/review-step-needs-a-round.md)
- derived-from [plan-review-round-unrecorded-plan](../../../../../source/plan-review-round-unrecorded-plan.md)

### In

- verifies [review-visit-records-round-verified](../../../../../evidence/review-visit-records-round-verified.md)
- verifies [review-visit-records-round-qa-1](../../../../../evidence/review-visit-records-round-qa-1.md)
- verifies [pr-452-fix-a-plan-review-round-whose](../../../../../evidence/pr-452-fix-a-plan-review-round-whose.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
