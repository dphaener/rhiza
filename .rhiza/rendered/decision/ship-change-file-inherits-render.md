# ship-change-file-inherits-render

decision · `01M66E7QQXGY17RMND0HHR6NYV`

## reasoning

The model check (review cycle 2, finding f1) confirmed that 83c71b7 moved ship-change's render step into the build workflow and that the review workflow now renders too, so the statement decided by ship-change-file-names-render no longer matched the code. The code is what this change set intends, so the statement changes.

## statement

Leaf ship-change-file's statement no longer lists a render command step of ship-change's own between the build and review workflow steps: ship-change runs build/build, build/render, build/check, review/code, review/review, review/render and then render-pr, so its render before review is inherited from build (see build-renders-branch) and review/render runs before render-pr.

## wording

decided by the orchestrator acting as arbiter, on the owner's word that the orchestrator signs off on this change set

## Links

### Out

- supersedes [ship-change-file-names-render](ship-change-file-names-render.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
