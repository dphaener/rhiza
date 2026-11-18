# byo-open-pr-needs-review

question · `01M7RBEHJ5MRCPKYMZGTTA2XG3`

## statement

The BYO example myteam-flow (docs/examples/byo/workflow.toml) cannot push: its push step runs rhiza ship open-pr, which refuses a run with no review journaled ("run <id> has no review journaled: ship open-pr opens the pull request once the review has run", internal/api/ship_openpr.go reviewJournaled), and myteam-flow has no checks step and its /myteam:implement invoke never records findings. Shown by a testscript run of the example: services, work, gate (with a lint hand-back) and signoff pass, then push fails and hands back to work. Which fix? (a) open-pr requires a journaled review only when the run's workflow declares a review (a checks step or conduct checks), so a team workflow with no review pushes; (b) the BYO example gains a review (a checks step, or its invoke tells the skill to record a round with rhiza run findings), which this leaf's must_not forbids; (c) a workflow-level opt-out (e.g. review = false) open-pr honours. byo-workflow-daemon-e2e cannot pass as written until one lands.

## Links

### Out

- blocks [byo-workflow-daemon-e2e](../capability/rhiza/workflows/workflow-runs/byo-workflow-daemon-e2e.md)

### In

- answers [open-pr-review-only-when-declared](../decision/open-pr-review-only-when-declared.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
