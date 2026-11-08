# daemon-fires-schedules

capability · `01M6W5HVM44MDKXV4ESTPX0KSX` · parent [scheduled-triggers](../scheduled-triggers.md)

## acceptance

Tests in internal/api with a fake clock show a daemon poll over two projects starting each due schedule's run, a due schedule waiting while maxRuns is full and then firing once, an unresolvable workflow starting nothing with a log line naming the schedule, and the daemon log carrying a scheduled run's start and its failure; for a scheduled command run that starts and fails, no Rex call is made

## statement

The daemon serving the root reads every registered project's schedules on each poll and starts each due schedule's run, taking a slot under rhiza.daemon.maxRuns and the project's maxRuns like any run; a due schedule with no free slot waits for one and still fires once, a schedule whose file no longer reads or whose workflow no longer resolves starts nothing and is logged; it logs each scheduled run's start, end and result, and Rhiza opens no Rex tab or pane of its own for a scheduled run, its start or its failure (failures go to the daemon log); an agent step inside a scheduled workflow opens its session as in any headless daemon run

## Links

### Out

- depends-on [schedule-due](schedule-due.md)
- depends-on [scheduled-run](scheduled-run.md)
- depends-on [schedules-file](schedules-file.md)
- decided-by [scheduled-triggers-shape](../../../../decision/scheduled-triggers-shape.md)

### In

- verifies [daemon-fires-schedules-verified](../../../../evidence/daemon-fires-schedules-verified.md)
- verifies [daemon-fires-schedules-qa-1](../../../../evidence/daemon-fires-schedules-qa-1.md)
- verifies [daemon-fires-schedules-verified-2](../../../../evidence/daemon-fires-schedules-verified-2.md)
- verifies [daemon-fires-schedules-qa-2](../../../../evidence/daemon-fires-schedules-qa-2.md)
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
