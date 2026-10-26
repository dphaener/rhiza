# launch-env-stamping

capability · `01M5TYCYNEBVFGQ1B2YH1KDDQP` · parent [actor-stamping](../actor-stamping.md)

## acceptance

A test reads the launched environment and finds all three.

## statement

Every agent Genesis launches gets `GENESIS_CHANGE`, `GENESIS_ACTOR` and, when it belongs to a run, `GENESIS_RUN` in its environment.

## Links

### Out

- decided-by [actor-set-at-launch](../../../../decision/actor-set-at-launch.md)

### In

- verifies [launch-env-stamping-test](../../../../evidence/launch-env-stamping-test.md)

## History

- `01M5SZHY68YGJZETXGP786ASR7` The Claude Code adapter: Genesis launches Claude Code in a Rex tab and its own worktree, with the permission mode always explicit, its session ID and environment stamped, and its status read from hooks, then Rex events, then the process exiting.
