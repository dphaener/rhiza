# workflow-step

capability · `01M5Z8XQV4XPRTXNA7KE1AKET4` · parent [ship-through-rhiza](../ship-through-rhiza.md)

## acceptance

Parse tests reject a workflow step with no `workflow`, one naming an unknown workflow, a workflow that includes itself directly or through another, `skip` naming a step the inner workflow lacks, and `workflow` or `skip` on another step kind; an engine test with a fake runner over an outer workflow whose workflow step names a three-step inner workflow skipping one finds the other two journaled in order under the inner workflow's name, the skipped one not run and the outer on.pass taken; an inner step failing fails the outer step, and `rhiza run resume` after it resumes at that inner step.

## statement

Workflow files take a `workflow` step kind naming another workflow, resolved as `rhiza run start` resolves it (a repo override included), and an optional `skip` list of its steps; the engine runs that workflow's steps inside the same run and journal, each skipped step taken as passed, and the step passes when the inner workflow ends done and fails when it ends failed.

## Links

### Out

- decided-by [engine-is-conductor](../../decision/engine-is-conductor.md)
- decided-by [answer-cb8bd934-1](../../decision/answer-cb8bd934-1.md)
- derived-from [ship-change-plan](../../source/ship-change-plan.md)

### In

- depends-on [ship-change-file](ship-change-file.md)
- verifies [workflow-step-verified](../../evidence/workflow-step-verified.md)
- verifies [workflow-step-qa-1](../../evidence/workflow-step-qa-1.md)
- verifies [workflow-step-qa-3](../../evidence/workflow-step-qa-3.md)
- verifies [workflow-step-qa-4](../../evidence/workflow-step-qa-4.md)
- verifies [workflow-step-qa-5](../../evidence/workflow-step-qa-5.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
