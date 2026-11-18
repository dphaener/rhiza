# m40-spawn-gate-reads-daemon-state

decision · `01M7QGA6GQ8FTEZKCRYE6PYX16`

## statement

A child session spawn (rhiza run spawn) and a workflow step's fresh agent session both wait on the holds in the daemon's last saved state (usage, load-high, memory-low, disk-low, tmp-low), the source rhiza run spawn already reads, not on a reading of their own, so every limit and the load hysteresis keep one meaning; each wait journals one spawn-held entry (a step session's with no parent), released by the spawn or the step's session entry

## wording

(assumed by the autonomous planner from the existing rhiza run spawn hold and Decision m40-host-limits-as-holds)

## Links

### Out

- derived-from [m40-spawn-holds-plan](../source/m40-spawn-holds-plan.md)

### In

- decided-by [spawn-waits-every-hold](../capability/rhiza/workflows/workflow-runs/spawn-waits-every-hold.md)
- decided-by [step-session-waits-holds](../capability/rhiza/workflows/workflow-runs/step-session-waits-holds.md)
- decided-by [inbox-held-spawns](../capability/rhiza/orchestration/inbox/inbox-list/inbox-held-spawns.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
