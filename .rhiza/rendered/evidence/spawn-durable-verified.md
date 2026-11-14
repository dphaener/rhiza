# spawn-durable-verified

evidence · `01M7BDWBHWT4391GQ3WZ4GC0E6` · actor [claude-code-8b038edb](../actor/claude-code-8b038edb.md)

## result

PASS

## statement

Verified at e66d8e9d4. Re-ran the Checks block: go test -trimpath ./internal/run -run TestSpawnResume (both subtests PASS: children across a driver restart; a pending hand-back delivered once after a restart), go test -trimpath ./internal/run/... ok, go test -trimpath ./internal/api -run 'TestRun|TestStepRunner|TestHandback|TestSpawn|TestLost|TestSessionRun' ok (incl. new TestStepRunnerChild), TestScripts/run-spawn, run-handback, run-resume, run-resume-agent-session all PASS. Read the diff: the driver is killed (panic) while children run and while a hand-back is being delivered, then restarted; K-live is polled until it ends (reattached), K-ended is recorded once even though its spawn process also records it (Store.AppendSpawnEnded skips an already-ended child, and the fold ignores a second spawn-ended entry), K-gone is recorded fail with 'session record is gone', the journaled hand-back (State.PendingHandback, cleared by the new handback-delivered entry) is delivered exactly once, and the run continues from the same step to Succeeded. The resume path lives in internal/api/run_lost.go (stepRunner.Child) and run_spawn.go rather than run.go; it is in scope. There are no stubs or TODOs, and only fakes are used. make self-check fails only on .rhiza/rendered actor files whose actor log entries the run keeps uncommitted in its root. The base commit HEAD~1 has the same drift, so the leaf's code does not cause it. Note: AppendSpawnEnded checks the journal and then appends without a lock, so a simultaneous driver and spawn process could each write an entry; the fold still records the child's end once.

## Links

### Out

- verifies [spawn-durable](../capability/rhiza/orchestration/daemon/spawn-durable.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
