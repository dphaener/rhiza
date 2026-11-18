# headless-tabs-builder-waits-index-lock

capability · `01M7PSQZS492GXWNYNX3EEK1KR`

## acceptance

Every go test below runs inside systemd-run --user --scope --quiet -p MemoryMax=6G -p TasksMax=300 with -p 2 and RHIZA_*/SHIP_*/REX_* unset, one at a time. (1) go test -trimpath -p 2 -count=30 ./cmd/rhiza -run 'TestScripts/^headless-run-agent-tabs$' passes; before the fix the same command failed about once in 15 to 20 runs with 'attempt 1 builder ... fail (exit 1)' and 'Unable to create ... index.lock: File exists' in $WORK/panes.out. (2) (failure case) With a temporary scratch edit, not committed, that makes the fake rex create $WORK/repo/.git/worktrees/repo.built/index.lock before it starts a builder's pane and remove it about 0.3 s later in the background, the script fails once with the old fake builder (the builder exits 1 on the lock) and passes once with the new one; the commit's Checks: block records both runs. (3) Every existing assertion of the script still holds unchanged.

## files

cmd/rhiza/testdata/script/headless-run-agent-tabs.txtar

## must_not

Change any Go code (internal/ or cmd/rhiza/*.go); drop, loosen or reorder any assertion of the script; skip or quarantine the script; edit other scripts' fake builders; use busy loops, CPU stress or parallel stress runs to reproduce

## statement

In cmd/rhiza/testdata/script/headless-run-agent-tabs.txtar, the fake bin/claude builder retries its git add and git commit of the leaf file (a bounded loop of short sleeps, about 5 s in all, then exit 1) while either fails, so a sibling's or the runner's hold on the shared run worktree's index.lock no longer fails a builder attempt, and the script's header comment says why

## Links

### Out

- decided-by [headless-tabs-fix-fixture-lock](../decision/headless-tabs-fix-fixture-lock.md)

### In

- verifies [headless-tabs-builder-waits-index-lock-verified](../evidence/headless-tabs-builder-waits-index-lock-verified.md)
- verifies [headless-tabs-builder-waits-index-lock-qa-1](../evidence/headless-tabs-builder-waits-index-lock-qa-1.md)

## History

- `01M7PSG0YN06CTQCB1Z4KZ77WQ` Fix flaky test TestScripts/headless-run-agent-tabs: failed then passed on a rerun in PR #513 at head 7a96ac202aae8f55b2f0adfcd46d0e0c70574172, PR #593 at head 4c727d1d1f52c4b94be6856422d2b1376284dd54; reproduce it locally, repeatedly, then fix it or quarantine it
