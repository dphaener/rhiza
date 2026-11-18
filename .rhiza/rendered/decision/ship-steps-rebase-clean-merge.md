# ship-steps-rebase-clean-merge

decision · `01M7MNGSXE269TM6M1G3V0AAVD`

## reasoning

ship rebase and ship check --settled already settle rendered-only conflicts, run the local gate and pre-push check on the merged head and push it, and record the pushed head so the next CI wait follows it; only the workflow routing sends a clean merge to the conductor. Splitting the merge step is needed because a command step has only pass and fail: ship check's red (CI red, unresolved conflict) must hand back while ship merge's main-moved exit must rebase. The merged head is re-checked by the local gate and the pre-push check inside ship rebase and by CI in wait-ci and pr-green before merge runs again.

## statement

The builtin ship-change brings origin/main into the branch in its own steps and hands back to conduct only when that cannot finish cleanly: wait-ci runs rhiza ship check --settled without --no-rebase, so a PR GitHub reports conflicting is merged with main there (a clean merge pushed and waited on, a code conflict left in progress and the wait settled); the old merge step splits into pr-green (rhiza ship check, fail hands back) and merge (rhiza ship merge alone, fail goes to a new rebase step); rebase runs plain rhiza ship rebase, which merges main, renders, runs the local gate and the pre-push check and pushes, then passes back to wait-ci for a CI round on the new head, and hands back on fail: a code conflict left in progress (the conflict fixer's) or the local gate or pre-push check failing on the merged head, kept unpushed (the gate fixer's, since only a fixer can change code). A merge that failed for a gh, fetch or git error also goes to rebase, which merges nothing when main has not moved, so the wait, check and merge retry in the steps. wait-ci, pr-green, merge and rebase each take max_visits 6, as plan.toml's merge-rebase loop does; a hand-back gives them their visits again. The conductor never runs rhiza ship rebase, with or without --no-push; the --no-push and --no-rebase flags stay for other workflows. Signoff is not asked again after a rebase, as in ship-change-classic.

## wording

Derived from the change set's intent and the code (internal/api/ship_rebase.go, shipSettleConflict in internal/api/ship.go, plan.toml's merge and rebase steps), autonomous plan run

## Links

### Out

- supersedes [conduct-main-moved-merge-no-push](conduct-main-moved-merge-no-push.md)
- derived-from [ship-steps-rebase-plan-2026-11-17](../source/ship-steps-rebase-plan-2026-11-17.md)

### In

- decided-by [ship-steps-rebase-clean](../capability/rhiza/workflows/built-in-workflows/ship-steps-rebase-clean.md)
- decided-by [ship-one-ci-rebase-in-steps](../capability/rhiza/orchestration/shipping/ship-one-ci-rebase-in-steps.md)
- decided-by [builtin-conductor-default](../capability/rhiza/workflows/built-in-workflows/builtin-conductor-default.md)
- decided-by [ship-one-ci-round](../capability/rhiza/orchestration/shipping/ship-one-ci-round.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
