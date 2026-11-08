# plan-fix-watchdog-unmerged-pr-at

source · `01M6XHVA4Q0021AZMQD4E0878G`

## quotes

```json
[]
```

## raw_sha256

e14ad724409bdc5a2d840568e683f3d1f80e7dabeb25485b22a40c4036cef8cd

## session

cf4d0c69-229d-4917-b1e9-01a1bb16ffc1

## summary

Autonomous plan run, no owner. Diagnosis: run 01M6X2EMMHKMQ9W4TDJH3ERXSD failed at step land because change set 01M6WMCCRC0KSCWEV17VY8KVQP (PR 383) updated two nodes and later retired them; rhiza retire based each retire on the node's landed version (Node.Version skips the change set's own ops), so model.checkBase refused the land as a conflict (exit 3) against the change set's own earlier update. rhiza check only warned (rule in-flight), and ship's head-build error dropped the child's error message, so the run log said only 'exit status 3'. The planner asked how to fix it (model tolerance of own ops, writer bases, or both); the arbiter decided retire-supersedes-own-pending-edits: rhiza retire drops the change set's pending update/move ops on the node and bases the retire on the landed version, as it already does for a pending create; then repair PR 383 by re-running rhiza retire in its change set with the fixed build; add a model/api test of update-then-retire landing; making the in-flight warning an error is out of scope. The planner kept the head-build error message as a second leaf within the intent.

## title

Plan: fix the land refused for an update-then-retire change set (autonomous)

## Links

### In

- derived-from [retire-supersedes-pending-edits](../capability/rhiza/change-sets/conflicts/retire-supersedes-pending-edits.md)
- derived-from [ship-head-build-error-detail](../capability/rhiza/orchestration/shipping/ship-head-build-error-detail.md)

## History

- `01M6XHCX4FRER8CR0MG5D4D76D` Fix watchdog unmerged PR at step land: run 01M6X2EMMHKMQ9W4TDJH3ERXSD; log: PR green at 6d6c4d1 and unmerged after step land failed: step "land" ended with outcome "fail" | exit 1 | rhiza: ship: land with the PR head's build: /home/rhiza/code/rhiza/.git/rhiza/ship/01M6WMCCRC0KSCWEV17VY8KVQP/rhiza land --change 01M6WMCCRC0KSCWEV17VY8KVQP --json: exit status 3
