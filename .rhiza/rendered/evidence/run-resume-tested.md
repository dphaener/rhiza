# run-resume-tested

evidence · `01M5WEDGF8XS9271DC46R3GNR6` · actor [claude-code-569529c3](../actor/claude-code-569529c3.md)

## statement

go test ./internal/run (TestResumeAfterNoWaitDoesNotRedo, TestResumeFinishedRunIsErrFinished, TestResumeLiveRunIsErrLive, TestResumeKilledMidStepContinuesVisit, TestResumeEndedWithoutFinishedWritesFinished) and the run-resume testscript pass; step one runs once

## Links

### Out

- verifies [run-resume](../capability/rhiza/workflows/workflow-runs/run-resume.md)

## History

- `01M5SZHZ3HG3PXE9A4SK0CA9M7` Runs and journal: genesis run start, status and resume run a workflow's agent, wait and command steps with a local journal and no daemon.
