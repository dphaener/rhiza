# main-checkout-change-warning-verified

evidence · `01M6WTE8VMCZGBWCAWGMP9Z1W5` · actor [claude-code-d01c8606](../actor/claude-code-d01c8606.md)

## result

PASS

## statement

At 089435ff, go test -trimpath ./internal/api -run TestScheduled -v passed all four scheduled tests, including TestScheduledRunWarnsOnMainCheckoutChange with subtests quiet, branch, commit and edit. quiet records no warning and writes nothing to the log. branch, commit and edit each end in state succeeded (the command's own status) with one warning in run.status Warnings and on stderr, naming 'branch from main to feature', 'HEAD from aaaa1111 to bbbb2222' and 'main.go' respectively. go test ./internal/run/... and TestScripts/describe-golden pass, and gofmt and vet are clean. From the diff: runScheduledStart snapshots the branch, HEAD and the hashes of uncommitted tracked files before and after runEngine; a change appends a journal entry of type warning, which State folds into Warnings and runInfo exposes. The run's state is never touched, and nothing is reverted or blocked.

## Links

### Out

- verifies [main-checkout-change-warning](../capability/rhiza/workflows/scheduled-triggers/main-checkout-change-warning.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
