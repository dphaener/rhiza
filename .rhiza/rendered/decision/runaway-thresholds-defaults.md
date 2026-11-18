# runaway-thresholds-defaults

decision · `01M7QN2R9KVQ1ZW3APMFYFQN9M`

## statement

Three watchdog keys, machine and project like the others: runawayProcs (default 400 processes), runawayMemory (default 8G of summed resident memory, a size as rhiza.daemon.minTmpFree parses it) and runawayGrace (default 5m); 0 turns a threshold off. The defaults sit above this VM's capped test scope (TasksMax=300, MemoryMax=6G), so a session running one capped test is not a runaway

## wording

(assumed, would have escalated) planner's choice under the autonomous rules

## Links

### Out

- derived-from [runaway-plan-notes](../source/runaway-plan-notes.md)

### In

- decided-by [runaway-threshold-keys](../capability/rhiza/orchestration/daemon/daemon-watchdog/runaway-threshold-keys.md)

## History

- `01M7NVMDMBVBKRW5B8SB0KMW6T` The daemon's watchdog treats a runaway agent session of any workflow, one whose process tree's process count or resident memory grows past a host threshold read portably (ps, on Linux and macOS), like a stall: it nudges the session, then fails the step with the reading as its reason
