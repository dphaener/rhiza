# spawn-wake-conductor

capability · `01M7QE2048S7N0FJFG3QWQZR5T` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/run -run 'TestSpawnWake|TestSpawnResume' passes: a no-wait child that said done and one that exited 1 each produce exactly one wake call to the parent with the child's role, leaf, outcome and detail, also when the poll repeats and after a driver restart; (failure case) a waiting (not no-wait) child's end recorded by the engine wakes no one, and a wake that returns an error is logged, leaves the spawn-ended entry in place and is not retried. go test -trimpath ./internal/api -run 'TestRunSpawnWake|TestRunSpawn' passes: the stepRunner's wake closes a said-done child's pane and writes a message file under the parent's .rhiza/local/messages folder naming the child, role, leaf and outcome, and leaves the pane of a child that exited alone.

## files

internal/run/engine.go, internal/run/spawnresume_test.go, internal/api/run_lost.go, internal/api/run_spawn_test.go

## must_not

Wake the conductor for a child's question or other block (askChildren and the arbiter keep those); message a parent twice for one child; let a wake failure fail the step or the run; block the engine's poll for longer than a short message timeout (session.message queues the pointer while the parent works)

## statement

When the run engine's reconcileSpawns journals the end of a no-wait spawned child (AppendSpawnEnded added it, so once per child), it closes the child's Rex pane if it said done, as a waiting spawn does, and sends the child's parent session a rhiza session message naming the child session, its role, leaf, outcome and detail and pointing at rhiza run status for the rest; the engine reaches both through an optional Runner interface (the engine runs no program), implemented by the api stepRunner over closeSessionPane and session.message.

## Links

### Out

- depends-on [spawn-no-wait](spawn-no-wait.md)
- decided-by [spawn-no-wait-wake](../../../../decision/spawn-no-wait-wake.md)

### In

- depends-on [conduct-no-wait-guide](conduct-no-wait-guide.md)
- verifies [spawn-wake-conductor-verified](../../../../evidence/spawn-wake-conductor-verified.md)
- verifies [spawn-wake-conductor-qa-1](../../../../evidence/spawn-wake-conductor-qa-1.md)

## History

- `01M7M2FB5BNRSDSN7HFWQSC2WN` rhiza run spawn can return once the child starts, and the child's rhiza session done (or a blocked question) wakes the conductor with rhiza session message, so a conductor ends its turn instead of blocking on spawn; every child stays a Rhiza session on the run's ledger (idea from rex-teams, passed on by the orchestrator 2026-11-17)
