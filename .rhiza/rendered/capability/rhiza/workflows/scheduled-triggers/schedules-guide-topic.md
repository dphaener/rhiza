# schedules-guide-topic

capability · `01M6W5HWFGG153SAAWP0XEK7JT` · parent [scheduled-triggers](../scheduled-triggers.md)

## acceptance

A test of the guide shows the workflows topic naming [[schedules]], CRON_TZ, catch_up, overlap and the main checkout convention, and the guide's existing size checks still pass

## statement

The workflows guide topic explains [[schedules]] in .rhiza/dispatch.toml: the keys, cron and every forms with CRON_TZ, catch_up and overlap, that a scheduled run has no change set and runs on the main checkout, that it should leave the checkout as it found it and acts on the model only through existing paths such as rhiza request, and where its health and logs show

## Links

### Out

- depends-on [schedules-file](schedules-file.md)
- decided-by [scheduled-triggers-shape](../../../../decision/scheduled-triggers-shape.md)
- decided-by [scheduled-run-readonly-cron-health](../../../../decision/scheduled-run-readonly-cron-health.md)

### In

- blocks [review-01m6w4k1-sched-r3-f1](../../../../question/review-01m6w4k1-sched-r3-f1.md)
- verifies [schedules-guide-topic-verified](../../../../evidence/schedules-guide-topic-verified.md)
- verifies [schedules-guide-topic-verified-2](../../../../evidence/schedules-guide-topic-verified-2.md)
- verifies [schedules-guide-topic-qa-2](../../../../evidence/schedules-guide-topic-qa-2.md)
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
