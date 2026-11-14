# spawn-visible-verified

evidence · `01M7BDWFX9KXY40V3G5MY93CDQ` · actor [claude-code-6bbdd14c](../actor/claude-code-6bbdd14c.md)

## result

PASS

## statement

At 123416aea (leaf commits 7655ff31c and fix 123416aea) re-ran every Checks: and acceptance command, all exit 0: go test -trimpath ./internal/api -run TestRunStatusSpawned (two children of step build with role, leaf, statement, model, state, outcome and session ID in the response and human text, the running child marked running, the failed child with 'exit 3', and session.show accepting the child's ID); ./internal/api -run TestSpawnedSessionOpenResume (fake Rex: session.pane reads the child's pane, session.message reaches it, and after it ends stepRunner.Handback resumes the same session ID); ./internal/web -run TestRunDetailSpawned (children under their step, leaf by statement, no session ID or handle on the page); ./internal/run/... (TestDeriveStatusSpawned); TestScripts/describe-golden with no drift. make self-check exits 0 in a clean detached worktree with env -i, which resolves attempt 1's finding. The diff touches internal/api/run.go (RunInfo.Spawned, RunSpawned, human output), internal/run/status.go, internal/web/rundetail.go, docs/reference/cli.md and the golden. The web page reuses the existing Session component, so no template change was needed. No stubs or TODOs.

## Links

### Out

- verifies [spawn-visible](../capability/rhiza/workflows/workflow-runs/spawn-visible.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
