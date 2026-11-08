# main-checkout-change-warning

capability · `01M6W5HVA0S11V8P7TJPCY6QH8` · parent [scheduled-triggers](../scheduled-triggers.md)

## acceptance

Tests in internal/api show a scheduled command that leaves the checkout alone recording no warning, and commands that switch branch, make a commit, and modify a tracked file each ending with the run's status as the command left it plus a warning in run status --json and the log naming the branch, the new HEAD or the file

## statement

When a scheduled run ends, Rhiza compares the main checkout's branch, HEAD and uncommitted tracked changes with their state when the run started, and if the run changed any of them the run's record and log carry a warning naming what changed; nothing is reverted or blocked

## Links

### Out

- depends-on [scheduled-run](scheduled-run.md)
- decided-by [scheduled-run-readonly-cron-health](../../../../decision/scheduled-run-readonly-cron-health.md)

### In

- depends-on [health-schedules-part](health-schedules-part.md)
- verifies [main-checkout-change-warning-verified](../../../../evidence/main-checkout-change-warning-verified.md)
- verifies [main-checkout-change-warning-qa-1](../../../../evidence/main-checkout-change-warning-qa-1.md)
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
