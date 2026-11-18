# stall-spawned-children-count

decision · `01M7QBVD59SZ7JP46RKEEGHFYN`

## statement

A run step's stall detection counts the live children of its open spawns (run.State.Spawns of the current step not yet ended) as the step's sessions, so an idle child with no done mark is nudged at the first limit and the step is stopped and resumed with --retry at the second, charged to the run's watchdog budget; a session waiting on an open spawn of its own is neither nudged nor counted stalled by its own silence, since its progress is its child's; the stop closes the panes of the step's children and of the waiting parent; the retried step's fresh conductor carries on from the journal as a part-way conduct does

## wording

(assumed, would have escalated) the conductor's plan of an autonomous run, derived from the milestone brief's point 2: nudge, then retry the step, for every workflow's agent steps

## Links

### Out

- derived-from [m40-stalls-conduct-plan](../source/m40-stalls-conduct-plan.md)

### In

- decided-by [stall-covers-spawned-children](../capability/rhiza/orchestration/daemon/daemon-watchdog/stall-covers-spawned-children.md)

## History

- `01M7NVMG0EJFRJV03G025C6366` The daemon's stall detection covers every agent session of any workflow's run, including the children a conductor spawns: an idle child with no done mark (the fixer that held a run 9 hours on 2026-11-16) is nudged, then its step retried or failed with a reason, replacing the check-in alarms of the withdrawn draft 01M7M2FE1S and the orchestrator's runwatch
