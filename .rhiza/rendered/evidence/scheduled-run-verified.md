# scheduled-run-verified

evidence · `01M6WS4PMEJQNP25X4KQXB7MKK` · actor [claude-code-f515264b](../actor/claude-code-f515264b.md)

## result

PASS

## statement

At 566ff696, read the diff and ran go test -trimpath -count=1 ./internal/api -run Scheduled and ./internal/run -run TestRunWithoutChangeSet: all pass. In a scratch repo, TestScheduledRunCommandAndWorkflow starts the command schedule intake and the repo-workflow schedule nightly (workflow tidy) with run start --schedule. Both succeed and are autonomous. Neither has a change set, worktree or branch in the started entry, and neither adds an entry to .rhiza/log. Both steps run with Dir set to the main checkout. The journal records trigger schedule and the schedule's name. run ls --json and run status --json show trigger and trigger_name with an empty change, and the run ls text says 'schedule intake'. The command's faked output reaches the run's stderr log. TestScheduledRunLeavesStepFails, on a schedule naming a copy of the built-in build workflow, and TestRunWithoutChangeSetFailsLeavesStep both fail the run with 'step "build" is a leaves step, which needs a change set', and no leaf starts. Unknown schedules and --workflow are refused. go vet and gofmt are clean, and the describe-golden testscript passes. In the full ./internal/api suite, 501 tests passed. Only TestReadBudget failed: the roadmap and ready read timings went over 50ms at VM load 25, a read path this diff does not touch. The suite then hit its timeout during TestLinkEndpointKinds.

## Links

### Out

- verifies [scheduled-run](../capability/rhiza/workflows/scheduled-triggers/scheduled-run.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
