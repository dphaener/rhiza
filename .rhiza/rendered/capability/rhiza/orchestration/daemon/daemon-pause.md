# daemon-pause

capability · `01M6CS1Y195R6RZEQ6BJEFC1WB` · parent [daemon](../daemon.md)

## acceptance

daemon-pause.txtar: with paused set, daemon --once on a ready change set starts nothing and daemon status shows Paused: yes and each candidate skipped as paused; after setting it false the next --once starts the run

## statement

With rhiza.daemon.paused true the daemon starts and resumes no runs but keeps polling, reporting and notifying, and daemon status says paused

## Links

### Out

- decided-by [stage-autonomy-meanings](../../../../decision/stage-autonomy-meanings.md)

### In

- verifies [daemon-pause-verified](../../../../evidence/daemon-pause-verified.md)
- verifies [daemon-pause-qa-1](../../../../evidence/daemon-pause-qa-1.md)
- verifies [pr-187-live-daemon-settings](../../../../evidence/pr-187-live-daemon-settings.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
