# failed-steps-merge-resolution

decision · `01M69QNAHHF8T0Q50E1ZE21G6R`

## statement

The conflict in internal/api/run.go is settled by keeping both: a command step exiting with run.StopExit keeps its last stop: line as detail (from #137), and any other failure keeps the capped tail of its output; the branch merges main instead of rebasing

## wording

(the orchestrator's ruling as arbiter, under the owner's standing instruction)

## Links

### Out

- answers [rebase-conflict-failed-steps-keep-output](../question/rebase-conflict-failed-steps-keep-output.md)

## History

- `01M686RRCJSQ57M82FENWT2HNM` When a run's command step fails, its journal detail keeps the tail of the step's output, shown by run status and the inbox, so an unattended failure says why
