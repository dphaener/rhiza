# ship-change-file-names-render

decision · `01M61BGA5SNFJNDHFHC1SX9B4W`

## reasoning

The model check confirmed that 179a96f added a render step the leaf's step-by-step statement leaves out. The step is needed (QA's cycle-1 blocker), so the leaf, not the code, changes. It is revised together with the fix for leaf-failed-01m60atz-1 so that the statement describes the final step.

## statement

Leaf ship-change-file's statement lists the render command step between the build and review workflow steps: it runs rhiza render and commits .rhiza/rendered and the uncommitted .rhiza/log entries on the branch (see render-commits-actor-records).

## wording

decided by the orchestrator acting as arbiter, on the owner's word that the orchestrator signs off on this change set

## Links

### Out

- answers [leaf-failed-01m60atz-3](../question/leaf-failed-01m60atz-3.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
