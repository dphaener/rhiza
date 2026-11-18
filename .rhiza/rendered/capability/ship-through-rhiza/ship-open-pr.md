# ship-open-pr

capability · `01M5Z8XR9JBAHCDW9FXRCVJRFW` · parent [ship-through-rhiza](../ship-through-rhiza.md)

## acceptance

go test -trimpath -p 2 ./internal/api -run 'TestShipOpenPR|TestReviewJournaled|TestReviewDeclared' and go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/(ship-open-pr.*|run-ship-change)$' pass: with a fake gh and a real git remote, a first call pushes the branch and creates one PR titled by the change set with the summary's text as body, a second call after a new commit pushes and edits that PR's body without creating another, and rhiza ship check after it reads the recorded PR; a run whose workflow declares no review (no checks step, no workflow step running a review, no agent step with [[check]] tables; e.g. myteam-flow from docs/examples/byo/workflow.toml) opens the PR with no review journaled; (failure case) a call outside a run, or in a run whose workflow declares a review (ship-change, bugfix, ship-change-classic) with none journaled, exits non-zero naming why

## files

internal/api/ship_openpr.go, internal/api/ship_openpr_test.go, internal/workflow/workflow.go (a declared-review test on the workflow's shape), internal/workflow/workflow_test.go, docs/reference (regenerated with rhiza docs if a description changes)

## must_not

Decide by a workflow's or step's name; drop the guard for a workflow that declares a review; add a workflow file key

## statement

`rhiza ship open-pr` pushes the change set's branch and opens a pull request titled by the change set, its body the run's `rhiza run summary`, or, when one is already open, pushes and replaces its body, recording the PR in the ship state file for the steps after it.

## Links

### Out

- decided-by [review-workflow-shape](../../decision/review-workflow-shape.md)
- decided-by [answer-cb8bd934-2](../../decision/answer-cb8bd934-2.md)
- decided-by [open-pr-review-only-when-declared](../../decision/open-pr-review-only-when-declared.md)
- derived-from [ship-change-plan](../../source/ship-change-plan.md)

### In

- depends-on [ship-signoff](ship-signoff.md)
- depends-on [ship-change-file](ship-change-file.md)
- verifies [ship-open-pr-verified](../../evidence/ship-open-pr-verified.md)
- verifies [ship-open-pr-qa-1](../../evidence/ship-open-pr-qa-1.md)
- verifies [ship-open-pr-qa-3](../../evidence/ship-open-pr-qa-3.md)
- verifies [ship-open-pr-qa-4](../../evidence/ship-open-pr-qa-4.md)
- verifies [ship-open-pr-qa-5](../../evidence/ship-open-pr-qa-5.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
