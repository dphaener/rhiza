# retire-after-update-lands-verified

evidence · `01M6XX53QV2J59XW472915PAXQ` · actor [claude-code-f06d750a](../actor/claude-code-f06d750a.md)

## result

PASS

## statement

Verifier at 5d4cc959: read the diff. buildRetire (internal/api/write_node.go) now drops n.Pending (the change set's own updates/moves on the node) and adds one node.retire with Base = the landed version. New test TestRetire/AfterUpdateLands records a capability, then in one change set sets fields and renames the handle (store -> stash), retires it by the new handle, asserts the change set holds exactly one op, a node.retire based on the landed version (so no pending update remains on the node), then runs propose, accept and land (land runs gmodel.WhatIf and refuses on any conflict) and confirms the node is retired after land. go test -trimpath ./internal/api/... -run 'TestRetire|TestStackedWrites' -v: PASS, including TestRetire/AfterUpdateLands; the updated TestRetire and TestStackedWrites/SeesDependency expectations (retire only) also pass. Full go test -trimpath ./internal/api/...: all tests pass except the timing budgets TestReadBudget and TestReadBudgetApply, which missed their medians at VM load average 7-8 with other runs active; they time reads/apply and do not touch the retire path.

## Links

### Out

- verifies [retire-after-update-lands](../capability/rhiza/change-sets/model-writes/retire-after-update-lands.md)

## History

- `01M6XMTZ6T3CVKWG5013HWA51F` Fix A change set that updates a node and then retires it cannot land: buildRetire (internal/api/write_node.go ~356) bases the node.retire op on the node's landed version (n.Version), ignoring the update this change set has pending on that node. At land the retire then counts as a stale base against its own change set's update. Evidence: change set 01M6WMCCRC0K (ship-change conflict fix, PR #383, green) fails land with 'model: conflict: base 01M68581Q6V7XS1VHP9Q4F011R of 01M68581Q6V7XS1VHP9Q48KQJS is stale: op 01M6WMRS1ZPY of change set 01M6WMCCRC0K since wrote field:acceptance, field:statement, handle'. Op 01M6WMRS1ZPY renames and updates the node, then op 01M6WMSB28 retires it on the old base. The same happens on node 01M6A0NG..CSVTX3B (ops ..3TSKP9 then ..34BK9A). No CLI path recovers: set --rebase refuses a retired node, retire refuses an already-retired node, and carve only moves nodes created in the source. Fix both: (1) retiring a node with pending updates in the same change set drops those updates (as it already does for a pending create) and bases the retire on the landed version; (2) land/fold accept an existing log where a change set's later op's base is superseded only by its own earlier op on the same node, so 01M6WMCCRC0K lands as it is. Add tests for update-then-retire and for the existing log shape. After it ships, retry run 01M6X2EMMHKMQ9W4TDJH3ERXSD (merge main into cs/fix-ship-change-stops-on-a first).
