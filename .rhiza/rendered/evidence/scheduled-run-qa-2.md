# scheduled-run-qa-2

evidence · `01M6WZ98Q1T7V2227XTJR5EN9Z` · actor [claude-code-170ee890](../actor/claude-code-170ee890.md)

## result

PASS

## statement

At 1738f418, the real daemon started scheduled command runs as 'rhiza run start --foreground --schedule <name>' on the scratch project's main checkout, with no change set. Six runs were journalled, and git status of the checkout stayed clean.

## Links

### Out

- verifies [scheduled-run](../capability/rhiza/workflows/scheduled-triggers/scheduled-run.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
