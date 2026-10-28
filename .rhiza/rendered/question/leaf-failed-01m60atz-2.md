# leaf-failed-01m60atz-2

question · `01M60EP8QZ10SMHX4GP950G2VJ`

## statement

Review check model still fails after 3 cycles: finding f1 (major at internal/api/run.go:450) on leaf ship-needs-leaves: Commit c2886a0 changed the refusal's next command to `rhiza run start --workflow plan`, because `rhiza run start plan` reads "plan" as a change set and fails. …. Its last cycle: cycle 3: confirmed f1 (major), f2 (major). How should it be fixed, or should leaf ship-needs-leaves change?

## Links

### Out

- blocks [ship-needs-leaves](../capability/ship-through-rhiza/ship-needs-leaves.md)

### In

- answers [needs-leaves-next-workflow-plan](../decision/needs-leaves-next-workflow-plan.md)

## History

- `01M5WPTZVV50AKZ4G5MZ7HMDKC` ship-change on Rhiza's workflows: its agent step runs Rhiza's build and review and waits for the owner's sign-off, in place of shipit.
