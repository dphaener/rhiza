# review-renders-branch

capability · `01M6506T856R5QRMHAYGP8QGJ6` · parent [workflows](../workflows.md)

## acceptance

A review run whose fixers or checks leave .rhiza/rendered stale ends with a commit 'Render after the review' after their last record, rhiza check passes on its branch, and no .rhiza/log entry is left uncommitted

## statement

A plain review run ends with its branch rendered and committed after its last fixer, verifier or check has written a record

## Links

### Out

- decided-by [render-at-workflow-end](../../../decision/render-at-workflow-end.md)

### In

- verifies [review-renders-branch-verified](../../../evidence/review-renders-branch-verified.md)
- verifies [review-renders-branch-qa-1](../../../evidence/review-renders-branch-qa-1.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
