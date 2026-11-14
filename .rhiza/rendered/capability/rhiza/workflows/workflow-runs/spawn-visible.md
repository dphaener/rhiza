# spawn-visible

capability · `01M79KPCC9NH1M9F5EFKVGJW69` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/api -run TestRunStatusSpawned shows rhiza run status and --json listing two children of one step with role, leaf, model, state, outcome and session ID, a running child marked running, and a failed child with its reason, and rhiza session show accepting a child's session ID; go test -trimpath ./internal/api -run TestSpawnedSessionOpenResume shows, with a fake Rex and a fake agent, rhiza session pane reading a running child's terminal by its session ID and rhiza session message reaching that running child, and an ended child resumed in the same conversation (same session ID) through the Runner's hand-back resume path from run-handback; go test -trimpath ./internal/web -run TestRunDetailSpawned shows the web run page listing the same children under their step, the leaf named by its statement, with no session, run or node ID on the page. Touches internal/api/run.go (RunInfo), internal/run/status.go, internal/web/rundetail.go and its template under internal/web/templates; the describe-golden testscript is regenerated with -update.

## statement

rhiza run status, its --json and the web run page list every spawned session under the step that spawned it, with its role, leaf, model, state and outcome, and open or resume it like any other session.

## Links

### Out

- depends-on [run-spawn](run-spawn.md)
- depends-on [run-handback](run-handback.md)
- decided-by [conductor-spawned-sessions](../../../../decision/conductor-spawned-sessions.md)

### In

- verifies [spawn-visible-verified](../../../../evidence/spawn-visible-verified.md)
- verifies [spawn-visible-qa-1](../../../../evidence/spawn-visible-qa-1.md)
- verifies [spawn-visible-qa-2](../../../../evidence/spawn-visible-qa-2.md)
- verifies [pr-572-the-daemon-offers-spawn](../../../../evidence/pr-572-the-daemon-offers-spawn.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
