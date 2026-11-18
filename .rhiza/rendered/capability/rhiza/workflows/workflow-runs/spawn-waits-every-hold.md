# spawn-waits-every-hold

capability · `01M7QGA6MTJS6AMR4V34MR7H2W` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/api -run Spawn and go test -trimpath ./internal/run: a test saves a memory-low machine hold, then (a second case) a tmp-low hold, then a disk-low hold, in the daemon state, runs the spawn with a ctx sleep that clears the hold on its second call, and sees exactly one spawn-held entry whose detail names the reading and the limit and the role, no child started while held, then the child started and the hold released in run state (SpawnHold carries Role); (failure case) with no hold saved the spawn journals no spawn-held entry and its response's held is false

## files

internal/api/usage_hold.go, internal/api/run_spawn.go, internal/api/run_spawn_test.go, internal/run/journal.go

## must_not

Count a spawn against rhiza.daemon.maxRuns; read the host's machine, disk or temp space in tests; change --no-wait or how a child's end is recorded

## statement

rhiza run spawn waits while the daemon's last saved state holds new starts for any pre-start check (usage, load-high, memory-low, disk-low or tmp-low), journaling one spawn-held entry per wait, with its role, whose detail names the reading and the limit, and starts the child once the saved state no longer holds

## Links

### Out

- depends-on [daemon-tmp-hold](../../orchestration/daemon/daemon-picks/daemon-tmp-hold.md)
- decided-by [m40-spawn-gate-reads-daemon-state](../../../../decision/m40-spawn-gate-reads-daemon-state.md)

### In

- depends-on [step-session-waits-holds](step-session-waits-holds.md)
- verifies [spawn-waits-every-hold-verified](../../../../evidence/spawn-waits-every-hold-verified.md)
- verifies [spawn-waits-every-hold-qa-1](../../../../evidence/spawn-waits-every-hold-qa-1.md)
- verifies [spawn-waits-every-hold-qa-2](../../../../evidence/spawn-waits-every-hold-qa-2.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
