# health-schedules-part

capability · `01M6W5HW3CJ9HJ6KKRS7G4CJXC` · parent [scheduled-triggers](../scheduled-triggers.md)

## acceptance

Tests in internal/api show a project with no schedules having no schedules part, a schedule not yet run listed ok with its next run, a failed last run making the part warn with next naming that run's log, a last run that changed the checkout making it warn, no daemon making it warn, and health exiting 0 in each; describe-golden regenerated for the new part enum value

## statement

rhiza health shows a schedules part for each project that has schedules: its summary counts them and each schedule is a detail line with its last run (time and run ID), result and next run; it is warn when any schedule's last run failed or changed the main checkout, or no daemon is running to fire them, with next pointing at the failing run's log, ok otherwise, never bad, and absent for a project without schedules

## Links

### Out

- depends-on [main-checkout-change-warning](main-checkout-change-warning.md)
- decided-by [scheduled-run-readonly-cron-health](../../../../decision/scheduled-run-readonly-cron-health.md)

### In

- verifies [health-schedules-part-verified](../../../../evidence/health-schedules-part-verified.md)
- verifies [health-schedules-part-qa-1](../../../../evidence/health-schedules-part-qa-1.md)
- verifies [health-schedules-part-qa-2](../../../../evidence/health-schedules-part-qa-2.md)
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
