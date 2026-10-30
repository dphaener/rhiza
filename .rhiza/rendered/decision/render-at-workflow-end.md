# render-at-workflow-end

decision · `01M6506Q1V4DQEXQY42A3XNK0W`

## statement

Plain build and review render their branch as the last step of their own workflow, as ship-change does, not in the run engine

## wording

the agent's option, assumed in an autonomous plan (would have escalated)

## Links

### Out

- derived-from [build-review-render-plan](../source/build-review-render-plan.md)

### In

- decided-by [build-renders-branch](../capability/rhiza/workflows/build-renders-branch.md)
- decided-by [review-renders-branch](../capability/rhiza/workflows/review-renders-branch.md)
- decided-by [clean-render-adds-no-commit](../capability/rhiza/workflows/clean-render-adds-no-commit.md)

## History

- `01M61GW7K1PDG7HMFXMN3YNA7B` Plain build and review runs leave their branch rendered: each renders and commits .rhiza/rendered and the uncommitted .rhiza/log entries before it ends, so rhiza check passes on the branch without a hand-made render commit.
