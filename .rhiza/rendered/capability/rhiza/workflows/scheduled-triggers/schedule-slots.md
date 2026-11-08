# schedule-slots

capability · `01M6W5HTVXNGQK1CJ943THTJDB` · parent [scheduled-triggers](../scheduled-triggers.md)

## acceptance

Unit tests in the new schedule package cover: next slot for a step, a list, a range, names and each macro; both day fields restricted matching either; CRON_TZ=America/Denver against a UTC clock; a spring-forward slot skipped and a fall-back slot fired once; every 15m giving :00, :15, :30, :45 whatever the start time; and parse errors for a six-field expression, an unknown zone, every 30s and every -5m

## statement

A schedule package computes a schedule's slots: cron takes the five-field dialect (minute hour day-of-month month day-of-week) with *, lists, ranges, steps, month and weekday names and the macros @@hourly, @@daily, @@weekly, @@monthly and @@yearly, read in the machine's local time zone unless prefixed CRON_TZ=<IANA zone>, a slot matching either day field when both are restricted, a slot a DST change skips not firing and one it repeats firing once; every takes a Go duration of at least 1m whose slots are multiples of it since the Unix epoch in UTC

## Links

### Out

- decided-by [scheduled-triggers-shape](../../../../decision/scheduled-triggers-shape.md)
- decided-by [scheduled-run-readonly-cron-health](../../../../decision/scheduled-run-readonly-cron-health.md)

### In

- depends-on [schedules-file](schedules-file.md)
- depends-on [schedule-due](schedule-due.md)
- verifies [schedule-slots-verified](../../../../evidence/schedule-slots-verified.md)
- verifies [schedule-slots-qa-1](../../../../evidence/schedule-slots-qa-1.md)
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
