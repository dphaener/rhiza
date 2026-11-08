# own-change-set-ops-not-stale-qa-1

evidence · `01M6XY8JH0073ZCWC7SXXCJDQR` · actor [claude-code-2ba9623a](../actor/claude-code-2ba9623a.md)

## result

PASS

## statement

QA at 90e4ec13: wrote the 01M6WMCCRC0K log shape (node.update with handle rename, then node.retire on the old landed base, no drop) with the merge-base binary in a scratch repo. Landing it with the merge-base binary fails exit 3 ('base ... is stale: op ... of change set <same> since wrote field:acceptance, field:statement, handle') and check warns. Landing the same log with the new binary exits 0, check is ok, and the node shows Retired: Gone. Cross-change-set conflicts still fire: change set B updates and retires store on the landed base, change set A then writes statement and lands first, and B's land fails exit 3 with 'stale: op ... of change set A since wrote field:statement'.

## Links

### Out

- verifies [own-change-set-ops-not-stale](../capability/rhiza/change-sets/conflicts/own-change-set-ops-not-stale.md)

## History

- `01M6XMTZ6T3CVKWG5013HWA51F` Fix A change set that updates a node and then retires it cannot land: buildRetire (internal/api/write_node.go ~356) bases the node.retire op on the node's landed version (n.Version), ignoring the update this change set has pending on that node. At land the retire then counts as a stale base against its own change set's update. Evidence: change set 01M6WMCCRC0K (ship-change conflict fix, PR #383, green) fails land with 'model: conflict: base 01M68581Q6V7XS1VHP9Q4F011R of 01M68581Q6V7XS1VHP9Q48KQJS is stale: op 01M6WMRS1ZPY of change set 01M6WMCCRC0K since wrote field:acceptance, field:statement, handle'. Op 01M6WMRS1ZPY renames and updates the node, then op 01M6WMSB28 retires it on the old base. The same happens on node 01M6A0NG..CSVTX3B (ops ..3TSKP9 then ..34BK9A). No CLI path recovers: set --rebase refuses a retired node, retire refuses an already-retired node, and carve only moves nodes created in the source. Fix both: (1) retiring a node with pending updates in the same change set drops those updates (as it already does for a pending create) and bases the retire on the landed version; (2) land/fold accept an existing log where a change set's later op's base is superseded only by its own earlier op on the same node, so 01M6WMCCRC0K lands as it is. Add tests for update-then-retire and for the existing log shape. After it ships, retry run 01M6X2EMMHKMQ9W4TDJH3ERXSD (merge main into cs/fix-ship-change-stops-on-a first).
