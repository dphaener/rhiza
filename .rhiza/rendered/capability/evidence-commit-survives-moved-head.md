# evidence-commit-survives-moved-head

capability · `01M7FT14Q1YV5YWP0RN8R7A8XA`

## acceptance

A new internal/api test (go test -trimpath -p 2 -count=1 ./internal/api -run '^TestCommitOnlyHeadMoved') injects a commit on HEAD through the hook between commitOnly's HEAD read and its update-ref: CommitEvidence succeeds, the Evidence commit's parent is the injected commit, both are on the branch, the injected commit's files are intact and the real index shows the Evidence paths committed; failure case: a hook that moves HEAD on every try fails after 5 tries with update-ref's error. go test -trimpath -p 2 ./internal/api -run 'CommitEvidence|CommitOnly|RunLeaves' passes, and go test -trimpath -p 2 -count=3 ./cmd/rhiza -run 'TestScripts/^headless-run-agent-tabs$' passes all three runs. Touches internal/api/run_leaves.go and its test.

## statement

A verifier's Evidence commit (commitOnly) survives another commit landing on the run worktree's HEAD while it builds: it rebuilds on the new HEAD and moves HEAD there, keeping both commits, instead of failing the verifier step.

## Links

### Out

- decided-by [evidence-commit-retries-moved-head](../decision/evidence-commit-retries-moved-head.md)

### In

- verifies [evidence-commit-survives-moved-head-verified](../evidence/evidence-commit-survives-moved-head-verified.md)
- verifies [evidence-commit-survives-moved-head-qa-1](../evidence/evidence-commit-survives-moved-head-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
