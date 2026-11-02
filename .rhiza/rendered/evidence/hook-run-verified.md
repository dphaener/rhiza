# hook-run-verified

evidence · `01M6DS95FA9ERCGBBG7YCEEREK` · actor [claude-code-0ef932fa](../actor/claude-code-0ef932fa.md)

## result

PASS

## statement

At 317ce347, go test ./cmd/rhiza -run 'TestScripts/(hook-commands|describe-golden)' -count=1 passed and make check exited 0. hook-commands.txtar runs hook run notify with a git-config command that appends stdin and RHIZA_INBOX_KIND: exit 0, output shown, kind=question logged exactly once next to the sample JSON; hook run event --command sets RHIZA_EVENT=accepted; a command ending 'exit 7' exits non-zero with 'boom' and 'exit status 7' shown; an unset event hook fails with 'unset' on stderr; daemon.state.json compares byte-equal before and after. By hand with a scratch build in a fresh repo: unset hook exit code 2 (with a git config next hint), unknown hook exit 2, failing command exit 7, notify exit 0 with one invocation logged, and no daemon state file or git change written. Diff read: hook_run.go reuses runHookCommandOutput, the daemon's runner, and writes nothing.

## Links

### Out

- verifies [hook-run](../capability/rhiza/hooks/hook-commands/hook-run.md)

## History

- `01M6BVTB161RJF58ZDY12SFBH3` rhiza hook ls shows which hook commands are set and where from, and rhiza hook run tries one with a sample
