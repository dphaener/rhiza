# spawn-waits-every-hold-verified

evidence · `01M7QHN5FPG30ZTE3XRDZDY1KB` · actor [claude-code-1e7c437e](../actor/claude-code-1e7c437e.md)

## result

PASS

## statement

Re-ran b4bc38687's Checks in capped scopes: go test -trimpath -p 2 ./internal/api -run Spawn (exit 0; TestRunSpawnHolds/MemoryLow, TmpLow, DiskLow and NoHold pass), go test ./internal/run (exit 0), go vet ./internal/api ./internal/run (exit 0), TestScripts/(run-spawn|describe-golden)$ (exit 0, golden current without -update), rhiza ship generate and rhiza render (exit 0, no diff beyond this verifier's own actor). Diff read: spawnHoldDetail now holds on MachineHold (prefixed with its reason), DiskLow and TmpLow from the saved daemon state, naming reading and limit; the test saves each hold with daemon.WriteState, clears it on the ctx sleep's second call, asserts no child starts while held, exactly one spawn-held entry with role verifier and the reading/limit in its detail, the hold released in run state with SpawnHold.Role, and one child spawned; NoHold journals no entry and held is false. maxRuns, --no-wait and child-end recording untouched; no host reads, no TODOs.

## Links

### Out

- verifies [spawn-waits-every-hold](../capability/rhiza/workflows/workflow-runs/spawn-waits-every-hold.md)

## History

- `01M7NVMAZ0PV2BP11NSNF04ANZ` The daemon's pre-start checks (load, memory, disk) gate every child session spawn of any workflow, not only run starts, and a free-tmp check joins them: a held spawn waits rather than starting into an OOM, and its run journal and the inbox say which reading held it
