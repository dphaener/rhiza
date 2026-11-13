# run-handback

capability · `01M79KPCETKJR0AX7B576P4C76` · parent [workflow-runs](../workflow-runs.md)

## acceptance

go test -trimpath ./internal/api -run TestHandbackDeliver shows handback.md written in the run directory with the cause, a red gate step's hand-back writing handback.md naming the failing command and the path of its kept output file, a red command step's hand-back writing handback.md naming the command and its output tail, a red ship-check step's hand-back writing the ship check cause, a live session messaged with the handback text, an ended session resumed with the same session ID, and a session that cannot be resumed failing the step with a message naming the session; go test -trimpath ./cmd/rhiza -run TestScripts/run-handback shows, with a fake agent, a command step failing twice, its agent step's session handed back each time with handback.md naming the failing command, and the run then passing. Touches internal/api/run.go (a Runner hand-back method), internal/run/engine.go, internal/api/session_message.go, the contract change set's handback.md writer (a cause argument only) and cmd/rhiza/testdata/script/run-handback.txtar; it writes no handback text of its own.

## statement

When a run hands back to an agent step, it passes the routing step's cause to the contract change set's handback.md writer, extending that writer to take it: for a gate step, the failing command and the path of its kept output file from the step's detail; for a command step, the command and the output tail its detail holds (a command step keeps no output file); for any other routing step, the cause from rhiza ship check --json as the writer reads it today. It then sends that step's session the handback text the contract change set's agent-invoke-template defines: a live session is messaged as rhiza session message does, an ended session is resumed in the same conversation through its adapter, and the run then waits on that session as on the step's first pass.

## Links

### Out

- depends-on [handback-route](handback-route.md)
- depends-on [gate-step](../workflow-definitions/gate-step.md)
- decided-by [conductor-milestone-direction](../../../../decision/conductor-milestone-direction.md)
- decided-by [spawn-reuses-contract-and-waits-on-holds](../../../../decision/spawn-reuses-contract-and-waits-on-holds.md)
- decided-by [handback-split-route-deliver](../../../../decision/handback-split-route-deliver.md)

### In

- depends-on [spawn-durable](../../orchestration/daemon/spawn-durable.md)
- depends-on [spawn-visible](spawn-visible.md)
- verifies [run-handback-verified](../../../../evidence/run-handback-verified.md)
- verifies [run-handback-qa-1](../../../../evidence/run-handback-qa-1.md)
- verifies [run-handback-qa-2](../../../../evidence/run-handback-qa-2.md)
- verifies [pr-572-the-daemon-offers-spawn](../../../../evidence/pr-572-the-daemon-offers-spawn.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
