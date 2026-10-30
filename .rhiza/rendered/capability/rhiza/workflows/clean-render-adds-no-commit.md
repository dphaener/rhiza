# clean-render-adds-no-commit

capability · `01M6506VDTVZ6SYN4WFQE44VFS` · parent [workflows](../workflows.md)

## acceptance

A build or review run on an already rendered branch ends with the same HEAD it started with

## statement

A plain build or review run whose branch is already rendered, with no uncommitted .rhiza/log entries, adds no commit

## Links

### Out

- decided-by [render-at-workflow-end](../../../decision/render-at-workflow-end.md)

### In

- verifies [clean-render-adds-no-commit-verified](../../../evidence/clean-render-adds-no-commit-verified.md)
- verifies [clean-render-adds-no-commit-qa-1](../../../evidence/clean-render-adds-no-commit-qa-1.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
