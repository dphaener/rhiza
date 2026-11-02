# settings-in-git-config

decision · `01M6CS0H7S57S9433XGAN4T0Y1`

## statement

Daemon and autonomy settings live in the repo's git config and are re-read every poll: rhiza.daemon.paused, maxRuns, minFree and permissionMode, and rhiza.autonomy.plan, accept, ship and signoff, each auto or owner; a daemon flag given on the command line wins over config for that run; defaults keep today's behaviour (plan owner, accept owner except an autonomous plan's gate, ship auto, signoff auto when --autonomous); rhiza daemon settings lists every setting with value and source, and rhiza daemon set <key> <value> validates and writes one

## wording

(assumed under the owner's standing instruction to decide on their behalf)

## Links

### Out

- derived-from [m22-run-from-rex-ask](../source/m22-run-from-rex-ask.md)

### In

- decided-by [daemon-settings](../capability/rhiza/orchestration/daemon/daemon-settings.md)

## History

- `01M6CRZP2E8H6EEEPVDAWAKQ6H` Plan milestone 22: run Rhiza from Rex, the owner's ask of 2026-11-01
