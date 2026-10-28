# ship-open-pr

capability · `01M5Z8XR9JBAHCDW9FXRCVJRFW` · parent [ship-through-rhiza](../ship-through-rhiza.md)

## acceptance

Tests with a fake gh and a real git remote find a first call pushing the branch and creating one PR whose title is the change set's title and whose body holds the summary's text, and a second call after a new commit pushing it and editing that PR's body without creating another; `rhiza ship check` after it reads the recorded PR; a call outside a run, or in a run with no review journaled, exits non-zero naming why.

## statement

`rhiza ship open-pr` pushes the change set's branch and opens a pull request titled by the change set, its body the run's `rhiza run summary`, or, when one is already open, pushes and replaces its body, recording the PR in the ship state file for the steps after it.

## Links

### Out

- decided-by [review-workflow-shape](../../decision/review-workflow-shape.md)
- decided-by [answer-cb8bd934-2](../../decision/answer-cb8bd934-2.md)
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
