# main-checkout-change-warning-qa-1

evidence · `01M6WX8X032H0192FYEWXYND2B` · actor [claude-code-64d1c225](../actor/claude-code-64d1c225.md)

## result

PASS

## statement

At 82b0f3e2 I built rhiza into a scratch dir and used it on a scratch git repo registered in a scratch Rhiza root, with RHIZA_ROOT/CHANGE/RUN unset. Each scheduled run ended in state succeeded. git checkout -b feature warned 'switched the main checkout's branch from main to feature'. A commit warned 'moved the main checkout's HEAD from b387179… to ae52673…'. Appending to tracked.txt warned 'changed uncommitted tracked files in the main checkout: tracked.txt'. Each warning was on stderr and in run start/status --json warnings. Nothing was reverted: the branch, the commit and the edit all stayed. A tick run in a checkout that was already dirty recorded no warning.

## Links

### Out

- verifies [main-checkout-change-warning](../capability/rhiza/workflows/scheduled-triggers/main-checkout-change-warning.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
