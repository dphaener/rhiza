# scheduled-triggers

capability · `01M6W5HTS2GPF0MXKPGWSQF9DM` · parent [workflows](../workflows.md)

## statement

A repo declares schedules in .rhiza/dispatch.toml and the daemon serving the root runs each one's workflow or command on its slots, as a run with no change set on the project's main checkout

## Links

### Out

- decided-by [scheduled-triggers-shape](../../../decision/scheduled-triggers-shape.md)

### In

- verifies [pr-374-scheduled-triggers](../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
