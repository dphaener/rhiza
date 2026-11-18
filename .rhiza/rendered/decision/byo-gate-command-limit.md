# byo-gate-command-limit

decision · `01M7RB61YSZ20BFXPPM6NEQ6P3`

## statement

A gate step with no stall= gets the watchdog's command limit (commandLimit, 30m default), as a command step does, instead of the agent retry limit it falls through to today; the built-in gates keep their explicit stall=70m

## wording

(planner, from the code: StallLimit in internal/api/watchdog_policy.go switches on command, wait and signoff only, so the BYO example's gate got the agent limit)

## Links

### Out

- derived-from [byo-daemon-plan-research](../source/byo-daemon-plan-research.md)

### In

- decided-by [byo-watchdog-stall-runaway](../capability/rhiza/workflows/workflow-runs/byo-watchdog-stall-runaway.md)

## History

- `01M7NVN06JNCCW1NVGPD4WZB16` A bring-your-own workflow (docs/examples/byo) runs end to end under the daemon in a test, with the spawn holds, stall and runaway detection, failure retry and intervention record applying to its steps as they do to the built-in workflows
