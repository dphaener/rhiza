# plan-review-done-unchained

capability · `01M76F8S79W69F3VS3VY9VTKYX` · parent [plan-workflow](../plan-workflow.md)

## acceptance

Go tests: internal/guide (plan_test.go) checks the plan-review topic says to run session done separately after the recording command succeeds and never to chain it with ';', staying within TopicCap; an internal/workflow builtin test checks the plan and plan-milestone review, critic and revise prompts carry the same instruction. Run: go test -trimpath ./internal/guide/... ./internal/workflow/... -run <the new tests>.

## statement

rhiza guide plan-review and the review, critic and revise step prompts of the built-in plan and plan-milestone workflows tell the agent to run rhiza session done as a separate command only once rhiza run findings, rhiza run verdict or the reviser's apply succeeded, never chained after one with ';'.

## Links

### Out

- decided-by [session-done-not-chained](../../../../../decision/session-done-not-chained.md)
- derived-from [plan-review-round-unrecorded-plan](../../../../../source/plan-review-round-unrecorded-plan.md)

### In

- verifies [plan-review-done-unchained-verified](../../../../../evidence/plan-review-done-unchained-verified.md)
- verifies [plan-review-done-unchained-qa-1](../../../../../evidence/plan-review-done-unchained-qa-1.md)
- verifies [pr-452-fix-a-plan-review-round-whose](../../../../../evidence/pr-452-fix-a-plan-review-round-whose.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
