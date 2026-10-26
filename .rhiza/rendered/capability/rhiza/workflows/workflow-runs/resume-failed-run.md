# resume-failed-run

capability · `01M5X3GCT3DRQ1XEGAZ630X0KS` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript where a command step fails, the cause is fixed, and resume --retry runs that step again and the run succeeds; the journal shows the failure and the retry.

## statement

rhiza run resume --retry restarts a failed run at the step that failed, once the cause is fixed, and journals the retry; without --retry, resume still refuses a finished run.

## Links

### Out

- decided-by [ship-change-holds-up](../../../../decision/ship-change-holds-up.md)

### In

- verifies [resume-failed-run-tests](../../../../evidence/resume-failed-run-tests.md)

## History

- `01M5X3GBH4ZQTZVXFEHR966747` ship-change that holds up: the agent step passes only when the agent says it is done, land renders with the branch's own build in Rhiza's repo, a PR that main has moved past is rebased, and a failed run can be resumed from the step that failed.
