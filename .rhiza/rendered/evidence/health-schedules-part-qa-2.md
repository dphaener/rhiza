# health-schedules-part-qa-2

evidence · `01M6WZ97TC6G0WJN77F0VA3Q0R` · actor [claude-code-170ee890](../actor/claude-code-170ee890.md)

## result

PASS

## statement

At 1738f418, against the same scratch root, rhiza health with no daemon warned '2 schedules; no daemon is running to fire them', with next run 06:56:00Z for each. With the daemon up it warned 'failer failed', listed each schedule's last run (time, run ID, failed/running/interrupted) and next run, and gave next 'rhiza run log <failed run>'. After the kill the sleeper line read 'interrupted'.

## Links

### Out

- verifies [health-schedules-part](../capability/rhiza/workflows/scheduled-triggers/health-schedules-part.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
