# hook-run

capability · `01M6BZW5XMG33S31H2BKT2JBZG` · parent [hook-commands](../hook-commands.md)

## acceptance

hook-commands.txtar: hook run notify with a fake command appending stdin and RHIZA_INBOX_KIND writes one line and exits 0; hook run event sets RHIZA_EVENT; a failing command exits non-zero with its output shown; an unset hook exits 2; the daemon state's handled items are unchanged

## statement

rhiza hook run <notify|event> [--command c] runs that hook's command once with a sample payload through the daemon's runner, showing its output and exit status, and records nothing

## Links

### Out

- decided-by [hook-ls-reads-daemon-state](../../../../decision/hook-ls-reads-daemon-state.md)

### In

- verifies [hook-run-verified](../../../../evidence/hook-run-verified.md)
- verifies [hook-run-qa-1](../../../../evidence/hook-run-qa-1.md)
- verifies [pr-196-hook-commands](../../../../evidence/pr-196-hook-commands.md)

## History

- `01M6BVTB161RJF58ZDY12SFBH3` rhiza hook ls shows which hook commands are set and where from, and rhiza hook run tries one with a sample
