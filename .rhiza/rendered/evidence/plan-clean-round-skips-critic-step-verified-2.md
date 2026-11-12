# plan-clean-round-skips-critic-step-verified-2

evidence · `01M77TYB60MB1GYZM2RD7DQXBT` · actor [claude-code-c663c2ab](../actor/claude-code-c663c2ab.md)

## result

PASS

## statement

Verifier, attempt 2, at a511eec72 (detached checkout). The fix commit updates TestPlanBuiltin to the plan order review -> has-findings -> critic/fresh, which resolves the quality finding that internal/workflow was red. Ran with GOTMPDIR outside the worktree path: go test -count=1 -trimpath ./internal/workflow/... ok (incl. TestPlanBuiltin, TestPlanCleanRoundSkipsCritic, which reads both built-in plan and plan-milestone); go test -trimpath ./internal/api -run 'TestReviewCheckFindingsFlag|TestReviewCheck' PASS (clean or absent round fails review-check --findings, a round with a finding passes, run verdict with {"verdicts":[]} errors); go test -count=1 -trimpath ./cmd/rhiza -run TestScripts/(plan-workflow|plan-workflow-autonomous|plan-workflow-pr|plan-milestone-backlog|plan-review-visit-round|describe-golden) ok (plan-workflow-autonomous: review -> has-findings fail -> fresh -> check -> gate with no critic session; plan-workflow, plan-review-visit-round and plan-milestone-backlog: a round with findings goes to the critic). Diff read: both plan.toml and plan-milestone.toml add a has-findings command step running rhiza run review-check --findings, with on.pass critic and on.fail fresh. Note: plan-review-visit-round fails at line 46 only when GOTMPDIR contains the word critic (this checkout's path does), because $WORK appears in run status output. That comes from the environment, not from the change.

## Links

### Out

- verifies [plan-clean-round-skips-critic-step](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-clean-round-skips-critic-step.md)

## History

- `01M77NRZEXWM1MTHM3AWYJFMQE` Fix plan workflows visit plan/critic after a review round with no findings, but rhiza run verdict refuses an empty verdicts list ('at least one'), so the critic cannot record anything and its prompt ('session done only once run verdict succeeded') cannot be satisfied. Seen in run 01M77N1ENKEAY668AM2STW90V7 (change set 01M77N1GAGBPBP7EH4EV2CT023), round 2. Either route plan/review to plan/fresh when the round is clean, or let run verdict accept an empty list for an empty round.
