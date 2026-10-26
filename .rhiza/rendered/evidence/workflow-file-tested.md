# workflow-file-tested

evidence · `01M5WEDEQYKCMH8GPJJ4Y8VSDF`

## statement

go test ./internal/workflow (TestLoadValid, TestLoadRefused, TestParseErrors, TestParseOrderAndStart, TestNext) and the TestScripts/run-start workflow refusals (api TestRunStartWorkflowRefusals) pass

## Links

### Out

- verifies [workflow-file](../capability/rhiza/workflows/workflow-definitions/workflow-file.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
