# evidence-commit-waits-index-lock-verified

evidence · `01M7PVW1S57XJYD4BYG46XRYHM` · actor [claude-code-f46866aa](../actor/claude-code-f46866aa.md)

## result

PASS

## statement

At 441a2631c, in a systemd-run scope (MemoryMax=6G TasksMax=300, -p 2, RHIZA_/SHIP_/REX_ unset): go test -trimpath -p 2 ./internal/api -run 'TestCommitOnly|TestEvidenceCommit' -count=1 -v exit 0 (TestCommitOnlyIndexLocked 0.54s, TestCommitOnlyIndexLockedThroughout 2.18s, HeadMoved/HeadMovedEveryTry/EvidenceCommit pass); go vet ./internal/api exit 0. Pre-fix (run_leaves.go from HEAD~1 in a scratch worktree, plus only the indexLockWait const so tests compile): TestCommitOnlyIndexLocked fails with git's "index.lock': File exists" and TestCommitOnlyIndexLockedThroughout fails after 59ms with no bounded wait. Diff read: retry only on "index.lock': File exists", 50ms poll bounded by 2s, update-ref/commit not repeated (hook asserted once), only run_leaves.go and run_leaves_test.go changed besides the builder's actor record; HeadMoved tests untouched; no testscript changed; no TODOs.

## Links

### Out

- verifies [evidence-commit-waits-index-lock](../capability/evidence-commit-waits-index-lock.md)

## History

- `01M7PSG0YN06CTQCB1Z4KZ77WQ` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172, PR #593 at head 4c727d1d1f52c4b94be6856422d2b1376284dd54; reproduce it locally, repeatedly, then fix it or quarantine it
