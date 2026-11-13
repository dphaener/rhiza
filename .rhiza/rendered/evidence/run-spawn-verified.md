# run-spawn-verified

evidence · `01M7BAGS6W5RFJGASWH2QTEEKG` · actor [claude-code-4824b997](../actor/claude-code-4824b997.md)

## result

PASS

## statement

Verified at 83d3c286d (leaf commits 8fe16199f, 5d3539d80, 83d3c286d). Re-ran: go test -trimpath ./internal/api -run TestRunSpawn (exit 0; Pass, Fail exit 1, OutsideRun names RHIZA_RUN, MissingBrief names path, UsageHold and MachineHold journaled as held then lifted, MaxRunsHoldsNothing); go test -trimpath ./cmd/rhiza -run TestScripts/run-spawn$ with a fake agent (exit 0); TestScripts/describe-golden after -update (exit 0, golden unchanged); go test ./internal/run/... ./internal/session/... and the scoped internal/api set (exit 0); rhiza ship generate (0 changed) and rhiza render (exit 0); make self-check from a clean detached worktree with RHIZA_* unset (exit 0; in the run worktree it fails only on this verifier's own unrendered actor). Read the diff: the child starts in the caller's process via session.start with the run's agent and permission mode, brief.md is RenderSpawnBrief (contract section then the caller's brief, no slice.json), spawn/spawn-held/spawn-ended journal entries under the calling step, holds read from daemon state UsageHold/MachineHold as the inbox does plus the governor limit, never maxRuns. No TODOs or stubs.

## Links

### Out

- verifies [run-spawn](../capability/rhiza/workflows/workflow-runs/run-spawn.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
