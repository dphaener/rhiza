# retire-supersedes-own-pending-edits

decision · `01M6XHSWN256M1Y406SGD500RB`

## reasoning

internal/api/write_node.go:388 bases a node retire on n.Version, the landed version, while model.checkBase (internal/model/apply.go:284) then finds the change set's own earlier update after that base and reports a conflict. The retire write already handles the analogous pending-create case by dropping the change set's own ops (e.drop), so extending that to pending updates/moves is the smallest fix consistent with existing code: a retire makes earlier edits of the same node in the same change set meaningless, and dropping them avoids basing on an unlanded version, which would break the landed-base invariant other writes rely on. Links are out of scope: internal/op/op.go defines only link.create and link.retire, and buildUnlink (internal/api/write_link.go) already drops a pending link.create, so no link retire can hit this conflict. Re-running retire on PR 383 uses the same supported drop mechanism, so no log entry is edited by hand. Making check fail would catch this in CI but doesn't fix the cause and widens the change set past the watchdog fix.

## statement

Fix it in the retire write: when the current change set already has pending update or move ops on the node being retired, `rhiza retire` drops those ops from the change set and bases the retire on the landed version, the same way it already drops a pending create. The retired node's own pending edits are moot, so the change set then holds one retire on a valid base and lands. Links need no change: the only link ops are link.create and link.retire, and unlink already drops a pending link.create, so a link retire cannot conflict with the change set's own earlier edit. Then repair PR 383 by re-running `rhiza retire` on that node in its change set with the fixed build, so its stale ops are dropped and it lands. Add a model/api test for update-then-retire in one change set landing cleanly. Turning rhiza check's in-flight warning into an error is out of scope for this change set.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-cf4d0c69-1](../question/question-cf4d0c69-1.md)

### In

- decided-by [retire-supersedes-pending-edits](../capability/rhiza/change-sets/conflicts/retire-supersedes-pending-edits.md)

## History

- `01M6XHCX4FRER8CR0MG5D4D76D` Fix watchdog unmerged PR at step land: run 01M6X2EMMHKMQ9W4TDJH3ERXSD; log: PR green at 6d6c4d1 and unmerged after step land failed: step "land" ended with outcome "fail" | exit 1 | rhiza: ship: land with the PR head's build: /home/rhiza/code/rhiza/.git/rhiza/ship/01M6WMCCRC0KSCWEV17VY8KVQP/rhiza land --change 01M6WMCCRC0KSCWEV17VY8KVQP --json: exit status 3
