# schedules-file

capability · `01M6W5HTZ9ZZ6PWNAJ739T6ASH` · parent [scheduled-triggers](../scheduled-triggers.md)

## acceptance

Unit tests in the package that parses .rhiza/dispatch.toml cover: no file or no [[schedules]] giving no schedules; an entry with every and workflow and one with cron and command read with defaults catch_up once and overlap skip; a file with both [request] and [[schedules]] reading both with no error; an unknown table still failing with an error naming it; each error case (unknown key, duplicate name, every with cron, neither workflow nor command, bad cron, every 30s, catch_up = "all", overlap = "parallel") failing with an error naming the schedule and the key; and a testscript where rhiza check passes on a repo whose .rhiza/dispatch.toml holds [request] and [[schedules]]

## statement

Rhiza reads [[schedules]] entries from the repo's committed .rhiza/dispatch.toml beside its [request] and [ready] dispatch tables, widening the tables that file allows to those three (amending dispatch-rules-file's rule that any other table is an error, so [[schedules]] is no longer one): each entry has a name unique in the repo, exactly one of every and cron, exactly one of workflow and command, and optionally catch_up (once, the default, or skip) and overlap (skip, the default, or queue); any other key, a duplicate name, both or neither of a pair, an unparsable every or cron, or an unknown catch_up or overlap value is an error naming the schedule and key, and any table other than [request], [ready] and [[schedules]] is still an error naming it

## Links

### Out

- depends-on [schedule-slots](schedule-slots.md)
- decided-by [scheduled-triggers-shape](../../../../decision/scheduled-triggers-shape.md)
- decided-by [scheduled-run-readonly-cron-health](../../../../decision/scheduled-run-readonly-cron-health.md)

### In

- depends-on [schedules-check](schedules-check.md)
- depends-on [scheduled-run](scheduled-run.md)
- depends-on [daemon-fires-schedules](daemon-fires-schedules.md)
- depends-on [schedules-guide-topic](schedules-guide-topic.md)
- blocks [review-01m6w4k1-sched-r3-f2](../../../../question/review-01m6w4k1-sched-r3-f2.md)
- verifies [schedules-file-verified](../../../../evidence/schedules-file-verified.md)
- verifies [schedules-file-qa-1](../../../../evidence/schedules-file-qa-1.md)
- verifies [schedules-file-verified-2](../../../../evidence/schedules-file-verified-2.md)
- verifies [pr-374-scheduled-triggers](../../../../evidence/pr-374-scheduled-triggers.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
