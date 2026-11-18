# ship-steps-rebase-clean-qa-2

evidence · `01M7P1AFJRPREY7801JC9HF9ET` · actor [claude-code-1bba01de](../actor/claude-code-1bba01de.md)

## result

PASS

## statement

QA with rhiza built from 3d2eb196f to a scratch path, in a throwaway repo, capped scope, RHIZA_/SHIP_/REX_ unset. Happy: workflow describe ship-change lists 12 steps conduct..deploy-check; merge on.fail->rebase, on.pass->main-ci; rebase on.fail->handback:conduct, on.pass->wait-ci; pr-green fail->handback:conduct; main-ci/close-out/ship/deploy-check no on.fail. workflow new --from ship-change copy shows wait-ci until 'rhiza ship check --settled' (no --no-rebase), merge run 'rhiza ship merge', rebase run 'rhiza ship rebase' stall 70m, max_visits 6 on wait-ci/pr-green/merge/rebase, handback_cap 3, conduct prompt without --no-push. In this repo ship-rhiza and bugfix describe/check valid with the same 12 included steps. describe --json repeated is byte-identical and equals api call workflow.describe. Failure: fork with [step.rebase] removed -> workflow check exit 1 'step merge: on.fail names "rebase", which is not a step'; workflow new onto an existing name refused exit 3; ship rebase with no change set exit 2 with next hint.

## Links

### Out

- verifies [ship-steps-rebase-clean](../capability/rhiza/workflows/built-in-workflows/ship-steps-rebase-clean.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
