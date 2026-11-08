# schedules-amend-dispatch-rules-here

decision · `01M6W6PW2F17SE6B31EYAWSCRQ`

## reasoning

Read the question, schedules-file (which already claims to amend dispatch-rules-file's unknown-table rule), and dispatch-rules-file in 01M6VY7G (its statement says 'any other table or key ... is an error'). Rewording only schedules-file would leave the two rules contradicting each other once both land, so dispatch-rules-file itself has to change. The two places that could do it: (1) edit 01M6VY7G. It is accepted and queued to build, its intent is workflow dispatch rather than schedules, and reopening another accepted plan from this one would churn work already approved. (2) Make the write here. Schedules are this change set's intent, it already depends on Workflow dispatch and builds after it, and it already claims the amendment. So (2) puts the write where the claim is. This change set's view cannot see dispatch-rules-file yet (rhiza show returns not_found without --change 01M6VY7G), so the write has to follow that change set's landing. The ordering this plan already follows covers that.

## statement

Scheduled triggers settles the contradiction itself: it adds a set on dispatch-rules-file's statement and acceptance so the tables .rhiza/dispatch.toml allows are [request], [ready] and [[schedules]] (whose entries schedules-file defines), and any other table is still an error naming it. Workflow dispatch (01M6VY7G) is left as it is. Because dispatch-rules-file only exists in Workflow dispatch's view until that change set lands, the write is made as soon as it lands (or through a --source on it), and before this change set is accepted; schedules-file keeps its wording, and that write now backs it.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m6w4k1-sched-r3-f2](../question/review-01m6w4k1-sched-r3-f2.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
