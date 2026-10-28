# needs-leaves-next-workflow-plan

decision · `01M61BG9TB3SH2PK9WTA1GKPJF`

## reasoning

QA confirmed in cycle 1 that `rhiza run start plan` reads plan as a change set prefix and fails, and c2886a0 fixed the code. The leaf's wording was shorthand for the plan workflow, not a literal argv, so revising the leaf to the working command keeps its intent and makes the model match the code and Evidence ship-needs-leaves-verified-3.

## statement

Leaf ship-needs-leaves' statement and acceptance name `rhiza run start --workflow plan` as the refusal's next command, matching the code at internal/api/run.go:450. The code stays as it is.

## wording

decided by the orchestrator acting as arbiter, on the owner's word that the orchestrator signs off on this change set

## Links

### Out

- answers [leaf-failed-01m60atz-2](../question/leaf-failed-01m60atz-2.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
