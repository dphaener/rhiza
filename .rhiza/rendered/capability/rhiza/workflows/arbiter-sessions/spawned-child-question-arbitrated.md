# spawned-child-question-arbitrated

capability · `01M7M5E14QHSACYM4E06M4E3BW` · parent [arbiter-sessions](../arbiter-sessions.md)

## acceptance

go test -trimpath ./internal/run -run 'TestSpawnedChildQuestion': a new TestSpawnedChildQuestionArbitrated runs a conductor step (agent step declaring checks) in a run not started autonomous whose journal has a spawn entry for a child session the fake runner reports blocked on a question while the conductor's own session works: AddQuestion gets the child's question, an asked entry names the child session, StartArbiter runs on the question's ID, Deliver goes to the child session and an answered entry follows, with no StartAdvisor call and no Notify, and the step then finishes on the conductor's own outcome; (failure case) a failing arbiter verdict on the child's question ends the step with an error outcome naming the question; a child with a spawn-ended entry is never polled.

## files

internal/run/engine.go (pollAgent, ask), internal/run/engine_test.go

## must_not

Arbitrate a child's non-question block (a permission prompt stays the owner's); launch an advisor; change how reconcileSpawns records a child's end; make the conductor step's outcome depend on a child's outcome

## statement

While a conductor step runs, the engine watches the step's live spawned children, and a child blocked on a question gets the same arbiter handling as the conductor's own: a question node, an asked entry keyed on the child's session, the arbiter, and the answer delivered to the child, its block ended; no advisor and no owner notice.

## Links

### Out

- depends-on [conductor-question-arbitrated](conductor-question-arbitrated.md)
- depends-on [conductor-sessions-launch-autonomous](conductor-sessions-launch-autonomous.md)
- decided-by [conductor-questions-to-arbiter](../../../../decision/conductor-questions-to-arbiter.md)

### In

- verifies [spawned-child-question-arbitrated-verified](../../../../evidence/spawned-child-question-arbitrated-verified.md)
- verifies [spawned-child-question-arbitrated-qa-1](../../../../evidence/spawned-child-question-arbitrated-qa-1.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand
