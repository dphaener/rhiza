# build-evidence-commit

capability · `01M5YX3E6THKN8CJCSJKQRFXC2` · parent [build-workflow](../build-workflow.md)

## acceptance

A test with staged and unstaged builder edits in the build worktree finds, after a verifier passes, one new commit holding only the Evidence entry's file, and the builder's edits still staged and unstaged as before.

## statement

When a verifier passes a leaf, the engine commits only that leaf's Evidence entry's path, as its own chore(rhiza) commit naming the leaf, never the builder's work in the worktree.

## Links

### Out

- depends-on [build-verifier-checkout](build-verifier-checkout.md)
- decided-by [verifier-writes-evidence](../../../../../decision/verifier-writes-evidence.md)
- decided-by [answer-ab54d544-7](../../../../../decision/answer-ab54d544-7.md)
- derived-from [build-workflow-plan](../../../../../source/build-workflow-plan.md)

### In

- depends-on [build-workflow-file](build-workflow-file.md)
- verifies [build-evidence-commit-verified](../../../../../evidence/build-evidence-commit-verified.md)
- verifies [pr-44-the-build-workflow](../../../../../evidence/pr-44-the-build-workflow.md)

## History

- `01M5WPTY7P14SSHQ1FJ0WA8ZR4` The build workflow: Rhiza's run engine conducts a change set's leaves, a builder session and then a fresh verifier session per leaf in panes of its tab, with only the verifier writing the leaf's Evidence and retries capped on stronger models.
