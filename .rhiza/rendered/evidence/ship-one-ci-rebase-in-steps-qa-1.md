# ship-one-ci-rebase-in-steps-qa-1

evidence · `01M7P15KQZSRBHKDT3H7QNXKEP` · actor [claude-code-90cc7b68](../actor/claude-code-90cc7b68.md)

## result

PASS

## statement

QA: TestScripts/ship-change-one-ci passes; its assertions cover merge fail -> rebase pass -> wait-ci/pr-green/merge visit 2 with conduct once, no handback, two pushes and one pr merge, and the conflict case rebase fail -> conduct visit 2 -> gate, signoff, push, merge; rhiza check in a clean detached checkout at HEAD exits 0 and rhiza render leaves the tree clean (fixer's rendered index order holds)

## Links

### Out

- verifies [ship-one-ci-rebase-in-steps](../capability/rhiza/orchestration/shipping/ship-one-ci-rebase-in-steps.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
