# review-01m7nwyk3pyjjy43bsg63bzyd1-r3-f1

question · `01M7NXYSW3GQH10DF4RHEKKQA3`

## statement

Plan review round 3 left finding f1 (contradiction on ship-steps-rebase-clean): The new rebase step runs plain rhiza ship rebase, which runs the local gate (make gate, up to an hour) on the merged head before the push (internal/api/ship_rebase.go header). But the leaf's statement and acceptance, and the decision ship-steps-rebase-clean-merge, give it no stall. With no stall, a command step gets the watchdog's DefaultCommandLimit of 30m (internal/api/watchdog_policy.go:26, StallLimit at :151), so the daemon's watchdog would stop and retry the rebase partway through the gate. ship-change-classic.toml:60-71 notes exactly this and sets stall = "70m" on every step that runs ship rebase or ship gate (rebase, rebase-merge, gate-rebase-merge and others); the builtin's gate step also sets stall = "70m". The acceptance should require stall = "70m" (at least the gate's) on the rebase step, and its failure case should fail when it is missing. The critic confirmed it: Checked: internal/api/ship_rebase.go's header says ship rebase runs the project's local gate on the merged head before pushing while CI is unstable. watchdog_policy.go:26 sets DefaultCommandLimit to 30m, and StallLimit (:151) uses it for a command step that has no stall. ship-change-classic.toml:60-71 says every step that runs the gate needs a stall above the gate's hour, and sets stall = "70m" on each one; the builtin's own gate step (ship-change.toml:50) does the same. The leaf's statement and acceptance give the new rebase step no stall, so the watchdog can stop and retry it partway through the gate. That is a real defect the acceptance would let through, so major stands. How should it be settled?

## Links

### Out

- blocks [ship-steps-rebase-clean](../capability/rhiza/workflows/built-in-workflows/ship-steps-rebase-clean.md)

### In

- answers [ship-steps-rebase-stall-70m](../decision/ship-steps-rebase-stall-70m.md)

## History

- `01M7MN70WCHDYN7H7T01AVZ2VK` Fix When a PR's merge is refused because origin/main moved past the checked head, or GitHub reports a conflict, the ship-change merge step hands the whole change set back to the conductor (a rejoin) instead of rebasing in the ship steps themselves: in milestone 38 every PR (#588, #590, #591, #593) took an extra conductor visit, gate and CI round for a clean rebase. Rebase (rhiza ship rebase) and re-check in the ship steps when the merge of main is clean, and hand back only on a real conflict
