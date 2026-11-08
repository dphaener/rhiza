# guide-topic-drops-warning-dependency

decision · `01M6W6MW114WEMH4Q0B4Y53AN4`

## reasoning

schedules-guide-topic's statement and acceptance document the convention that a scheduled run should leave the main checkout as it found it, plus the schedule keys, CRON_TZ, catch_up, overlap and where health and logs show; neither names the end-of-run warning that main-checkout-change-warning builds. Writing the guide text needs none of that leaf's code (the branch/HEAD/tracked-change comparison), and the keys it documents come from schedules-file, which it already depends on. The plan reviewer and the critic both found the link unreal. Keeping it would only serialize the build without reason. The other option, widening the topic to describe the warning, would add scope that the change set's intent and the owner never asked for, so dropping the link is the smaller and better-supported fix.

## statement

Remove the depends-on link from @schedules-guide-topic to @main-checkout-change-warning; the topic keeps its depends-on @schedules-file, and its statement and acceptance stay as they are

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6w4k1-sched-r3-f1](../question/review-01m6w4k1-sched-r3-f1.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
