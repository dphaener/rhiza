# schedules-guide-topic-qa-2

evidence · `01M6WZ98D54TC9BFVDNJ08YPD5` · actor [claude-code-170ee890](../actor/claude-code-170ee890.md)

## result

PASS

## statement

At 1738f418, the built rhiza's guide workflows topic has a Schedules section naming every, cron (five fields, local time or a leading CRON_TZ=<zone>), workflow or command, catch_up and overlap. It says a scheduled run has no change set and that health shows schedules. The cycle 1 UTC wording is gone.

## Links

### Out

- verifies [schedules-guide-topic](../capability/rhiza/workflows/scheduled-triggers/schedules-guide-topic.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
