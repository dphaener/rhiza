# ship-land-records-model-conflict

capability · `01M77NA4P0Z0VX0NW4RM4J29Z5` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

Go tests in internal/api: a ship land whose change set updates a node field that another change set landed a write to since its base is refused and records the conflict (node, op, other change set, aspects, head) in the ship state; the same land through the head-build path (a fake head build via ctx.Exec, as in TestShipLandWithHeadBuild, that exits 3 printing the --json error body rhiza land gives for that conflict) records the same conflict; rhiza land --json on a model conflict prints an error body whose structured fields name the node, op, other change set and aspects; ship check --model-conflict passes on it, fails with no recorded conflict, fails once the branch moved past the recorded head, fails once the branch holds 3 commits with the run's Review: and Model-Conflict: trailers, and fails after a land refused for a missing dependency (in-process or through the head build); ship check --json shows both change sets with their intents; after rhiza set --rebase in the change set, ship land succeeds and clears the record. Run with go test -trimpath ./internal/api/... -run 'Ship|Land' only.

## statement

When rhiza ship land is refused because the change set's what-if view has a model conflict (rhiza land exit 3 on a stale base: another change set landed a write to the same node aspect), it records the conflict in the change set's ship state (the node, the conflicting op and its change set, the aspects, the head it was refused at) and names it in its error. This holds both when ship land lands in-process and when it lands through the PR head's build (Rhiza's own repo, where shipLandWith runs `rhiza land --change <cs> --json` as a subprocess): rhiza land's --json error body for a model conflict carries the conflict as structured fields (node, op, other change set, aspects) beside its code and message, and ship land reads them from the child's body rather than from headFailure's flattened text. rhiza ship check --model-conflict exits 0 only when the ship state records such a conflict, the branch is still at that head and the run has made fewer than 3 model-conflict fixes (commits on the branch carrying this run's Review: trailer and a Model-Conflict: trailer), else exits 1 naming why, and rhiza ship check --json shows the recorded conflict with both change sets' IDs and intents for the fixer; a land refused for any other reason, or a later land that succeeds, leaves no conflict recorded.

## Links

### Out

- decided-by [ship-land-model-conflict-rebase](../../../../../decision/ship-land-model-conflict-rebase.md)

### In

- depends-on [ship-change-model-conflict-fix-step](ship-change-model-conflict-fix-step.md)
- verifies [ship-land-records-model-conflict-verified](../../../../../evidence/ship-land-records-model-conflict-verified.md)
- verifies [ship-land-records-model-conflict-qa-1](../../../../../evidence/ship-land-records-model-conflict-qa-1.md)
- verifies [pr-484-fix-watchdog-unmerged-pr-at-01m77n1gagbp](../../../../../evidence/pr-484-fix-watchdog-unmerged-pr-at-01m77n1gagbp.md)

## History

- `01M77N1GAGBPBP7EH4EV2CT023` Fix watchdog unmerged PR at step land: run 01M7786QKAZMH7PBJZP7QQA6HX; log: R head's build: conflict: cannot land change set 01M763YX47S5XY0W3P8FVJ54S8: model: conflict: base 01M6RN85C17S5VE1A6SQVD7WE1 of 01M6RN85C17S5VE1A6SNZAXK89 is stale: op 01M77C4T7EBWGBK547B9S31A74 of change set 01M763Z16T0ZGQMY1EC2Y81V0K since wrote field:statement: /home/rhiza/code/rhiza/.git/rhiza/ship/01M763YX47S5XY0W3P8FVJ54S8/rhiza land --change 01M763YX47S5XY0W3P8FVJ54S8 --json: exit status 3
