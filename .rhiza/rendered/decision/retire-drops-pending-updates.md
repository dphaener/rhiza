# retire-drops-pending-updates

decision · `01M6XN76EA3EFYEXPD998RVP8Q`

## reasoning

The change set's intent prescribes it (part 1). A retired node takes no more writes, so the dropped updates would change nothing that survives; basing the retire on the landed version keeps it from conflicting with its own change set at land.

## statement

Retiring a landed node that this change set has pending updates or moves on drops those pending ops, as retiring a node created in the change set already drops its create, and writes the node.retire on the node's landed version, so the retire is never based behind its own change set's update.

## wording

(from the change set's intent)

## Links

### Out

- derived-from [fix-01m6xmt-plan](../source/fix-01m6xmt-plan.md)

### In

- decided-by [retire-after-update-lands](../capability/rhiza/change-sets/model-writes/retire-after-update-lands.md)

## History

- `01M6XMTZ6T3CVKWG5013HWA51F` Fix A change set that updates a node and then retires it cannot land: buildRetire (internal/api/write_node.go ~356) bases the node.retire op on the node's landed version (n.Version), ignoring the update this change set has pending on that node. At land the retire then counts as a stale base against its own change set's update. Evidence: change set 01M6WMCCRC0K (ship-change conflict fix, PR #383, green) fails land with 'model: conflict: base 01M68581Q6V7XS1VHP9Q4F011R of 01M68581Q6V7XS1VHP9Q48KQJS is stale: op 01M6WMRS1ZPY of change set 01M6WMCCRC0K since wrote field:acceptance, field:statement, handle'. Op 01M6WMRS1ZPY renames and updates the node, then op 01M6WMSB28 retires it on the old base. The same happens on node 01M6A0NG..CSVTX3B (ops ..3TSKP9 then ..34BK9A). No CLI path recovers: set --rebase refuses a retired node, retire refuses an already-retired node, and carve only moves nodes created in the source. Fix both: (1) retiring a node with pending updates in the same change set drops those updates (as it already does for a pending create) and bases the retire on the landed version; (2) land/fold accept an existing log where a change set's later op's base is superseded only by its own earlier op on the same node, so 01M6WMCCRC0K lands as it is. Add tests for update-then-retire and for the existing log shape. After it ships, retry run 01M6X2EMMHKMQ9W4TDJH3ERXSD (merge main into cs/fix-ship-change-stops-on-a first).
