# review-check-refuses-stale-round

capability · `01M76F8RTCWVC288HQ1BJSJ2PE` · parent [plan-workflow](../plan-workflow.md)

## acceptance

A test in internal/api (run_review_test.go) with a journal holding round 2 with a confirmed finding, then a revise step visit and no later round: review-check and review-check --fresh exit non-zero, the message names round 2 as predating the revise, and neither adds a question or asked entry, even when the round is 3; with a round 3 journaled after the revise visit both exit 0 for --fresh and review-check judges round 3 as before. An internal/workflow builtin test checks that plan and plan-milestone have the fresh step after critic running rhiza run review-check --fresh with on.pass = check and on.fail = review, and that critic's on.pass is fresh. The describe-golden testscript is regenerated for the new flag. Run: go test -trimpath ./internal/api ./internal/workflow/... -run <the new tests>, and go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden.

## statement

rhiza run review-check refuses to judge, exiting non-zero and writing nothing, when the run's latest findings round was journaled before its latest revise visit began, saying that round predates the revise and the review recorded no round since; rhiza run review-check --fresh makes only that test, exiting 0 when the latest round postdates the latest revise visit (or the run has none). The built-in plan and plan-milestone workflows run it as a command step fresh between critic and check, with on.pass = check and on.fail = review, so a stale round sends the run back to a fresh review, never to revise, and check judges only a fresh round.

## Links

### Out

- decided-by [stale-round-not-judged](../../../../../decision/stale-round-not-judged.md)
- derived-from [plan-review-round-unrecorded-plan](../../../../../source/plan-review-round-unrecorded-plan.md)

### In

- verifies [review-check-refuses-stale-round-verified](../../../../../evidence/review-check-refuses-stale-round-verified.md)
- verifies [review-check-refuses-stale-round-qa-1](../../../../../evidence/review-check-refuses-stale-round-qa-1.md)
- verifies [pr-452-fix-a-plan-review-round-whose](../../../../../evidence/pr-452-fix-a-plan-review-round-whose.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
