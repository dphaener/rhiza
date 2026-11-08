# schedules-check-qa-1

evidence · `01M6WX8WCVB7P5NGE90N47JYVX` · actor [claude-code-64d1c225](../actor/claude-code-64d1c225.md)

## result

PASS

## statement

At 82b0f3e2 I built rhiza into a scratch dir and used it on a scratch git repo registered in a scratch Rhiza root, with RHIZA_ROOT/CHANGE/RUN unset. rhiza check passed on a schedule naming the built-in ship-change, one naming a repo workflow (.rhiza/workflows/lw.toml) and a command schedule. It failed on workflow = "missing-wf" with 'schedule "nightly": workflow: names workflow "missing-wf", which is neither in .rhiza/workflows nor built in'.

## Links

### Out

- verifies [schedules-check](../capability/rhiza/workflows/scheduled-triggers/schedules-check.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
