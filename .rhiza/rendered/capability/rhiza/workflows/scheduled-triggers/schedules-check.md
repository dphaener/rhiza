# schedules-check

capability · `01M6W5HV4DHZKJGNMDBNYDEBTN` · parent [scheduled-triggers](../scheduled-triggers.md)

## acceptance

A testscript (schedules-check) shows rhiza check passing with a schedule naming a repo workflow and one with a command, and failing with a schedule naming a missing workflow, with an unknown key and with a bad cron expression, each failure naming the schedule

## statement

rhiza check fails on a [[schedules]] entry that does not read, or whose workflow neither the repo layer (.rhiza/workflows/) nor the built-in layer resolves, naming .rhiza/dispatch.toml, the schedule's name and the key, as it does for the dispatch rules

## Links

### Out

- depends-on [schedules-file](schedules-file.md)
- decided-by [scheduled-triggers-shape](../../../../decision/scheduled-triggers-shape.md)
- decided-by [scheduled-run-readonly-cron-health](../../../../decision/scheduled-run-readonly-cron-health.md)

### In

- verifies [schedules-check-verified](../../../../evidence/schedules-check-verified.md)
- verifies [schedules-check-qa-1](../../../../evidence/schedules-check-qa-1.md)
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
