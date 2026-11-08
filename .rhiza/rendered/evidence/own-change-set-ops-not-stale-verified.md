# own-change-set-ops-not-stale-verified

evidence · `01M6XW1191E42JD8Z6QMZ50G7Q` · actor [claude-code-8f1f73eb](../actor/claude-code-8f1f73eb.md)

## result

PASS

## statement

At 1e23bd85, read the diff from 5d4cc959: checkBase in internal/model/apply.go now skips history entries of the op's own change set when looking for stale overlaps, and returns ErrRetired when only own ops (including its own retire) follow the base. Ran go test -trimpath -count=1 -v ./internal/model/...: all pass, including TestOwnOpsOnOldBaseLand (01M6WMCCRC0K shape: update, rename, retire on the landed base; WhatIf has no problems, as-is after landing has none, node retired as pki; with another change set's landed statement write before it, WhatIf reports ErrConflict on csB), TestApplyCases 'a change set's own earlier ops never make its base stale', 'another change set's op between a change set's own ops still conflicts' (ErrConflict), 'a change set writing a node again after retiring it on an older base is skipped' (ErrRetired) and 'a change set retiring a link twice on an older base is skipped' (ErrRetired). go vet and gofmt -l on internal/model are clean.

## Links

### Out

- verifies [own-change-set-ops-not-stale](../capability/rhiza/change-sets/conflicts/own-change-set-ops-not-stale.md)

## History

- `01M6XMTZ6T3CVKWG5013HWA51F` Fix A change set that updates a node and then retires it cannot land: buildRetire (internal/api/write_node.go ~356) bases the node.retire op on the node's landed version (n.Version), ignoring the update this change set has pending on that node. At land the retire then counts as a stale base against its own change set's update. Evidence: change set 01M6WMCCRC0K (ship-change conflict fix, PR #383, green) fails land with 'model: conflict: base 01M68581Q6V7XS1VHP9Q4F011R of 01M68581Q6V7XS1VHP9Q48KQJS is stale: op 01M6WMRS1ZPY of change set 01M6WMCCRC0K since wrote field:acceptance, field:statement, handle'. Op 01M6WMRS1ZPY renames and updates the node, then op 01M6WMSB28 retires it on the old base. The same happens on node 01M6A0NG..CSVTX3B (ops ..3TSKP9 then ..34BK9A). No CLI path recovers: set --rebase refuses a retired node, retire refuses an already-retired node, and carve only moves nodes created in the source. Fix both: (1) retiring a node with pending updates in the same change set drops those updates (as it already does for a pending create) and bases the retire on the landed version; (2) land/fold accept an existing log where a change set's later op's base is superseded only by its own earlier op on the same node, so 01M6WMCCRC0K lands as it is. Add tests for update-then-retire and for the existing log shape. After it ships, retry run 01M6X2EMMHKMQ9W4TDJH3ERXSD (merge main into cs/fix-ship-change-stops-on-a first).
