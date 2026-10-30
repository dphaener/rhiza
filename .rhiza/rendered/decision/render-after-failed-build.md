# render-after-failed-build

decision · `01M66BMTMSV59BC3THADY6KDDF`

## statement

A failed build still renders and commits its branch before it ends, so a failed branch also passes rhiza check and a retry starts from a clean branch.

## Links

### Out

- answers [build-render-after-failed-build](../question/build-render-after-failed-build.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
