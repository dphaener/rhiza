# review-01m6w4k1-sched-r3-f2

question · `01M6W6KPX1HP0R9B3DRGWYQGJ9`

## statement

Plan review round 3 left finding f2 (contradiction on schedules-file): schedules-file allows [[schedules]] in .rhiza/dispatch.toml and says it amends dispatch-rules-file. But dispatch-rules-file (01M6VYTACHSR, in the accepted, unbuilt change set 01M6VY7G6HPJ, Workflow dispatch) still says 'any other table or key ... is an error', and this plan writes nothing to it. Once both land, the model holds two statements that contradict each other. rhiza overlaps cannot catch this because no node is written twice. dispatch-rules-file needs a change, either in 01M6VY7G or in a write this plan makes, so it names [[schedules]] as an allowed table. The critic confirmed it: schedules-file's statement says it amends dispatch-rules-file's unknown-table rule. But rhiza change show lists no write by this plan to dispatch-rules-file (01M6VYTACHSRCG9FCR51ESF67H in 01M6VY7G), whose statement still says 'any other table or key ... is an error naming it'. rhiza overlaps reports nothing because no node is written twice. Once both change sets land, the model holds two contradicting rules for .rhiza/dispatch.toml. dispatch-rules-file's statement needs to name [[schedules]] as an allowed table. How should it be settled?

## Links

### Out

- blocks [schedules-file](../capability/rhiza/workflows/scheduled-triggers/schedules-file.md)

### In

- answers [schedules-amend-dispatch-rules-here](../decision/schedules-amend-dispatch-rules-here.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
