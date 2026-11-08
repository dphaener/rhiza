# scheduled-triggers-shape

decision · `01M6W4VFJ8Z4Y2K49NCY21ASDZ`

## reasoning

Milestone: the sibling change set Workflow dispatch (01M6VY7G) serves milestone-25 by its arbiter's decision, Scheduling as policy landed beside it, and the intent names Nebari's issue intake as the example: owner-requested daemon work, which is milestone-25's statement; no other milestone fits. Location: internal/rootcfg's Config holds only projects, self_update_project and machine-wide daemon settings, is uncommitted and shared by several projects, so per-repo schedules there would leak across projects. Workflow dispatch already ruled the same way and moved its rules to .rhiza/dispatch.toml (commits 15cf9c5a, 907c325a), and the intent itself calls a schedule another trigger kind, so the dispatch file is where triggers live; a second file would split one concept. Where it runs: every run today is on a change set in its worktree (decision change-set-per-worktree), but a 15-minute schedule opening a draft change set per firing would flood the model, which the question names; the main checkout is the natural place for work with no change set, and keeping it read-only respects the rule that the main tree stays clean for close-out. Writing work through requests keeps the model's change-set discipline intact. Limits: the daemon already gates every start on maxRuns (internal/api/daemon.go), and a schedule bypassing it would let a fast schedule starve change-set work; warn-dont-safeguard argues against new safeguards but not for removing existing limits. The overrides are exactly the two defaults the intent names as overridable, given explicit keys so rhiza check can validate them.

## statement

Scheduled triggers serve @milestone-25 (the backlog milestone), like Workflow dispatch, which it depends on. (1) Schedules are declared as [[schedules]] in the committed .rhiza/dispatch.toml, beside the dispatch rules, not in the root's rhiza.toml; each entry has a name (unique in the repo), every = '<duration>' or cron = '<expr>' (exactly one), and workflow = '<name>' or command = '<anything>' (exactly one). rhiza check reports a bad entry (both or neither, unparsable duration or cron, unknown workflow, unknown key) like other dispatch-file errors. (2) A scheduled run has no change set and opens none: it runs on the repo's main checkout, read-only by default, records a run with the trigger 'schedule' and the schedule's name in place of a change set, and writes its logs like any run. Work a scheduled run wants done (for example issue intake) goes through the existing paths, such as queuing requests, never by writing to the model from the schedule itself. (3) Scheduled runs take slots under the same rhiza.daemon.maxRuns and per-project maxRuns as other runs and gain no override of them; when no slot is free the due slot waits for one and still counts as one firing. Per schedule, two defaults can be overridden in the entry: catch_up = 'once' (default; after downtime run once) or 'skip' (after downtime wait for the next slot), and overlap = 'skip' (default; skip a slot while the previous run of the same schedule is going) or 'queue' (run once when it ends). Nothing else is overridable by a schedule.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-d0780667-1](../question/question-d0780667-1.md)
- derived-from [scheduled-triggers-plan-interview](../source/scheduled-triggers-plan-interview.md)

### In

- decided-by [scheduled-triggers](../capability/rhiza/workflows/scheduled-triggers.md)
- decided-by [schedule-slots](../capability/rhiza/workflows/scheduled-triggers/schedule-slots.md)
- decided-by [schedules-file](../capability/rhiza/workflows/scheduled-triggers/schedules-file.md)
- decided-by [schedules-check](../capability/rhiza/workflows/scheduled-triggers/schedules-check.md)
- decided-by [scheduled-run](../capability/rhiza/workflows/scheduled-triggers/scheduled-run.md)
- decided-by [schedule-due](../capability/rhiza/workflows/scheduled-triggers/schedule-due.md)
- decided-by [daemon-fires-schedules](../capability/rhiza/workflows/scheduled-triggers/daemon-fires-schedules.md)
- decided-by [schedules-guide-topic](../capability/rhiza/workflows/scheduled-triggers/schedules-guide-topic.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
