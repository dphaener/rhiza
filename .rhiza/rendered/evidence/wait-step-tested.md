# wait-step-tested

evidence · `01M5WEDFC4Q9APZS6NNWZDV6RK` · actor [claude-code-569529c3](../actor/claude-code-569529c3.md)

## statement

go test ./internal/run (TestWaitHeldOnThirdPoll, TestWaitTimesOut, TestWaitSleepsOnlyToDeadline) with a fake clock, and go test ./internal/api (TestRunStopsAtWaitAndResumes) pass

## Links

### Out

- verifies [wait-step](../capability/rhiza/workflows/workflow-definitions/wait-step.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
