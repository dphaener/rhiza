# build-renders-branch-verified

evidence · `01M66CXMTE1C3D5A562V143MKS` · actor [claude-code-2b55ed4c](../actor/claude-code-2b55ed4c.md)

## result

PASS

## statement

At 83c71b7: make check passed (an earlier run's only failure was the timing-sensitive tui/modelapp read budget, which passed on rerun). Read the diff: build.toml gains a render step (build -> render -> check) that runs rhiza render, stages .rhiza/rendered and .rhiza/log, and commits 'Render after the build' with a Change: trailer when anything changed; ship-change drops its own render step and inherits build/render. Ran build-workflow.txtar in a scratch worktree with added assertions. After the built-in build run, the build worktree's HEAD is 'Render after the build' / Change: <change id>, committing 6 .rhiza/log entries and the rendered actor pages and index. git status --porcelain -- .rhiza is empty, rhiza check prints ok on the branch, and fold-ok passes.

## Links

### Out

- verifies [build-renders-branch](../capability/rhiza/workflows/build-renders-branch.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
