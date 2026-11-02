# live-settings-rebase-resolution

decision · `01M6D0ZHQXWR9XY1KYEDD56RTD`

## statement

The rebase conflict in runDaemon's start state is resolved by merging origin/main: daemonStartState from main is used with the live settings' max runs, and the settings rows are recorded before the state is first written

## wording

(arbiter, the orchestrator, under the owner's standing instruction to decide on their behalf)

## Links

### Out

- answers [rebase-conflict-live-daemon-settings](../question/rebase-conflict-live-daemon-settings.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
