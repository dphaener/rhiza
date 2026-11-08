# own-ops-never-stale

decision · `01M6XN5972NYED9NYJFMQNEHV3`

## reasoning

The intent (change set 01M6XMT, part 2) asks that land/fold accept a log where a later op's base was superseded only by its own change set's earlier op, so 01M6WMCCRC0K lands as written. model.checkBase (internal/model/apply.go:284) is the one place every node and link op's base is checked, and history entries already carry their ChangeSet, so a same-change-set skip there is a small, uniform change. Limiting it to node.retire would leave the same trap for link retires and for update-after-update written on an old base by a parallel session in one change set, where the ops are one unit of intent and cannot conflict with themselves in any meaningful sense. Conflict detection between change sets, which is what the stale-base check protects, is untouched.

## statement

Apply the rule to every op type on nodes and links, checked in model.checkBase: when walking the target's history after an op's base, skip entries written by the op's own change set, so an op is never stale against its own change set's ops. Entries from other change sets after the base still conflict on overlapping aspects exactly as today, including ones interleaved between the change set's own ops. The other checks are unchanged: a base that is not an operation on the target still fails, and writing a retired target still fails.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-8d4fbe69-1](../question/question-8d4fbe69-1.md)
- derived-from [fix-01m6xmt-plan](../source/fix-01m6xmt-plan.md)

### In

- decided-by [own-change-set-ops-not-stale](../capability/rhiza/change-sets/conflicts/own-change-set-ops-not-stale.md)

## History

- `01M6XMTZ6T3CVKWG5013HWA51F` Fix A change set that updates a node and then retires it cannot land: buildRetire (internal/api/write_node.go ~356) bases the node.retire op on the node's landed version (n.Version), ignoring the update this change set has pending on that node. At land the retire then counts as a stale base against its own change set's update. Evidence: change set 01M6WMCCRC0K (ship-change conflict fix, PR #383, green) fails land with 'model: conflict: base 01M68581Q6V7XS1VHP9Q4F011R of 01M68581Q6V7XS1VHP9Q48KQJS is stale: op 01M6WMRS1ZPY of change set 01M6WMCCRC0K since wrote field:acceptance, field:statement, handle'. Op 01M6WMRS1ZPY renames and updates the node, then op 01M6WMSB28 retires it on the old base. The same happens on node 01M6A0NG..CSVTX3B (ops ..3TSKP9 then ..34BK9A). No CLI path recovers: set --rebase refuses a retired node, retire refuses an already-retired node, and carve only moves nodes created in the source. Fix both: (1) retiring a node with pending updates in the same change set drops those updates (as it already does for a pending create) and bases the retire on the landed version; (2) land/fold accept an existing log where a change set's later op's base is superseded only by its own earlier op on the same node, so 01M6WMCCRC0K lands as it is. Add tests for update-then-retire and for the existing log shape. After it ships, retry run 01M6X2EMMHKMQ9W4TDJH3ERXSD (merge main into cs/fix-ship-change-stops-on-a first).
