# schedule-due

capability · `01M6W5HVEXS3EA891B7V2FAT01` · parent [scheduled-triggers](../scheduled-triggers.md)

## acceptance

Unit tests in internal/daemon cover: a slot due on time; three missed slots after downtime firing once under catch_up once and none under skip; a new schedule not firing before its first slot; a slot during a running run dropped under overlap skip and fired once at its end under queue, with two such slots still firing once

## statement

A pure daemon decision says which schedules are due at a time, from each schedule's slots and its runs' journals (the last slot fired and whether its run is still going), so the state survives daemon restarts: one firing per due slot however many slots were missed, catch_up once firing once after downtime and skip waiting for the next slot, a schedule with no run yet waiting for its next slot, overlap skip dropping a slot while the schedule's previous run is going and queue keeping one firing to start when it ends

## Links

### Out

- depends-on [schedule-slots](schedule-slots.md)
- decided-by [scheduled-triggers-shape](../../../../decision/scheduled-triggers-shape.md)

### In

- depends-on [daemon-fires-schedules](daemon-fires-schedules.md)
- verifies [schedule-due-verified](../../../../evidence/schedule-due-verified.md)
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
