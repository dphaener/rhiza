# spawn-wake-conductor-verified

evidence · `01M7QEM6AEJ14K8337XS1MKNZY` · actor [claude-code-b5491234](../actor/claude-code-b5491234.md)

## result

PASS

## statement

At e94526667 re-ran the Checks: block: go test -p 2 -trimpath ./internal/run -run 'TestSpawnWake|TestSpawnResume' (PASS, both TestSpawnWake subtests), go test -p 2 -trimpath ./internal/run (ok), go test -p 2 -trimpath ./internal/api -run 'TestRunSpawnWake|TestRunSpawn' (TestRunSpawnWake SaidDone/Exit1 PASS), go vet ./internal/run ./internal/api, rhiza ship generate, rhiza render: all exit 0. Diff read: engine.reconcileSpawns calls the optional SpawnWaker only when AppendSpawnEnded added the entry and the spawn is NoWait, so one wake per child, none for a waiting child, a wake error is only logged; TestSpawnWake covers said-done and exit-1 children across repeated polls and a driver restart, a waiting child waking no one, and an erroring wake logged once with its spawn-ended entry kept. stepRunner.WakeParent closes a said-done child's pane via closeSessionPane and messages the parent (queued while working, typed with a 30s timeout when idle) naming child, role, leaf, outcome, detail and rhiza run status; TestRunSpawnWake checks the pane close, the untouched pane of an exited child, the message file under .rhiza/local/messages/parent-1 and the queued pointer. Only the leaf's four files plus its model log and render changed; no stubs or TODOs.

## Links

### Out

- verifies [spawn-wake-conductor](../capability/rhiza/workflows/workflow-runs/spawn-wake-conductor.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)
