# health-schedules-part-qa-1

evidence · `01M6WX8XVJH69YD1T2J0PH6ZDE` · actor [claude-code-64d1c225](../actor/claude-code-64d1c225.md)

## result

PASS

## statement

At 82b0f3e2 I built rhiza into a scratch dir and used it on a scratch git repo registered in a scratch Rhiza root, with RHIZA_ROOT/CHANGE/RUN unset. rhiza health exited 0 every time, and its schedules part was always warn or ok, never bad. With no daemon it warned 'no daemon is running to fire them', with next 'rhiza daemon --autonomous'. Each schedule got a detail line with its last run time, run ID, result and next run, or 'none yet'. Three schedules that changed the checkout and one that failed were each named in the summary. With the daemon up and only failer failing, it warned 'failer failed', with next 'rhiza run log <that run>'. When a schedule earlier in the file had changed the checkout, next pointed at that schedule's run status, not at the failing run (reported as a minor finding).

## Links

### Out

- verifies [health-schedules-part](../capability/rhiza/workflows/scheduled-triggers/health-schedules-part.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
