# plan-review-done-unchained-verified

evidence · `01M76J9VEWYA0BRSQZ6F61M206` · actor [claude-code-9d42dfd1](../actor/claude-code-9d42dfd1.md)

## result

PASS

## statement

At b364a6a7b (attempt 2), ran go test -trimpath -count=1 ./internal/guide/... ./internal/workflow/... -run Unchained: TestPlanReviewDoneUnchained (plan-review topic says session done runs as its own command only after rhiza run findings / run verdict / the reviser's apply succeeded, never chained with ';', and stays within TopicCap) and TestPlanReviewPromptsUnchained (plan and plan-milestone review, critic and revise prompts carry 'as a separate command', the recording command, 'succeeded' and 'never chained after it with `;`') both PASS. Read the diffs of 5a50d91ce and b364a6a7b: guide text and both builtin TOMLs carry the instruction. Attempt 1's finding is fixed: plan-workflow, plan-workflow-pr and plan-milestone-backlog now assert the new reviewer, critic and reviser prompts (and the critic->fresh->check steps), and with GOTMPDIR overridden to a scratch dir, go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-workflow|plan-workflow-pr|plan-milestone-backlog|plan-workflow-ci-infra-wait|describe-golden)$' all PASS.

## Links

### Out

- verifies [plan-review-done-unchained](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-review-done-unchained.md)

## History

- `01M76F4G0JQ649RR29M9QPFX68` Fix A plan review round whose findings fail to record still passes the review step, so review-check judges the previous round again and the plan run fails at revise's max_visits with its fixes in. Evidence 2026-11-11, plan run 01M76583MQA31HW51PH7396Q4Y (milestone 33, release channels, change set 01M765854D7DGJDBW5CYGJ4SR1): round 2 (22:34Z) had f1 release-provenance and f2 install-script; revise visit 2 fixed both (ops 22:37-22:38Z in 01M76AZ94M8HWY6GCD06H8ES46 and the install-script change set); the round-3 reviewer (session 1f769929-ee04-4c59-961f-01a1ccc326ee) found that round 1's confirmed f6 (an unneeded depends-on link) was never retired, but 'rhiza run findings' rejected its finding because it named the link's ID ('it accepts only node IDs'), and the reviewer had chained 'rhiza session done' with ';', so the session ended and the review step passed ('said done and until held') with no round 3 journaled. check visit 3 then read LatestRound (internal/api/run_review.go:991), which was still round 2, reported its already-fixed findings, and revise went over max_visits 2: 'run 01M76583 failed: step "revise" went over max_visits 2'. Fix: (1) the review step's until (and session done for a reviewer) requires that this visit journaled a new findings round, empty for a clean review, else the step does not pass; (2) review-check refuses to judge when the latest round predates the latest revise; (3) run findings accepts a link (or change-set link) as a finding's target, since unreal-depends-on findings are about links; (4) the review prompt says not to chain session done after a command that can fail. Unstuck by rhiza run resume --retry 01M76583.
