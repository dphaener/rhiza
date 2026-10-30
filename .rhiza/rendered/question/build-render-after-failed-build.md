# build-render-after-failed-build

question · `01M6506WEYXT6D39X1J6ZPGZ55`

## statement

When a build fails, should the run still render and commit its branch, so the failed branch also passes rhiza check?

## Links

### Out

- blocks [build-renders-branch](../capability/rhiza/workflows/build-renders-branch.md)

### In

- answers [render-after-failed-build](../decision/render-after-failed-build.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
