# rebase-conflict-live-daemon-settings

question · `01M6D0ABE47ZN0RP4KZGW578SE`

## statement

Rebasing branch cs/live-daemon-settings onto origin/main conflicts outside the rendered files, so ship-change stopped. How should these conflicts be resolved: internal/api/daemon.go?

## Links

### In

- answers [live-settings-rebase-resolution](../decision/live-settings-rebase-resolution.md)

## History

- `01M6CRZQ0YKB2QV4V5TWY71G4V` The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy
