# schedules-guide-topic-verified-2

evidence · `01M6WYCGVY4W66J7WT0YVDH2D4` · actor [claude-code-65b80ed2](../actor/claude-code-65b80ed2.md)

## result

PASS

## statement

At 0e9166f6 ran go test -trimpath -count=1 ./internal/guide/...: ok. TestWorkflowsText checks the workflows topic's Schedules section names [[schedules]], CRON_TZ, catch_up, overlap, every, cron, main checkout, no change set, as it found it, rhiza request, rhiza health and daemon log; the package's size/cap tests pass too. Read the diff: the fix now says cron without CRON_TZ reads local time, matching schedule.ParseCron (nil location -> time.Local, as dispatch.go and daemon_schedules.go call it). The depends-on main-checkout-change-warning is gone per Decision guide-topic-drops-warning-dependency; only depends-on schedules-file remains.

## Links

### Out

- verifies [schedules-guide-topic](../capability/rhiza/workflows/scheduled-triggers/schedules-guide-topic.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
