# ship-steps-rebase-clean-qa-1

evidence · `01M7P15KBW4V9MVKE86Z8YTP3F` · actor [claude-code-90cc7b68](../actor/claude-code-90cc7b68.md)

## result

PASS

## statement

QA: built ./cmd/rhiza to a scratch path; rhiza workflow describe ship-change lists the 12 steps with merge fail -> rebase, rebase pass -> wait-ci / fail -> handback:conduct, pr-green fail -> handback:conduct, wait-ci timeout -> handback:conduct, no on.fail on main-ci/close-out/ship/deploy-check; toml has wait-ci until 'rhiza ship check --settled' (no --no-rebase), rebase stall 70m, max_visits 6, handback_cap 3, conduct prompt without --no-push; go test ./internal/workflow -run 'TestBuiltinShipChange|TestPlainRedLoopsShipChangeHasNone|TestMergeWithoutGate' and TestScripts/workflow-new pass

## Links

### Out

- verifies [ship-steps-rebase-clean](../capability/rhiza/workflows/built-in-workflows/ship-steps-rebase-clean.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
