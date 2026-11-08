# schedules-check-verified

evidence · `01M6WRWQRKEW15PMCA66Q2FN1P` · actor [claude-code-e2eed522](../actor/claude-code-e2eed522.md)

## result

PASS

## statement

At 421532c6 ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(schedules-check|schedules-file|dispatch)' -count=1: ok. The schedules-check script passes rhiza check with a schedule naming a repo workflow (.rhiza/workflows/bugfix.toml) and one with a command, and fails it on a missing workflow, an unknown key and a bad cron, each stdout naming the schedule. By hand with a scratch build, rhiza check exits 1 with: 'dispatch.toml: schedule "x": workflow: names workflow "nothing-here", which is neither in .rhiza/workflows nor built in', '[[schedules]] schedule "x": unknown key "evry"', '[[schedules]] schedule "x": cron: cron "99 * * * *": minute: 99 out of range 0-59', each naming .rhiza/dispatch.toml, the schedule and the key. go test ./internal/api -run 'Check|Dispatch|Schedule' passed (its only FAIL was the root-state guard on ~/code/.rhiza-root/usage.live.json, unrelated); internal/dispatch ok. Diff read: dispatchProblems now checks each schedule's workflow with repoOrBuiltinWorkflow, the same check the dispatch rules use.

## Links

### Out

- verifies [schedules-check](../capability/rhiza/workflows/scheduled-triggers/schedules-check.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
