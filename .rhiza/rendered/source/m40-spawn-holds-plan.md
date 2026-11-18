# m40-spawn-holds-plan

source · `01M7QG8KE8VEP99YXQSVEMM8SC`

## quotes

```json
[]
```

## raw_sha256

19238cc54189a3164e74481f913cf280460741d8ecf33050f9b4d6afb8884c9a

## session

6997790c-fb5f-4105-9548-01a1ef028399

## summary

Autonomous conductor plan, no owner present. Research found: the daemon's poll holds new starts and resumes on disk-low (rhiza.daemon.minFree, repo filesystem) and on load-high or memory-low (rhiza.daemon.maxLoad, rhiza.daemon.minMemory, internal/daemon/load.go Evaluate), saved in the daemon state file; rhiza run spawn (internal/api/run_spawn.go) already waits on the saved usage hold and machine hold (spawnHoldDetail in internal/api/usage_hold.go) and journals one spawn-held entry per wait, but not on disk-low; the run engine starts a workflow step's agent session (internal/run/engine.go agent, Runner.StartAgent) with no hold at all; no free-tmp check exists; the journal's SpawnHolds reach no inbox item. The plan adds a free-tmp hold to the daemon poll, makes rhiza run spawn and the engine's fresh step sessions wait on every saved pre-start hold, and shows a held spawn or step session in the inbox naming the reading and limit.

## title

Autonomous plan of the daemon's pre-start checks (m40 spawn holds)

## Links

### In

- derived-from [m40-tmp-hold-shape](../decision/m40-tmp-hold-shape.md)
- derived-from [m40-spawn-gate-reads-daemon-state](../decision/m40-spawn-gate-reads-daemon-state.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
