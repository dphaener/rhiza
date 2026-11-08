# question-d0780667-2

question · `01M6W4WZHD7J536GFBW9DD7HHE`

## session

d0780667-0b40-44e7-a17d-01a1b849f0e6

## statement

Decision scheduled-triggers-shape says a scheduled run uses the repo's main checkout, 'read-only by default'. Following warn-dont-safeguard, what does read-only mean in practice? / Which cron dialect and time zone does cron = '<expr>' use, and which every = '<duration>' values are allowed? / How does rhiza health rate a schedule's part?

## Links

### In

- answers [scheduled-run-readonly-cron-health](../decision/scheduled-run-readonly-cron-health.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
