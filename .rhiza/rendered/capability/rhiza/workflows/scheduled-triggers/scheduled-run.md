# scheduled-run

capability · `01M6W5HV83HMQEQ15319XERHJ6` · parent [scheduled-triggers](../scheduled-triggers.md)

## acceptance

Tests in internal/api and internal/run start a scheduled run of a command and of a repo workflow in a scratch repo: neither opens a change set or a worktree, both run in the main checkout, each journal records trigger schedule and the name, run ls --json and run status --json show them with no change set, the command's output reaches the run's log, and a workflow with a leaves step fails naming that step

## statement

A run can be started for a schedule with no change set: it opens none, works in the project's main checkout with no worktree or branch, runs the schedule's workflow, or its command wrapped as an implicit one-step workflow with a single command step, autonomous, and its journal records trigger schedule and the schedule's name in place of a change set, shown by rhiza run ls and rhiza run status; a step that needs a change set (such as a leaves step) fails the run with an error naming the step

## Links

### Out

- depends-on [schedules-file](schedules-file.md)
- decided-by [scheduled-triggers-shape](../../../../decision/scheduled-triggers-shape.md)

### In

- depends-on [main-checkout-change-warning](main-checkout-change-warning.md)
- depends-on [daemon-fires-schedules](daemon-fires-schedules.md)
- verifies [scheduled-run-verified](../../../../evidence/scheduled-run-verified.md)
- verifies [scheduled-run-qa-1](../../../../evidence/scheduled-run-qa-1.md)
- verifies [scheduled-run-qa-2](../../../../evidence/scheduled-run-qa-2.md)
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
