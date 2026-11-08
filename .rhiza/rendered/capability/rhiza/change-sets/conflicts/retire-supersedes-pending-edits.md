# retire-supersedes-pending-edits

capability · `01M6XHVSGJT8C0VDSS1ZQAWBM9` · parent [conflicts](../conflicts.md)

## acceptance

Tests in internal/api (and internal/model where the fold is checked): a change set sets a node's statement and handle, then retires it in a later entry; the change set holds only the retire, based on the landed version, and rhiza land succeeds with no conflict; the same with a node.move before the retire; a change set holding an update and a retire on the stale landed base (written as before this fix) is repaired by rhiza retire on that node, after which it lands; a testscript retire-after-update.txtar shows update, retire and land in one change set exit 0 with fold-ok; run with go test -trimpath ./internal/api/... ./internal/model/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/retire-after-update'

## statement

rhiza retire on a node the current change set has pending node.update or node.move ops on drops those ops from the change set and writes the retire on the node's landed version, as it already drops a pending create, so a change set that edits a node and later retires it holds one retire on a valid base and lands; run again on a node the change set already retired on a stale base (its own earlier edits after the landed version), it drops those edits and rewrites the retire onto the landed version, which repairs a change set written before this fix

## Links

### Out

- decided-by [retire-supersedes-own-pending-edits](../../../../decision/retire-supersedes-own-pending-edits.md)
- derived-from [plan-fix-watchdog-unmerged-pr-at](../../../../source/plan-fix-watchdog-unmerged-pr-at.md)

### In

- verifies [retire-supersedes-pending-edits-verified](../../../../evidence/retire-supersedes-pending-edits-verified.md)
- verifies [retire-supersedes-pending-edits-qa-1](../../../../evidence/retire-supersedes-pending-edits-qa-1.md)
- verifies [pr-389-fix-watchdog-unmerged-pr-at](../../../../evidence/pr-389-fix-watchdog-unmerged-pr-at.md)

## History

- `01M6XHCX4FRER8CR0MG5D4D76D` Fix watchdog unmerged PR at step land: run 01M6X2EMMHKMQ9W4TDJH3ERXSD; log: PR green at 6d6c4d1 and unmerged after step land failed: step "land" ended with outcome "fail" | exit 1 | rhiza: ship: land with the PR head's build: /home/rhiza/code/rhiza/.git/rhiza/ship/01M6WMCCRC0KSCWEV17VY8KVQP/rhiza land --change 01M6WMCCRC0KSCWEV17VY8KVQP --json: exit status 3
