# schedule-due-verified

evidence · `01M6WPMFD3334VACQP25Y0KVP7` · actor [claude-code-998ee67f](../actor/claude-code-998ee67f.md)

## result

PASS

## statement

At b0b8e269 I ran go test -trimpath -count=1 -run TestScheduleDue -v ./internal/daemon (4 tests PASS) and the whole internal/daemon package (ok); go vet and gofmt clean. Read internal/daemon/scheduledue.go and its tests: ScheduleDue is pure over dispatch.Schedule, schedule.Schedule, a journal-derived ScheduleState and now. Tests cover each acceptance point: TestScheduleDueOnTime (10:00 slot fires at 10:00, not at 9:59); TestScheduleDueMissedSlots (slots 2,3,4 missed fire once for 4:00 under catch_up once, a follow-up decision is quiet, none under catch_up skip with next 5:00); TestScheduleDueNewSchedule (Since 9:30 no fire at 9:45, fires 10:00); TestScheduleDueOverlap (no fire while running; after the run ends at 11:20 the 10:00 and 11:00 slots are dropped under overlap skip and fire once for 11:00 under queue, quiet after that run is journaled).

## Links

### Out

- verifies [schedule-due](../capability/rhiza/workflows/scheduled-triggers/schedule-due.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
