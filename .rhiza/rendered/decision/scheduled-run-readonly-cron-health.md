# scheduled-run-readonly-cron-health

decision · `01M6W5013PBWWW9NTYV38FRH8R`

## reasoning

Read-only: the owner's rule warn-dont-safeguard (commit 79cc0be1) quotes 'i don't think we can reliably safeguard against users shooting themselves in the foot. best we can do is warn people.' and 'a schedule can run anything. who cares?', so read-only cannot mean enforcement; it can only be a stated default plus a plain warning when broken. The warning hooks into what exists: healthMain (internal/api/health.go) already rates a dirty, off-main or behind main checkout warn, and the memory/rule that the main tree stays clean for close-out is the concrete harm a schedule writing there causes, so detecting the change after the run and reporting it in the run record and health is the useful, non-blocking signal. Cron and every: five-field cron is the dialect owners already know (crontab, GitHub Actions schedules), and the macros cost nothing; local time matches crontab and what an owner on one VM expects, while CRON_TZ inside the expression gives an explicit zone without adding a schedule key, keeping scheduled-triggers-shape's 'Nothing else is overridable by a schedule'. Every uses time.ParseDuration like workflow.go's positiveDuration and daemon_settings.go, so durations read the same across Rhiza; a 1m floor matches minute-granular cron and keeps a typo like '15s' from hammering the daemon (refused by check as an invalid value, not a safeguard against a choice: sub-minute schedules are not a meaningful daemon slot). Epoch alignment makes slots deterministic, which the intent's catch-up and overlap defaults need to count missed slots after downtime. Health: health.go's parts are ok/warn/bad where bad exits 1, and every comparable part (inbox, requests, main, usage) uses warn for 'worth a look'; the intent says rhiza health lists each schedule with its last run, result and next run and that failures go to the logs, which fits detail lines plus a next to the log. Making a failed schedule bad would turn a flaky 15-minute intake into a failing health check, against warn-dont-safeguard's spirit; parse errors already surface through the check part, so the schedules part needs no bad state.

## statement

(1) 'Read-only by default' is a convention Rhiza states and warns about, never a fence: a scheduled run gets no sandbox, no changed permissions and no refused command, and a command or workflow runs in the main checkout as written. Rhiza documents in the guides that a schedule should leave the main checkout as it found it. After each scheduled run ends, Rhiza compares the main checkout with its state at the start (branch, HEAD, uncommitted tracked changes); if the run changed it, the run's record and log carry a warning naming what changed, the schedule's health detail says so, and the existing main part of rhiza health already warns on a dirty or off-main checkout. Nothing is reverted or blocked. (2) cron = '<expr>' takes the standard five-field cron dialect (minute hour day-of-month month day-of-week) with *, lists, ranges, steps and month/weekday names, plus the macros @@hourly, @@daily, @@weekly, @@monthly and @@yearly; there is no seconds or year field, and when both day fields are restricted a slot matches either (classic cron). It is read in the daemon machine's local time zone unless the expression starts with CRON_TZ=<IANA zone> (for example 'CRON_TZ=America/Denver 0 9 * * 1-5'); across a DST change a slot that does not exist is skipped and a slot that occurs twice fires once. every = '<duration>' takes a Go duration (as time.ParseDuration, like Rhiza's other duration settings) of at least 1m; its slots are aligned to multiples of the duration since the Unix epoch in UTC, so '15m' fires at :00, :15, :30 and :45 and the slots survive daemon restarts. rhiza check reports an unparsable expression, an unknown CRON_TZ zone, or an every below 1m or not positive as a dispatch-file error. (3) rhiza health gains one part, schedules, per project with schedules: its summary counts the schedules, and each schedule is a detail line with its last run (time and run ID), result and next run. The part is ok when no schedule's last run failed and none left the main checkout changed (a schedule that has not run yet is ok); it is warn when any schedule's last run failed, its last run changed the main checkout, or no daemon is running to fire it, with next pointing at the failing run's log. It is never bad: a failing schedule is worth a look, not an exit 1, and an invalid [[schedules]] entry is already reported by the rhiza check part. A repo with no schedules shows no schedules part.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-d0780667-2](../question/question-d0780667-2.md)
- derived-from [scheduled-triggers-plan-interview](../source/scheduled-triggers-plan-interview.md)

### In

- decided-by [schedule-slots](../capability/rhiza/workflows/scheduled-triggers/schedule-slots.md)
- decided-by [schedules-file](../capability/rhiza/workflows/scheduled-triggers/schedules-file.md)
- decided-by [schedules-check](../capability/rhiza/workflows/scheduled-triggers/schedules-check.md)
- decided-by [main-checkout-change-warning](../capability/rhiza/workflows/scheduled-triggers/main-checkout-change-warning.md)
- decided-by [health-schedules-part](../capability/rhiza/workflows/scheduled-triggers/health-schedules-part.md)
- decided-by [schedules-guide-topic](../capability/rhiza/workflows/scheduled-triggers/schedules-guide-topic.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
