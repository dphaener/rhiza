# agent-step-tested

evidence · `01M5WEDEYP93QEAWNTFV0S9ABN`

## statement

go test ./internal/run (TestAgentPollsUntilDone, TestAgentFailEndsRunFailed, TestAgentResumedDoesNotStartAgain, TestAgentStartErrorFailsRun), go test ./internal/api (TestRunAgentStepEndsWithTheSessionExit, which checks the session records the rendered prompt), go test ./internal/workflow (TestAgentPrompt), go test ./internal/adapter/claudecode (TestCommandPrompt) and the run-agent-step testscript pass; in the testscript the fake Claude Code's last argument is the step's prompt with {change}, {slug}, {title}, {step} and {run} filled, and a workflow with no prompt or an unknown placeholder is refused before anything runs

## Links

### Out

- verifies [agent-step](../capability/rhiza/workflows/workflow-definitions/agent-step.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
