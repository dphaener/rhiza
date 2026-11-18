# evidence-commit-waits-index-lock

capability · `01M7PSZC9T64QYFHRD9KP08P3A`

## acceptance

Every go test below runs inside systemd-run --user --scope --quiet -p MemoryMax=6G -p TasksMax=300 with -p 2 and RHIZA_*/SHIP_*/REX_* unset, one at a time. (1) A new test TestCommitOnlyIndexLocked in internal/api/run_leaves_test.go, built on evidenceCommitRunner and raceEvidence, creates the run worktree's index.lock (git rev-parse --git-path index.lock) from the commitOnlyHeadRead hook and removes it about 300 ms later; CommitEvidence succeeds, HEAD is the Evidence commit, the hook ran once (the commit was built once), and git status shows the Evidence entries committed with the real index equal to HEAD. Before the fix the same test fails with git's 'index.lock': File exists error. (2) (failure case) A new test TestCommitOnlyIndexLockedThroughout holds index.lock for the whole call: CommitEvidence fails with an error naming index.lock after the bounded wait (well under 10 s), the hook ran once, and the lock file is left as the test made it. (3) go test -trimpath -p 2 ./internal/api -run 'TestCommitOnly|TestEvidenceCommit' passes, TestCommitOnlyHeadMoved and TestCommitOnlyHeadMovedEveryTry unchanged.

## files

internal/api/run_leaves.go, internal/api/run_leaves_test.go

## must_not

Retry update-index on any error but an existing index.lock; rebuild or re-move the commit on that retry; wait unbounded or over about 2 s; change the headless-run-agent-tabs script or any other testscript; use busy loops, CPU stress or parallel stress runs to reproduce

## statement

In internal/api/run_leaves.go, commitOnly's update-index --add of the committed paths on the run worktree's real index, run after a successful update-ref, retries while git fails because index.lock exists (a builder's concurrent git add or commit in the shared run worktree), with short sleeps for a bounded time of about 2 s, then fails with git's error; any other update-index failure still fails at once, and the commit is not built again (HEAD already holds it). Its doc comment says why

## Links

### Out

- derived-from [headless-tabs-index-lock-plan](../source/headless-tabs-index-lock-plan.md)

### In

- verifies [evidence-commit-waits-index-lock-verified](../evidence/evidence-commit-waits-index-lock-verified.md)
- verifies [evidence-commit-waits-index-lock-qa-1](../evidence/evidence-commit-waits-index-lock-qa-1.md)

## History

- `01M7PSG0YN06CTQCB1Z4KZ77WQ` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172, PR #593 at head 4c727d1d1f52c4b94be6856422d2b1376284dd54; reproduce it locally, repeatedly, then fix it or quarantine it
