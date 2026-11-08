# schedule-slots-qa-1

evidence · `01M6WX8VRJMMHXT14D0MR55VQV` · actor [claude-code-64d1c225](../actor/claude-code-64d1c225.md)

## result

PASS

## statement

At 82b0f3e2 I built rhiza into a scratch dir and used it on a scratch git repo registered in a scratch Rhiza root, with RHIZA_ROOT/CHANGE/RUN unset. rhiza check reported 'cron "0 0 * * * *": want 5 fields', 'every "30s": must be at least 1m0s' and 'unknown time zone "Mars/Base"'; it accepted CRON_TZ=America/Denver 0 2 * * *. rhiza health showed next runs on the right boundaries at 06:13Z: every 15m gave 06:15:00Z, every 1m gave 06:14:00Z, @@daily gave 2026-11-09T00:00:00Z, @@hourly 07:00:00Z, @@weekly 2026-11-15T00:00:00Z (a Sunday) and 0 3 * * * gave 2026-11-09T03:00:00Z.

## Links

### Out

- verifies [schedule-slots](../capability/rhiza/workflows/scheduled-triggers/schedule-slots.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
