# open-pr-review-when-declared

capability · `01M7RBM7A9D06RDP95JE3DJCXW` · parent [ship-through-rhiza](../ship-through-rhiza.md)

## acceptance

go test -trimpath -p 2 ./internal/api -run 'TestShipOpenPR|TestReviewJournaled|TestReviewDeclared' and go test -trimpath -p 2 ./internal/workflow pass, plus env -u'd of RHIZA_*/SHIP_*/REX_* go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/run-ship-change$': a run whose workflow declares no review (myteam-flow from docs/examples/byo/workflow.toml, a review-less workflow in a Go test) opens the PR with no review journaled; the declared-review test is true for ship-change (conduct checks), ship-change-classic and review (a review workflow step or checks step) and bugfix; (failure case) a run whose workflow declares a review with none journaled still exits non-zero naming why, and a call outside a run is still refused

## files

internal/api/ship_openpr.go, internal/api/ship_openpr_test.go, internal/workflow/workflow.go, internal/workflow/workflow_test.go, docs/reference (only if rhiza docs regenerates a changed description)

## must_not

Decide by a workflow's or step's name; drop the guard for a workflow that declares a review; add a workflow file key; change docs/examples/byo/*.toml

## statement

rhiza ship open-pr requires a journaled review only when the run's workflow declares one (a checks step, a workflow step running a review, or an agent step with [[check]] tables), judged by the workflow's declared shape and never a name, so a workflow with no review such as the BYO example myteam-flow opens its PR after its own steps pass

## Links

### Out

- decided-by [open-pr-review-only-when-declared](../../decision/open-pr-review-only-when-declared.md)

### In

- depends-on [byo-workflow-daemon-e2e](../rhiza/workflows/workflow-runs/byo-workflow-daemon-e2e.md)
- verifies [open-pr-review-when-declared-verified](../../evidence/open-pr-review-when-declared-verified.md)
- verifies [open-pr-review-when-declared-qa-1](../../evidence/open-pr-review-when-declared-qa-1.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
