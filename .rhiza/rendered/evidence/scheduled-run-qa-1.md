# scheduled-run-qa-1

evidence · `01M6WX8WPS9CJAACQD6H6YGG4E` · actor [claude-code-64d1c225](../actor/claude-code-64d1c225.md)

## result

PASS

## statement

At 82b0f3e2 I built rhiza into a scratch dir and used it on a scratch git repo registered in a scratch Rhiza root, with RHIZA_ROOT/CHANGE/RUN unset. rhiza run start --schedule tick (a command) ran in the main checkout: no worktree was added, no change set was opened, and its output was printed. run start --json, run ls --json and run status --json all show change "", trigger schedule and trigger_name tick, and run ls says 'schedule tick'. A repo workflow copied from the built-in build ran under schedule leafy and failed with 'step "build" is a leaves step, which needs a change set, and this run has none'. For daemon-started runs, rhiza run log shows the command's output ('slow start', 'slow done'). Human run status prints 'Change set: : ' and leaves out the schedule name and the warnings (reported as a minor finding).

## Links

### Out

- verifies [scheduled-run](../capability/rhiza/workflows/scheduled-triggers/scheduled-run.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
