# hook-ls-qa-1

evidence · `01M6DSTYNRKQSAM265D50P5YXW` · actor [claude-code-dcf481da](../actor/claude-code-dcf481da.md)

## result

PASS

## statement

QA at 59869bfc with a scratch build (go build -o $TMPDIR/rhiza-x) in a fresh git repo, fake rex on PATH: hook ls printed 'notify: unset' / 'event: unset' and --json source unset; after git config rhiza.notifyCommand it showed '(git config)'; with a live 'rhiza daemon --permission-mode manual --notify-command "echo flag-notify"' and rhiza.eventCommand in config it showed 'notify: echo flag-notify (daemon flag)' and 'event: echo cfg-ev (git config)', and kept the daemon's value after git config changed under it; daemon status showed 'Hook notify: ...' lines and --json hooks with name/command/source; after SIGINT stop hook ls fell back to git config. go test hook-commands and describe-golden passed; make check exited 0.

## Links

### Out

- verifies [hook-ls](../capability/rhiza/hooks/hook-commands/hook-ls.md)

## History

- `01M6BVTB161RJF58ZDY12SFBH3` rhiza hook ls shows which hook commands are set and where from, and rhiza hook run tries one with a sample
