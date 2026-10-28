# run-one-live

capability · `01M61R7G3NZN2827J6Z1GD0BQ1` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript stops a fake run at a wait step with --no-wait, then a second rhiza run start on the change set exits 2 naming that run and rhiza run resume, and after the first run finishes a new start succeeds; a testscript fails a fake run, starts a second one stopped at a wait step, and rhiza run resume --retry on the failed run exits 2 naming the second; a Go test starts two runs on one change set concurrently and exactly one begins.

## statement

A change set has at most one live run: rhiza run start refuses a change set whose run's journal has not finished (done or failed), an interrupted run included, exiting 2 with the run's ID and a next step of rhiza run resume, and rhiza run resume (with or without --retry) refuses a run whose change set has another live run, exiting 2 naming that run; each checks and opens or reopens the journal under one lock so two concurrent starts or resumes never both begin.

## Links

### Out

- decided-by [answer-8cd68bc7-11](../../../../decision/answer-8cd68bc7-11.md)

### In

- depends-on [run-abandon](run-abandon.md)
- verifies [run-one-live-verified](../../../../evidence/run-one-live-verified.md)
- verifies [run-one-live-qa-1](../../../../evidence/run-one-live-qa-1.md)
- verifies [pr-50-a-change-set-has-at-most-one](../../../../evidence/pr-50-a-change-set-has-at-most-one.md)

## History

- `01M5ZNV4E4AJR52EX631180DMS` A change set has at most one live run: rhiza run start refuses a second run on a change set whose run is still going, naming that run.
