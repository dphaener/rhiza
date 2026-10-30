# build-renders-branch

capability · `01M6506RT8ZF7BTAQNKA5X7XV4` · parent [workflows](../workflows.md)

## acceptance

A build run whose sessions leave .rhiza/rendered stale ends with a commit 'Render after the build' naming the change set, rhiza check passes on its branch, and no .rhiza/log entry is left uncommitted

## statement

A plain build run ends with its branch rendered and committed, so rhiza check passes on it without a hand-made render commit

## Links

### Out

- decided-by [render-at-workflow-end](../../../decision/render-at-workflow-end.md)

### In

- blocks [build-render-after-failed-build](../../../question/build-render-after-failed-build.md)
- verifies [build-renders-branch-verified](../../../evidence/build-renders-branch-verified.md)
- verifies [build-renders-branch-qa-1](../../../evidence/build-renders-branch-qa-1.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
