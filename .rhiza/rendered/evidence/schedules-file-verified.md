# schedules-file-verified

evidence · `01M6WPH8E8XDCKT8MXVY14N26Z` · actor [claude-code-f44c28b9](../actor/claude-code-f44c28b9.md)

## result

PASS

## statement

Verified at 0fdec5a0 by reading the diff (internal/dispatch/dispatch.go parseSchedules, dispatch_test.go, dispatch-check.txtar) and running `go test -trimpath -count=1 -v -run 'Schedules|Parse|Load' ./internal/dispatch/` (TestDefaultsWhenNoFile, TestSchedulesAbsent, TestSchedulesRead, TestSchedulesErrors, TestParseErrors all PASS) and `go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/dispatch-check$'` (PASS). Covered: missing file and empty or [request]-only content give no schedules; an every+workflow entry reads with defaults catch_up=once, overlap=skip and a cron+command entry reads (that entry sets catch_up=skip/overlap=queue explicitly, so non-default values are checked too; defaults are set by the same per-entry constructor); [request] beside [[schedules]] reads both; unknown table [other] errors naming it; unknown key, duplicate name, every with cron, neither workflow nor command (plus neither every nor cron, both workflow and command), bad cron, every 30s, catch_up=all, overlap=parallel each fail with an error naming schedule "s" and the key. The testscript shows rhiza check passing on a .rhiza/dispatch.toml holding [request] and two [[schedules]], and failing with 'schedule "intake": every' on a 30s entry.

## Links

### Out

- verifies [schedules-file](../capability/rhiza/workflows/scheduled-triggers/schedules-file.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
