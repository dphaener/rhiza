# daemon-fires-schedules-verified

evidence · `01M6WVK2G3X8XRQCK3B59NAN84` · actor [claude-code-71d35ce8](../actor/claude-code-71d35ce8.md)

## result

PASS

## statement

At bb9a5544 I read the diff (internal/api/daemon_schedules.go, daemon_projects.go, daemon.go) and ran go vet plus go test -trimpath -count=1 ./internal/api -run 'Schedul|TestDaemon'. All pass on a fake clock with a fake Detach and fake exec. TestDaemonPollStartsEachProjectsSchedule: one poll over two projects starts each project's due schedule. TestDaemonPollFiresSchedules: with maxRuns=1, one schedule starts and the other waits while the first run is journaled and live; when that run ends, the waiting schedule fires exactly once over two polls, makes no rex call, and logs 'started a run of schedule'. TestDaemonScheduleUnresolvedWorkflow: workflow 'gone' starts nothing, and 'schedule nightly not started' is logged once. TestDaemonScheduledRunFailureLogged: the detached run start --schedule of a 'false' command fails with no rex call, and the daemon log carries both the start and 'of schedule bad failed'. Scheduled starts go through daemonDetachNoRex, so they never reach ensureRepoSession or Rex. Inside the ~/code root, the package's root-state guard sometimes flagged .rhiza-root/usage.live.json, which live status lines rewrite; no test failed, and the same run passes cleanly in a clone outside the root.

## Links

### Out

- verifies [daemon-fires-schedules](../capability/rhiza/workflows/scheduled-triggers/daemon-fires-schedules.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
