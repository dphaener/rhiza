# run-start-refuses-agent-session-verified

evidence · `01M65FQFAVF7BMNZ3ZGARNA4A5` · actor [claude-code-e37730d2](../actor/claude-code-e37730d2.md)

## result

PASS

## statement

At b00d77a, go test ./cmd/rhiza -run TestScripts/run-start-agent-session passes: with RHIZA_SESSION=01SESSIONFORTESTING, rhiza run start --workflow one exits 2 with 'run start refused: this is agent session 01SESSIONFORTESTING, and no agent session starts a run' and no .git/rhiza/runs; with RHIZA_SESSION empty the same start succeeds. By hand in a scratch repo with RHIZA_SESSION=abc-123, run start (plain and --json) exits 2 naming abc-123, and neither .git/rhiza nor a sibling run worktree is created. The diff puts the guard first in runRunStart, before the workflow is resolved or anything is written. make check exits 0.

## Links

### Out

- verifies [run-start-refuses-agent-session](../capability/rhiza/workflows/workflow-runs/run-start/run-start-refuses-agent-session.md)

## History

- `01M61GW73H05YHE4DK0CXFD1KH` No agent session starts a run: rhiza run start refuses inside an agent session (RHIZA_SESSION set), naming the session, and the build and review guides tell verifiers and checks to test runs with fakes only, never real Rex, agents or gh.
