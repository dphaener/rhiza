# ship-steps-rebase-clean-verified

evidence · `01M7NYTH2A4A9AWCPCZE6Q35X3` · actor [claude-code-5153b929](../actor/claude-code-5153b929.md)

## result

PASS

## statement

In a detached checkout at 48b07ac34 (capped scope, RHIZA_/SHIP_/REX_ unset): go test -trimpath ./internal/workflow -run 'TestBuiltinShipChange|TestPlainRedLoopsShipChangeHasNone|TestMergeWithoutGate' exit 0; go test ./internal/workflow/... exit 0; go test ./cmd/rhiza -run TestScripts/workflow-new exit 0. Failure cases: TestBuiltinShipChange exits 1 with ship-change.toml mutated to merge on.fail handback:conduct, wait-ci until with --no-rebase, rebase on.fail wait-ci, rebase stall removed, and rebase stall 40m. Diff read: the test asserts the 12 steps in order, wait-ci until exactly rhiza ship check --settled with timeout to handback:conduct, pr-green running rhiza ship check (pass merge, fail handback), merge running exactly rhiza ship merge (only pass main-ci, fail rebase), rebase running rhiza ship rebase with stall 70m >= gate's (pass wait-ci, fail handback), max_visits 6 on the four steps, handback_cap 3, no fail route or hand-back after the merge, conduct prompt naming the conflict fixer without --no-push. No change to internal/api, ship-change-classic.toml, plan.toml, .rhiza/workflows or docs; no TODOs. rhiza ship generate and rhiza render (scratch build) exit 0 with a clean tree. ship-change-one-ci.txtar still asserts the old shape; that is leaf ship-one-ci-rebase-in-steps's file.

## Links

### Out

- verifies [ship-steps-rebase-clean](../capability/rhiza/workflows/built-in-workflows/ship-steps-rebase-clean.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
