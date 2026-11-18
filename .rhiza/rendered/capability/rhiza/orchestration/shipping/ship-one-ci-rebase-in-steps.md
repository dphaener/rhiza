# ship-one-ci-rebase-in-steps

capability · `01M7MNGTG0TSZM85WSJ79R55P3` · parent [shipping](../shipping.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/ship-change-one-ci passes: the happy path keeps one CI wait, one push and one merge; when main moves between the CI check and the merge, rhiza run status shows merge (visit 1) fail then rebase, rebase (visit 1) pass then wait-ci, wait-ci (visit 2), pr-green (visit 2) and merge (visit 2) pass then main-ci, conduct visited once and no handback in the status or the fake conductor's log, two pushes of the PR branch, one pr merge and the change set landed; (failure case) when the commit that moves main changes a file the branch also changes, rebase fails and hands back to conduct with the merge of main in progress (rhiza run status shows rebase fail and conduct (visit 2)), the fake conductor resolves and commits it, and the run goes through gate, signoff and push to a merge.

## files

cmd/rhiza/testdata/script/ship-change-one-ci.txtar

## must_not

Start any process the script's processScripts entry in cmd/rhiza/main_test.go does not already allow; touch the network, Rex, agents or the real gh; change ship-change.toml or internal/api

## statement

The ship-change-one-ci script shows a main-moved merge rebased in the ship steps with no hand-back, and a code conflict with main handed back to conduct with the merge left in progress.

## Links

### Out

- depends-on [ship-steps-rebase-clean](../../workflows/built-in-workflows/ship-steps-rebase-clean.md)
- decided-by [ship-steps-rebase-clean-merge](../../../../decision/ship-steps-rebase-clean-merge.md)

### In

- verifies [ship-one-ci-rebase-in-steps-verified](../../../../evidence/ship-one-ci-rebase-in-steps-verified.md)
- verifies [ship-one-ci-rebase-in-steps-qa-1](../../../../evidence/ship-one-ci-rebase-in-steps-qa-1.md)
- verifies [ship-one-ci-rebase-in-steps-qa-2](../../../../evidence/ship-one-ci-rebase-in-steps-qa-2.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
