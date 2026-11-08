# daemon-fires-schedules-verified-2

evidence · `01M6WYT02NAA45F11NW3RY56Z3` · actor [claude-code-6201f631](../actor/claude-code-6201f631.md)

## result

PASS

## statement

Verified attempt 2 at a14eb564 (fix on bb9a5544). I read the diff: an interrupted scheduled run's LastEnd is now its last journal entry, fixed across polls, not each poll's now, and inbox.go leaves scheduled runs with no change set out of the inbox and its newer-run key, so no Rex blocked notice. I ran go test -trimpath ./internal/api -run 'TestDaemon(Poll|Schedule|Scheduled|Interrupted)|TestScheduledRun|TestHealthSchedules|TestInbox' -count=1 from an archive of HEAD outside ~/code: ok. From the checkout every test passes, but rootguard fails on ~/code/.rhiza-root/usage.live.json changing, a known environment leak that is not this leaf's code. internal/daemon: ok. Acceptance, with a fake clock: TestDaemonPollFiresSchedules (each due schedule starts across two projects; with maxRuns full the other waits, then fires once over repeated polls, no rex calls), TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow (starts nothing, log names the schedule), TestDaemonScheduledRunFailureLogged (daemon log has the start and the failure, zero rex calls, now also no inbox item and no Rex notice), and the new TestDaemonInterruptedScheduledRunNextSlotFires (the slot after an interrupted run fires). Review findings quality f1 and qa f1/f4 are fixed. quality f2 (catch_up=skip drops a slot that waits more than 2m) is unchanged and was minor. Process note: the fix commit's Actor claude-code-7a69f316 has no tracked record, though rhiza check reports nothing for it. make self-check's only failures are stale rendered files, as expected partway through a run.

## Links

### Out

- verifies [daemon-fires-schedules](../capability/rhiza/workflows/scheduled-triggers/daemon-fires-schedules.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
