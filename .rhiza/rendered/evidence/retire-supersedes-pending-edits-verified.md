# retire-supersedes-pending-edits-verified

evidence · `01M6XSJ4KS5TFVVF3K99CGPS53` · actor [claude-code-d682c09a](../actor/claude-code-d682c09a.md)

## result

PASS

## statement

Verifier at b7cf6e6c. Read the diff: buildRetire drops the change set's own pending update/move ops and writes the retire on the landed version; staleRetire now counts a retire as stale when the change set also holds its own update or move on the node, which fixes attempt 1's finding f1. Tests: TestRetire/SupersedesPendingEdits (handle+move+retire leaves one retire on the landed base), SupersedesPendingEditsLands (lands, folds clean), and RepairsStaleRetire / RepairsStaleRetireAfterMove / RepairsRetireOnUpdate build the old shape (update [+move] and retire all on the landed base), repair it with retire, and land it (covers f2). go test -trimpath ./internal/api -run 'TestRetire|TestStackedWrites' passed; go test -trimpath ./internal/api -run 'Retire|Write|Land|Fold|Stacked|Set|Mv|Move' passed every test (afterwards the package's harness flagged only the known .rhiza-root/usage.live.json leak from the enclosing ~/code root, unrelated to this leaf); go test -trimpath ./internal/model/... ok; go test -trimpath ./cmd/rhiza -run TestScripts/retire-after-update ok. The whole internal/api package needs longer than go test's default 10m (the Makefile uses 20m), so I scoped the run. By hand, in scratch repos with RHIZA_* unset: the pre-fix binary (4b2082fa) wrote set --handle (and in a second repo, mv) then retire, all on the landed base, and its land exited 3; then the new binary's retire <node-id> exited 0, left one node.retire on the landed base, and land exited 0, in both repos.

## Links

### Out

- verifies [retire-supersedes-pending-edits](../capability/rhiza/change-sets/conflicts/retire-supersedes-pending-edits.md)

## History

- `01M6XHCX4FRER8CR0MG5D4D76D` Fix watchdog unmerged PR at step land: run 01M6X2EMMHKMQ9W4TDJH3ERXSD; log: PR green at 6d6c4d1 and unmerged after step land failed: step "land" ended with outcome "fail" | exit 1 | rhiza: ship: land with the PR head's build: /home/rhiza/code/rhiza/.git/rhiza/ship/01M6WMCCRC0KSCWEV17VY8KVQP/rhiza land --change 01M6WMCCRC0KSCWEV17VY8KVQP --json: exit status 3
