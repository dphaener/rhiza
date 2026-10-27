# build-evidence-commit-verified

evidence · `01M5ZMT64Y4FXZBXZEA30WE404` · actor [claude-code-c1b3cd7b-verifier](../actor/claude-code-c1b3cd7b-verifier.md)

## result

PASS

## statement

TestEvidenceCommit, on a real git worktree with staged, unstaged and untracked builder edits, finds after CommitEvidence one new commit on the old HEAD holding only the verifier's Evidence entry files (its add and its link) plus the entries of the verifier's actor record (the 23e31f5 fix, so the Actor: trailer resolves), no builder file; git diff --cached, git diff and status of the builder's files are byte-identical before and after; a second call makes no new commit. Testscript build-evidence-commit finds the same end to end (HEAD^ is the builder's commit, MM notes.txt and ?? scratch.txt unchanged). TestLeavesStepEvidenceCommit covers the engine calling it once on a pass. All passed. Note: 'only the Evidence entry's file' is met as the Evidence entries plus the verifier actor's record entries.

## tests

go test -count=1 ./internal/api -run TestEvidenceCommit; ./cmd/rhiza -run 'TestScripts/build-evidence-commit$'; ./internal/run -run TestLeavesStepEvidenceCommit (re-run 2026-10-27 at commit c095271 by claude-code-c1b3cd7b-verifier; go test ./... -count=1 also passed)

## Links

### Out

- verifies [build-evidence-commit](../capability/rhiza/workflows/built-in-workflows/build-workflow/build-evidence-commit.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
