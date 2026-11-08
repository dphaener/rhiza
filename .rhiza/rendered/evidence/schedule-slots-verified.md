# schedule-slots-verified

evidence · `01M6WPAYYM3ATZB1AT42CDN72P` · actor [claude-code-298ba213](../actor/claude-code-298ba213.md)

## result

PASS

## statement

Verifier at 3f7ad85b ran go test -trimpath -v ./internal/schedule/...: all 8 tests PASS (TestCronFields: step */20, list 5,35, range 9-11, range-step 10-30/10, month and weekday names jan-feb mon-wed, 7 as Sunday, 29 FEB; TestCronMacros: @@hourly @@daily @@weekly @@monthly @@yearly; TestCronEitherDayField: '0 0 15 * fri' fires on the 15th or Fridays, a star field requires both; TestCronTimeZone: CRON_TZ=America/Denver 0 9 * * * against UTC clocks gives 16:00Z winter and 15:00Z summer; TestCronDST: 02:30 skipped on 2026-03-08 spring-forward, 01:30 fired once on 2026-11-01 fall-back and hourly slots in the repeated hour fire once; TestEvery: 15m gives only :00/:15/:30/:45 from four start times, slots are epoch multiples (90m from epoch = 5400s); TestParseErrors: six fields, CRON_TZ=Mars/Olympus, every 30s and every -5m all rejected). gofmt -l and go vet on the package are clean. Read the diff: pure package (no clock or file reads), time/tzdata embedded, Cron.Next walks month/day/hour/minute in the zone and skips a repeated wall reading via ZoneBounds, Every.Next floors to epoch multiples in UTC; only the builder's own actor record is added under .rhiza/log.

## Links

### Out

- verifies [schedule-slots](../capability/rhiza/workflows/scheduled-triggers/schedule-slots.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
