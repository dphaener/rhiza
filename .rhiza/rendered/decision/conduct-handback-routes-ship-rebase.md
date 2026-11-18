# conduct-handback-routes-ship-rebase

decision · `01M7MNGT1EYZJED01XCCWWXEET`

## statement

A conductor handed back by ship-change reads handback.md and spawns, never fixes: gate's rhiza run build-check a verifier for the leaf lacking independent Evidence (a builder retry if it fails); gate's rhiza run review-check a critic for the check with an open blocker or major, then a fixer for confirmed findings and the checks again; gate's rhiza check a fixer with the output as its findings; gate's rhiza ship gate and the push step's rhiza ship open-pr (its pre-push check) the gate-fix role; signoff changes a fixer with the owner's requested changes, then the checks; wait-ci's timeout no role (read rhiza ship check --json, change nothing, finish); pr-green's settled red CI the ci-fix role, its unresolved merge of main (left by the CI wait on a conflicting PR) the conflict-fix role, and a merged head the local gate or pre-push check failed on the gate-fix role; rebase's code conflict left in progress the conflict-fix role, its local gate or pre-push check failing on the merged head the gate-fix role, and any other rebase failure (a gh, fetch or git error) no role: change nothing and finish. The merge step no longer hands back. The conductor never runs rhiza ship rebase itself. Every route lands and renders again, then rhiza session done; handback_cap and the steps after the merge are as before.

## wording

Hand-back routes as ship-change's new pr-green and rebase steps and the fixer role topics give them (derived from the code, autonomous plan run)

## Links

### Out

- supersedes [conduct-handback-routes](conduct-handback-routes.md)
- derived-from [ship-steps-rebase-plan-2026-11-17](../source/ship-steps-rebase-plan-2026-11-17.md)

### In

- decided-by [conduct-handback-ship-rebase](../capability/rhiza/guide/agent-guide/conduct-handback-ship-rebase.md)
- decided-by [conduct-handback-topic](../capability/rhiza/guide/agent-guide/conduct-handback-topic.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
