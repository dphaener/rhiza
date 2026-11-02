# hook-ls

capability · `01M6BZW3EF9VGT73M6T3FB655T` · parent [hook-commands](../hook-commands.md)

## acceptance

hook-commands.txtar: hook ls shows both unset; after git config rhiza.notifyCommand it shows notify from git config; with a recorded daemon state carrying --notify-command it shows daemon flag; daemon status --json shows the recorded commands; describe-golden.txtar is regenerated

## statement

rhiza hook ls lists the notify and event hooks, each with its command and source: daemon flag (read from a live daemon's state), git config, or unset

## Links

### Out

- decided-by [hook-ls-reads-daemon-state](../../../../decision/hook-ls-reads-daemon-state.md)

### In

- verifies [hook-ls-verified](../../../../evidence/hook-ls-verified.md)
- verifies [hook-ls-qa-1](../../../../evidence/hook-ls-qa-1.md)
- verifies [pr-196-hook-commands](../../../../evidence/pr-196-hook-commands.md)

## History

- `01M6BVTB161RJF58ZDY12SFBH3` rhiza hook ls shows which hook commands are set and where from, and rhiza hook run tries one with a sample
