# plan-clean-round-skips-critic-step-verified

evidence · `01M77SE7P5TRT5C3EDRDJ7P5B3` · actor [claude-code-69f6e1c7](../actor/claude-code-69f6e1c7.md)

## result

PASS

## statement

Verifier at cf18b5eb9: read the diff (plan.toml and plan-milestone.toml route review on.pass to a new has-findings command step running `rhiza run review-check --findings`, pass->critic, fail->fresh; run_review.go adds the --findings test). Ran go test -trimpath ./internal/api/... -run TestReviewCheckFindingsFlag (no round and a clean {"findings":[]} round fail --findings, a round with a finding passes, run verdict with {"verdicts":[]} still errors): PASS. go test -trimpath ./internal/workflow/... -run TestPlanCleanRoundSkipsCritic (reads both built-in plan and plan-milestone files): PASS. go test -trimpath ./cmd/rhiza -run TestScripts/(plan-review-visit-round|plan-workflow|plan-workflow-pr|plan-workflow-autonomous|plan-milestone-backlog|describe-golden): PASS; plan-workflow shows review visit 2 (clean round) -> has-findings fail -> fresh -> check pass -> gate with no critic visit 2 and one critic prompt, while visit 1 with findings still goes to the critic. describe-golden regenerated for the new flag.

## Links

### Out

- verifies [plan-clean-round-skips-critic-step](../capability/rhiza/workflows/built-in-workflows/plan-workflow/plan-clean-round-skips-critic-step.md)

## History

- `01M77NRZEXWM1MTHM3AWYJFMQE` Fix plan workflows visit plan/critic after a review round with no findings, but rhiza run verdict refuses an empty verdicts list ('at least one'), so the critic cannot record anything and its prompt ('session done only once run verdict succeeded') cannot be satisfied. Seen in run 01M77N1ENKEAY668AM2STW90V7 (change set 01M77N1GAGBPBP7EH4EV2CT023), round 2. Either route plan/review to plan/fresh when the round is clean, or let run verdict accept an empty list for an empty round.
