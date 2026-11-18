# ship-one-ci-round

capability · `01M79KPE2V3AZWKX2HR1X1PZ1M` · parent [shipping](../shipping.md)

## acceptance

go test -trimpath ./cmd/rhiza -run TestScripts/ship-change-one-ci shows, with a fake gh and a fake conductor that lands and renders on the branch, exactly one CI wait and one push on the happy path and the change set landed on main after the merge with no ship land step; (failure case) when main moves before the merge, rhiza ship merge fails and the run shows two CI waits, two pushes and no hand-back to conduct.

## files

cmd/rhiza/testdata/script/ship-change-one-ci.txtar

## must_not

Start any process the script's processScripts entry in cmd/rhiza/main_test.go does not already allow; touch the network, Rex, agents or the real gh; change ship-change.toml or internal/api

## statement

On the conductor default the change set lands on the branch before the first push, so a run whose PR is green and whose main has not moved waits on CI once; a second CI round happens only when the merge finds main has moved, and the ship steps rebase and wait for it without a conductor visit.

## Links

### Out

- depends-on [builtin-conductor-default](../../workflows/built-in-workflows/builtin-conductor-default.md)
- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)
- decided-by [ship-steps-rebase-clean-merge](../../../../decision/ship-steps-rebase-clean-merge.md)

### In

- verifies [ship-one-ci-round-verified](../../../../evidence/ship-one-ci-round-verified.md)
- verifies [ship-one-ci-round-qa-2](../../../../evidence/ship-one-ci-round-qa-2.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
