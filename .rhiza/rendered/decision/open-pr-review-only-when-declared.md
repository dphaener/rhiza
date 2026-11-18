# open-pr-review-only-when-declared

decision · `01M7RBJ05HVWR5VPAV1XQAZ8XQ`

## reasoning

The change set's intent is that a bring-your-own workflow runs end to end under the daemon, and the owner's words (source byo-daemon-plan-research) ask that a BYO workflow run end to end. The leaf byo-workflow-daemon-e2e forbids changing docs/examples/byo/*.toml, which rules out (b), and the example is the documented level-3 team workflow, so making every team add a Rhiza review would contradict bring-your-own. (c) adds a new workflow field and file-format surface only to restate what the workflow already says by having no review step. The open-pr guard (ship-open-pr, decided by review-workflow-shape) exists so the built-in review workflow cannot be skipped before a PR; tying it to whether the workflow declares a review keeps that protection for every workflow that has one, and myteam-flow still has its own gate and owner sign-off before push. Deciding by declared step kinds rather than names honours the leaf's must_not on special-casing a workflow or step name.

## statement

Option (a): rhiza ship open-pr requires a journaled review only when the run's workflow declares one (a checks step, a workflow step running a review, or conduct checks). A workflow that declares no review, such as the BYO example myteam-flow, pushes after its own steps (gate, signoff) pass. A workflow that declares a review still refuses open-pr until it is journaled, so ship-change, bugfix and the other built-ins keep today's guard. The test is the workflow's declared shape, never its or a step's name. ship-open-pr's acceptance is reworded to match: a run whose workflow declares a review and has none journaled exits non-zero naming why; a run whose workflow declares none opens the PR.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-2debf9c0-1](../question/question-2debf9c0-1.md)
- answers [byo-open-pr-needs-review](../question/byo-open-pr-needs-review.md)

### In

- decided-by [ship-open-pr](../capability/ship-through-rhiza/ship-open-pr.md)
- decided-by [open-pr-review-when-declared](../capability/ship-through-rhiza/open-pr-review-when-declared.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
