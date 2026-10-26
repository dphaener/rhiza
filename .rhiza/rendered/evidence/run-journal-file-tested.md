# run-journal-file-tested

evidence · `01M5WEDG17KHB6F80M9ZAMM5B1`

## statement

go test ./internal/run (TestRoundTrip, TestReadPrefix, TestTornLineIgnored, TestFoldFinished, TestTryLockExclusion) and the run-journal testscript pass, which finds the change set's log unchanged after a run

## Links

### Out

- verifies [run-journal-file](../capability/rhiza/workflows/workflow-runs/run-journal-file.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
