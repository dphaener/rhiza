# hook-ls-verified

evidence · `01M6DS523QG6YA481Z1HY9YD2G` · actor [claude-code-e885070a](../actor/claude-code-e885070a.md)

## result

PASS

## statement

At 2bf7d5e0, ran go test ./cmd/rhiza -run 'TestScripts/(hook-commands|describe-golden)$' -count=1: both PASS. hook-commands.txtar checks hook ls shows notify and event unset (text and --json source unset); after git config rhiza.notifyCommand it shows 'echo cfg (git config)'; with a live daemon started with --notify-command (state recorded in daemon.state.json) it shows 'echo flag (daemon flag)' with event still from git config; daemon status --json shows name/command/source of the recorded hooks and text status shows 'Hook notify: echo flag (daemon flag)'; after the daemon stops, config answers again. Re-ran describe-golden with -update: no diff, so the golden is current. make check (fmt-check, vet, test, self-check) exited 0 with a clean tree. Read the diff: internal/api/hook.go resolves flag > git config > unset, hook ls prefers a live daemon's recorded state.Hooks only when the state's process matches the lock holder; daemon start records hooks in state; daemon status lists them.

## Links

### Out

- verifies [hook-ls](../capability/rhiza/hooks/hook-commands/hook-ls.md)

## History

- `01M6BVTB161RJF58ZDY12SFBH3` rhiza hook ls shows which hook commands are set and where from, and rhiza hook run tries one with a sample
