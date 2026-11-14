# review-three-checks

capability · `01M79KPDAWDKT4JQK1KDNR5K0W` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

go test -trimpath ./internal/workflow -run TestAgentStepChecks shows an agent step's declared checks parsed and rhiza workflow check refusing a check with no name or a duplicate name; go test -trimpath ./internal/api -run TestRunFindingsCheckNames shows, in a run in an agent step declaring plan, code and qa, findings and verdict with --check plan recorded as plan's rounds, and --check scope refused with a message naming plan, code and qa. Touches internal/workflow/workflow.go, internal/workflow/workflow_test.go, internal/api/run_review.go, internal/api/run_review_test.go.

## statement

An agent step may declare its review checks in the workflow definition the way a checks step does ([[step.<name>.check]] with a name); rhiza run findings --check and rhiza run verdict --check accept a run in such a step and validate the name against that step's declared checks, refusing any other (scope included) with a message naming them; nothing hard-codes plan, code and qa.

## Links

### Out

- decided-by [conductor-three-checks](../../../../decision/conductor-three-checks.md)
- decided-by [conductor-checks-gate-pick](../../../../decision/conductor-checks-gate-pick.md)

### In

- depends-on [critic-on-major](critic-on-major.md)
- depends-on [conduct-plans-leaves](../../orchestration/daemon/conduct-plans-leaves.md)
- verifies [review-three-checks-verified](../../../../evidence/review-three-checks-verified.md)
- verifies [review-three-checks-qa-1](../../../../evidence/review-three-checks-qa-1.md)
- verifies [review-three-checks-qa-2](../../../../evidence/review-three-checks-qa-2.md)
- verifies [pr-575-a-rhiza-conduct-skill-conducts](../../../../evidence/pr-575-a-rhiza-conduct-skill-conducts.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
