# workflow-cleanup-steps-qa-1

evidence · `01M7758GT4E9ZHFCDFYZXNJCRZ` · actor [claude-code-3c7e77a4](../actor/claude-code-3c7e77a4.md)

## result

PASS

## statement

QA at 333b5bf65 with a scratch build in a scratch root and repo: run start --foreground of a repo workflow whose command step exits 3 ended failed (exit 1) and then journalled both [[cleanup]] steps (named 'first' and default 'cleanup-2') as pass, in order, after the finished entry; the first ran with pwd = the run worktree and RHIZA_RUN, RHIZA_CHANGE and RHIZA_STEP=first set. A workflow whose step passes with cleanups 'exit 7', an unknown command and a touch ended succeeded (exit 0), logged 'cleanup step breaks failed: exit 7' and 'nonexistent-cmd failed: exit 127' on stderr, journalled both as fail with output tail, and the third still ran (file created). A run held at a wait step (--foreground --no-wait) then run abandon: finished failed 'abandoned by', cleanup 'mark' journalled pass with pwd the worktree, before the worktree was removed; a second abandon is refused and does not rerun it. workflow describe --json of every built-in and repo workflow shows an empty cleanup list.

## Links

### Out

- verifies [workflow-cleanup-steps](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-cleanup-steps.md)

## History

- `01M763YS08ZM13JZ065Q95YA2E` Workflow cleanup phase: a workflow file declares command steps that run once when a run of it ends, done, failed or abandoned; a failing cleanup step is logged and never changes the run's result; workflow check and describe cover it.
