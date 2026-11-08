# daemon-fires-schedules-qa-2

evidence · `01M6WZ97KQ6MHD578R5QXE0XRP` · actor [claude-code-170ee890](../actor/claude-code-170ee890.md)

## result

PASS

## statement

At 1738f418 I built rhiza into a scratch dir and ran it on a scratch git repo registered in a scratch Rhiza root, with RHIZA_ROOT/CHANGE/RUN unset and a logging fake rex first on PATH. dispatch.toml had failer (every 1m, command exits 3) and sleeper (every 1m, sleep 600). I ran rhiza daemon --autonomous --max-runs 3 --min-free 0 --permission-mode manual for real. Both fired at 06:56. The daemon log said 'run … of schedule failer failed: step "command" ended with outcome "fail"'. rhiza inbox from the root and from the project said 'Nothing waits on you'. The daemon output had no escape bytes, so no OSC 7501 status. The fake rex was never called. I then SIGKILLed the daemon and the live sleeper run, and health showed sleeper 'interrupted'. I restarted the daemon at 06:58:16 and sleeper fired again for slot 06:58, so an interrupted run no longer stalls its schedule. The 06:59 sleeper slot was skipped while that run was live, and failer fired and failed again with the inbox still empty.

## Links

### Out

- verifies [daemon-fires-schedules](../capability/rhiza/workflows/scheduled-triggers/daemon-fires-schedules.md)

## History

- `01M6W4K49SDWJPW9NRF27QVBP3` Scheduled triggers: a repo's rhiza.toml declares [[schedules]] with every = '<duration>' or cron = '<expr>' and either workflow = '<name>' or command = '<anything>'; the daemon serving the root runs every repo's schedules, wrapping a command as an implicit one-step workflow so it gets the same run record and logs. Defaults a schedule can override: after downtime it runs once, not once per missed slot, and a slot is skipped while the previous scheduled run is still going. rhiza health lists each schedule with its last run, result and next run; failures go to the logs, never Rex tabs. Depends on the workflow dispatch request 01M6VY7CAJPC (a schedule is another trigger kind). Example use: Nebari's issue intake every 15m. Follows decisions workflows-are-the-rules and warn-dont-safeguard (PR #353).
