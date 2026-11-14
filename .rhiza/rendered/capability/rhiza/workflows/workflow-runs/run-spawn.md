# run-spawn

capability · `01M79KPCATRS9TD6XGS3C3KYY6` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/api -run TestRunSpawn covers a spawned child's journal entry under the calling step, its brief.md holding the contract section then the caller's brief, the outcome returned (pass, and a failed child with exit 1 and its reason), a refusal outside a run naming RHIZA_RUN, a refusal of a missing --brief file naming the path, a wait journaled as held while a usage hold is on that ends when the hold lifts, and a wait journaled as held while the daemon state records a machine hold that ends when the state clears it, with maxRuns at 1 holding neither; go test -trimpath ./cmd/rhiza -run TestScripts/run-spawn passes with a fake agent; go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden passes after -update. Touches internal/api/run_spawn.go (new), internal/api/run_leaves.go, internal/api/usage_hold.go, internal/run/journal.go, the contract change set's brief.md writer (a spawn mode taking the caller's brief only), cmd/rhiza/testdata/script/run-spawn.txtar and describe-golden.txtar; it must not start the child from the daemon process or count it against rhiza.daemon.maxRuns.

## statement

An agent session in a run starts a child session with rhiza run spawn --role <role> [--leaf <leaf>] --brief <file> [--model <model>]: it starts the child in the run's worktree through the run's agent adapter with the run's permission mode, writes the child's brief.md through the contract change set's brief writer, extending that writer with a spawn mode that writes only its contract section followed by the caller's brief file, waits journaled as held while a usage hold or the machine load hold (the daemon state's MachineHold, read as the inbox reads it) is on (never on maxRuns), journals the child as a child of the calling step, and prints the child's outcome as JSON when it ends.

## Links

### Out

- decided-by [conductor-spawned-sessions](../../../../decision/conductor-spawned-sessions.md)
- decided-by [spawn-reuses-contract-and-waits-on-holds](../../../../decision/spawn-reuses-contract-and-waits-on-holds.md)
- decided-by [spawn-starts-from-caller](../../../../decision/spawn-starts-from-caller.md)

### In

- depends-on [spawn-visible](spawn-visible.md)
- depends-on [spawn-durable](../../orchestration/daemon/spawn-durable.md)
- verifies [run-spawn-verified](../../../../evidence/run-spawn-verified.md)
- verifies [run-spawn-qa-1](../../../../evidence/run-spawn-qa-1.md)
- verifies [run-spawn-qa-2](../../../../evidence/run-spawn-qa-2.md)
- verifies [pr-572-the-daemon-offers-spawn](../../../../evidence/pr-572-the-daemon-offers-spawn.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
