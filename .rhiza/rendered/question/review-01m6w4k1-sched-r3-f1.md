# review-01m6w4k1-sched-r3-f1

question · `01M6W6KPTEX7QHZCPMD01P63EN`

## statement

Plan review round 3 left finding f1 (unreal-depends-on on schedules-guide-topic): schedules-guide-topic depends-on main-checkout-change-warning, but the guide only states the convention that a schedule should leave the main checkout as it found it. Its statement and acceptance never mention the warning, and writing the workflows topic needs none of that leaf's code. Its depends-on schedules-file is enough. The critic confirmed it: schedules-guide-topic's statement and acceptance document the convention that a scheduled run should leave the main checkout as it found it, and never mention the warning main-checkout-change-warning records. Writing that guide text needs none of that leaf's code (the comparison of branch, HEAD and tracked changes at run end). The depends-on on schedules-file covers the keys the topic documents, so the link to main-checkout-change-warning is unreal. How should it be settled?

## Links

### Out

- blocks [schedules-guide-topic](../capability/rhiza/workflows/scheduled-triggers/schedules-guide-topic.md)

### In

- answers [guide-topic-drops-warning-dependency](../decision/guide-topic-drops-warning-dependency.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
