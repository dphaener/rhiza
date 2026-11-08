# daemon-fires-schedules-qa-1

evidence · `01M6WX8YCX24HRCGGKPKPV2GRM` · actor [claude-code-64d1c225](../actor/claude-code-64d1c225.md)

## result

FAIL

## statement

At 82b0f3e2 I built rhiza into a scratch dir and used it on a scratch git repo registered in a scratch Rhiza root, with RHIZA_ROOT/CHANGE/RUN unset. I ran rhiza daemon --autonomous --max-runs 1 --min-free 0 --permission-mode manual for real, with a logging fake rex first on PATH. slow (every 1m, sleep 70) fired at 06:21. failer waited while slow held the only slot, then fired once for slot 06:22. slow's 06:22 slot came during its run and was dropped. A broken dispatch.toml logged 'schedules not read, none started'. workflow = "vanished" logged 'schedule nightly not started: no workflow "vanished"' and started nothing. The daemon log carried every start and every end, including 'run … of schedule failer failed: step "command" ended with outcome "fail"'. The fake rex CLI was never called. But the failed scheduled run also became an owner inbox item ('failed run …, change set , step command', next run resume --retry), and the daemon wrote an OSC 7501 state=blocked status reading 'a change set: a run failed' (reported as a major finding). Separately, with the default rhiza.daemon.trimGoCache each poll took about 80s on this VM.

## Links

### Out

- verifies [daemon-fires-schedules](../capability/rhiza/workflows/scheduled-triggers/daemon-fires-schedules.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
