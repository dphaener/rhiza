# scheduled-triggers-plan-interview

source · `01M6W5GEV7HWZXNE5R20ZJNA6C`

## quotes

```json
[]
```

## raw_sha256

0a38b955868470e8f805d52bdb05cbe245081efed76c49a34a5f71956031e2dd

## session

d0780667-0b40-44e7-a17d-01a1b849f0e6

## summary

Autonomous plan run on change set 01M6W4K (Scheduled triggers). No owner took part: the planner raised two batches of questions and the arbiter answered each with a Decision worded (assumed, would have escalated). Facts from the code: rhiza.toml is the root's machine-wide config (internal/rootcfg), not committed in any repo and shared by several projects; workflow dispatch (01M6VY7G, accepted, not yet built) moved its rules to a committed .rhiza/dispatch.toml for that reason and serves milestone-25; every run today needs a change set and its worktree (RunStartRequest.Change, internal/api/run.go); rhiza health reports fixed parts (internal/api/health.go); the repo has no .rhiza/policy.toml, so backlog work builds after milestone work. Settled in scheduled-triggers-shape: the change set serves milestone-25; [[schedules]] live in .rhiza/dispatch.toml with name, every or cron, workflow or command; a scheduled run has no change set, runs on the main checkout, records trigger schedule and the schedule's name, and acts on the model only through existing paths such as queuing requests; scheduled runs share maxRuns, a due slot waits for a free slot and counts as one firing; per schedule only catch_up (once or skip) and overlap (skip or queue) are overridable. Settled in scheduled-run-readonly-cron-health: read-only is a documented convention, never a fence, with a warning when a run changed the main checkout's branch, HEAD or tracked files; cron is five-field with names, macros and CRON_TZ, local time by default, classic day-field OR and DST rules; every is a Go duration of at least 1m aligned to epoch multiples in UTC; health gains a schedules part per project with schedules, warn at worst. Also found: rhiza change show fails on this change set with 'not canonical' because of the known @-reference re-encoding bug, already tracked by draft change set 01M6VYTS (Fix change show HTML escape).

## title

Scheduled triggers plan, autonomous run with the arbiter

## Links

### In

- derived-from [scheduled-triggers-shape](../decision/scheduled-triggers-shape.md)
- derived-from [scheduled-run-readonly-cron-health](../decision/scheduled-run-readonly-cron-health.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
