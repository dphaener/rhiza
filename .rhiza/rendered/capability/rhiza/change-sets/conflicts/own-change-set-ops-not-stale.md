# own-change-set-ops-not-stale

capability · `01M6XN774723KQDY2E0HZ0HXFJ` · parent [conflicts](../conflicts.md)

## acceptance

Tests in internal/model fold a log of the 01M6WMCCRC0K shape (one change set updates and renames a node, then retires it on the old base) with no problems and land it; still report a conflict when another change set's overlapping op sits after the base, including one interleaved between the change set's own ops; and still report the target as retired (ErrRetired or a conflict, never an applied write) when one change set retires a target and then writes it again on a base older than the retire, for both node and link ops; go test -trimpath ./internal/model/... passes.

## statement

An op is never stale against its own change set's ops: a base superseded only by the same change set's earlier ops on the target folds and lands with no conflict, while other change sets' overlapping ops after the base still conflict.

## Links

### Out

- decided-by [own-ops-never-stale](../../../../decision/own-ops-never-stale.md)

### In

- verifies [own-change-set-ops-not-stale-verified](../../../../evidence/own-change-set-ops-not-stale-verified.md)
- verifies [own-change-set-ops-not-stale-qa-1](../../../../evidence/own-change-set-ops-not-stale-qa-1.md)
- verifies [pr-394-fix-a-change-set-that-updates](../../../../evidence/pr-394-fix-a-change-set-that-updates.md)

## History

- `01M6XMTZ6T3CVKWG5013HWA51F` Fix A change set that updates a node and then retires it cannot land: buildRetire (internal/api/write_node.go ~356) bases the node.retire op on the node's landed version (n.Version), ignoring the update this change set has pending on that node. At land the retire then counts as a stale base against its own change set's update. Evidence: change set 01M6WMCCRC0K (ship-change conflict fix, PR #383, green) fails land with 'model: conflict: base 01M68581Q6V7XS1VHP9Q4F011R of 01M68581Q6V7XS1VHP9Q48KQJS is stale: op 01M6WMRS1ZPY of change set 01M6WMCCRC0K since wrote field:acceptance, field:statement, handle'. Op 01M6WMRS1ZPY renames and updates the node, then op 01M6WMSB28 retires it on the old base. The same happens on node 01M6A0NG..CSVTX3B (ops ..3TSKP9 then ..34BK9A). No CLI path recovers: set --rebase refuses a retired node, retire refuses an already-retired node, and carve only moves nodes created in the source. Fix both: (1) retiring a node with pending updates in the same change set drops those updates (as it already does for a pending create) and bases the retire on the landed version; (2) land/fold accept an existing log where a change set's later op's base is superseded only by its own earlier op on the same node, so 01M6WMCCRC0K lands as it is. Add tests for update-then-retire and for the existing log shape. After it ships, retry run 01M6X2EMMHKMQ9W4TDJH3ERXSD (merge main into cs/fix-ship-change-stops-on-a first).
