# conductor-question-arbitrated

capability · `01M7M5E0WTFG4BFTK6Z7J15C8H` · parent [arbiter-sessions](../arbiter-sessions.md)

## acceptance

go test -trimpath ./internal/run -run 'TestConductorQuestion|TestNormalRunQuestionNotArbitrated|TestAutonomousQuestionGoesToArbiter|TestAutonomousArbiterFailFailsRun': a new TestConductorQuestionArbitratedInNormalRun runs a workflow not started autonomous whose agent step declares checks (once at top level and once included through a workflow step, as ship-rhiza includes ship-change/conduct); its session blocked on a question gets an AddQuestion with the asked text, an asked entry, StartArbiter on the question's ID and Deliver of the answer to that session, with no StartAdvisor call, no advisor entry and no Notify; (failure case) a failing arbiter verdict ends that step with an error outcome naming the question, as in an autonomous run; TestNormalRunQuestionNotArbitrated still shows a plain agent step's (no checks) question notified to the owner with no arbiter and no asked entry.

## files

internal/run/engine.go (arbitrates, block, unblock, pollAgent), internal/run/engine_test.go, internal/workflow/workflow.go (a Step predicate for a conductor step, if one is added)

## must_not

Change State.Autonomous, the signoff step, the gates or daemonStart; force --autonomous on any run; launch an advisor or notify the owner for a conductor step's question; arbitrate a non-question block

## statement

In any run, autonomous or not, a question the session of a conductor step (an agent step that declares checks, such as ship-change's conduct, nested through a workflow step included) is blocked on goes to the arbiter exactly as an autonomous run's does: a question node, an asked entry, the arbiter, the answer delivered, with no advisor launched and no owner notice; every other step's question in a run not started --autonomous still goes to the owner with an advisor.

## Links

### Out

- decided-by [conductor-questions-to-arbiter](../../../../decision/conductor-questions-to-arbiter.md)

### In

- depends-on [spawned-child-question-arbitrated](spawned-child-question-arbitrated.md)
- verifies [conductor-question-arbitrated-verified](../../../../evidence/conductor-question-arbitrated-verified.md)
- verifies [conductor-question-arbitrated-qa-1](../../../../evidence/conductor-question-arbitrated-qa-1.md)

## History

- `01M7M51W2S30H06K0BFZ6MK7PF` Fix In a conductor run, the conductor's planning questions go through AskUserQuestion and block the run: the engine spawns an advisor (which only recommends to the owner) and no arbiter decides, so an autonomous run with no owner at the pane stalls until someone answers by hand
