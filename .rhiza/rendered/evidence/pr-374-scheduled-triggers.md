# pr-374-scheduled-triggers

evidence · `01M6X1YBPN0GAYA6EM5B30KB5P` · actor [dphaener](../actor/dphaener.md)

## head

c3d67419598f2c48fe2079ec67841295aee5e751

## pr

https://github.com/dphaener/rhiza/pull/374

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #374 (https://github.com/dphaener/rhiza/pull/374) passed CI at c3d6741

## status

verified

## Links

### Out

- verifies [scheduled-triggers](../capability/rhiza/workflows/scheduled-triggers.md)
- verifies [schedule-slots](../capability/rhiza/workflows/scheduled-triggers/schedule-slots.md)
- verifies [schedules-file](../capability/rhiza/workflows/scheduled-triggers/schedules-file.md)
- verifies [schedules-check](../capability/rhiza/workflows/scheduled-triggers/schedules-check.md)
- verifies [scheduled-run](../capability/rhiza/workflows/scheduled-triggers/scheduled-run.md)
- verifies [main-checkout-change-warning](../capability/rhiza/workflows/scheduled-triggers/main-checkout-change-warning.md)
- verifies [schedule-due](../capability/rhiza/workflows/scheduled-triggers/schedule-due.md)
- verifies [daemon-fires-schedules](../capability/rhiza/workflows/scheduled-triggers/daemon-fires-schedules.md)
- verifies [health-schedules-part](../capability/rhiza/workflows/scheduled-triggers/health-schedules-part.md)
- verifies [schedules-guide-topic](../capability/rhiza/workflows/scheduled-triggers/schedules-guide-topic.md)
- verifies [dispatch-rules-file](../capability/rhiza/workflows/workflow-dispatch/dispatch-rules-file.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
