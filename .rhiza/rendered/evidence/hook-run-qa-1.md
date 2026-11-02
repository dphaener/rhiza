# hook-run-qa-1

evidence · `01M6DSV0BEFK34MRZJDSTZE29R` · actor [claude-code-dcf481da](../actor/claude-code-dcf481da.md)

## result

PASS

## statement

QA at 59869bfc with a scratch build in a fresh repo: hook run notify (git config command) got the sample inbox JSON on stdin and RHIZA_INBOX_KIND/ID/CHANGE/TITLE in env, printed its stdout and stderr then 'notify hook: exit status 0', exit 0; hook run event --command got the lifecycle JSON and RHIZA_EVENT=accepted; 'echo boom; exit 7' printed boom and exited 7 (--json code 7); unset event hook and unknown hook both exited 2 with next hints; api call hook.run gave the same response; with a live daemon, hook run notify used the daemon flag command and daemon.state.json was unchanged before and after.

## Links

### Out

- verifies [hook-run](../capability/rhiza/hooks/hook-commands/hook-run.md)

## History

- `01M6BVTB161RJF58ZDY12SFBH3` rhiza hook ls shows which hook commands are set and where from, and rhiza hook run tries one with a sample
