# ship-steps-rebase-plan-2026-11-17

source · `01M7MNEHW3HGWMD6S0K4X9Z15P`

## quotes

```json
[]
```

## raw_sha256

47ed8493b9d8a15c511ab6c6c2abd21a53d97449cc4c7b69474cc1ed03962646

## session

04560a76-2703-48d7-8ab3-01a1e953bfe4

## summary

Autonomous plan run for the bug that ship-change hands a whole change set back to the conductor when its merge finds origin/main moved or GitHub reports the PR conflicting. No owner took part. Research: internal/workflow/builtin/ship-change.toml routes the merge step (rhiza ship check --no-rebase && rhiza ship merge) to handback:conduct on every fail, and wait-ci passes --no-rebase so a conflicting PR is never merged with main there; the conductor then runs rhiza ship rebase --no-push and goes through gate, signoff, push and CI again. rhiza ship rebase (internal/api/ship_rebase.go) already merges origin/main, settles rendered-only conflicts, runs the local gate and pre-push check and pushes, leaving a code conflict in progress (exit 1) and a merged head the gate failed on unpushed (exit 1); rhiza ship check --settled without --no-rebase already merges main on a conflicting PR the same way (shipSettleConflict) and keeps waiting on the pushed head. plan.toml's merge step already fails to a rebase step that loops back through a CI wait (max_visits 6). Plan: wait-ci drops --no-rebase; the merge step is split into pr-green (rhiza ship check, fail hands back) and merge (rhiza ship merge, fail goes to a new rebase step running rhiza ship rebase that passes back to wait-ci and hands back on fail: a code conflict left in progress, or the local gate or pre-push check failing on the merged head); the conduct prompt, the conduct-handback and conduct-roles topics, the workflows topic, docs/primitives.md and the ship-change-one-ci script follow. rhiza ship rebase --no-push stays as a flag. Not breaking, not severe.

## title

Plan: ship steps rebase a clean merge of main (autonomous plan run)

## Links

### In

- derived-from [ship-steps-rebase-clean-merge](../decision/ship-steps-rebase-clean-merge.md)
- derived-from [conduct-handback-routes-ship-rebase](../decision/conduct-handback-routes-ship-rebase.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
