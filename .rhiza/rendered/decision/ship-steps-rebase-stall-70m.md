# ship-steps-rebase-stall-70m

decision · `01M7NXZQ6T5HC52MJWNX3M75FR`

## reasoning

Confirmed the finding against the code: watchdog_policy.go sets DefaultCommandLimit to 30m and StallLimit returns it for a command step with no stall of its own, while ship rebase runs the local gate (make gate, up to an hour) on the merged head before pushing. Without a stall the watchdog would stop and retry the rebase partway through the gate. The repo's own convention (ship-change-classic.toml and docs.toml comments, the builtin gate step) is stall = "70m" on every step that runs the gate, so following it is the least surprising answer and needs no new mechanism; the acceptance must check it so the defect cannot ship.

## statement

The new rebase step that runs rhiza ship rebase sets stall = "70m" (at least the gate step's stall), like every other step that runs ship rebase or ship gate. The leaf ship-steps-rebase-clean's acceptance requires that stall on the step, and its failure case fails when the stall is missing or below the gate's.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7nwyk3pyjjy43bsg63bzyd1-r3-f1](../question/review-01m7nwyk3pyjjy43bsg63bzyd1-r3-f1.md)

### In

- decided-by [ship-steps-rebase-clean](../capability/rhiza/workflows/built-in-workflows/ship-steps-rebase-clean.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
