# fix-01m6xmt-plan

source · `01M6XN6HQP02JPRD38Z06CJCY5`

## quotes

```json
[]
```

## raw_sha256

c173074b77a8388b06e65d15da16e3b635c1a8fc2fff132815a15d3f2c3b765a

## session

8d4fbe69-c91e-4026-a866-01a1bb4e1527

## summary

Autonomous plan run, no owner. The planner read internal/api/write_node.go: buildRetire drops a node's pending ops only when the node was created in this change set; for a landed node it writes node.retire with Base n.Version (the landed version), ignoring this change set's pending update or move on that node. At land, model.checkBase (internal/model/apply.go) walks the target's history after the base and reports a conflict against the change set's own earlier update, which is how change set 01M6WMCCRC0K (PR #383) fails to land. The intent prescribes both fixes. One question went to the arbiter, who ruled (assumed, would have escalated) in Decision own-ops-never-stale: checkBase skips history entries of the op's own change set for every node and link op; other change sets' entries after the base still conflict on overlapping aspects, including interleaved ones; the ErrBase and retired-target checks are unchanged. The planner split the work into two leaves: retire drops the change set's pending updates and moves on a landed node and bases the retire on the landed version (model-writes), and an op is never stale against its own change set's ops (conflicts), with a test of the 01M6WMCCRC0K log shape. The change set serves milestone-25, like the other backlog fixes.

## title

Plan of fix 01M6XMT: update-then-retire in one change set lands (autonomous)

## Links

### In

- derived-from [retire-drops-pending-updates](../decision/retire-drops-pending-updates.md)
- derived-from [own-ops-never-stale](../decision/own-ops-never-stale.md)

## History

- `01M6XMTZ6T3CVKWG5013HWA51F` Fix A change set that updates a node and then retires it cannot land: buildRetire (internal/api/write_node.go ~356) bases the node.retire op on the node's landed version (n.Version), ignoring the update this change set has pending on that node. At land the retire then counts as a stale base against its own change set's update. Evidence: change set 01M6WMCCRC0K (ship-change conflict fix, PR #383, green) fails land with 'model: conflict: base 01M68581Q6V7XS1VHP9Q4F011R of 01M68581Q6V7XS1VHP9Q48KQJS is stale: op 01M6WMRS1ZPY of change set 01M6WMCCRC0K since wrote field:acceptance, field:statement, handle'. Op 01M6WMRS1ZPY renames and updates the node, then op 01M6WMSB28 retires it on the old base. The same happens on node 01M6A0NG..CSVTX3B (ops ..3TSKP9 then ..34BK9A). No CLI path recovers: set --rebase refuses a retired node, retire refuses an already-retired node, and carve only moves nodes created in the source. Fix both: (1) retiring a node with pending updates in the same change set drops those updates (as it already does for a pending create) and bases the retire on the landed version; (2) land/fold accept an existing log where a change set's later op's base is superseded only by its own earlier op on the same node, so 01M6WMCCRC0K lands as it is. Add tests for update-then-retire and for the existing log shape. After it ships, retry run 01M6X2EMMHKMQ9W4TDJH3ERXSD (merge main into cs/fix-ship-change-stops-on-a first).
