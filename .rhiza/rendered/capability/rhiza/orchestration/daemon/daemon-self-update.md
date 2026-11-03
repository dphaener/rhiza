# daemon-self-update

capability · `01M6FES89JR3Z5XA751AGCD3RV` · parent [daemon](../daemon.md)

## acceptance

A Go test with injected git and make: a moved origin/main fetches, fast-forwards and runs make update once; an unmoved one runs nothing; a pending start delays it to a later poll; a dirty checkout or a failing make update is logged and not retried until origin/main moves again; with selfUpdate off nothing runs; rhiza daemon settings lists rhiza.daemon.selfUpdate = false (default)

## statement

With the setting rhiza.daemon.selfUpdate on (off by default, shown and edited in the autonomy app), when origin/main has moved past the main checkout the daemon, at the end of a poll with no start pending, fast-forwards the main checkout and runs make update there, so its re-exec on a changed binary takes over; a dirty or diverged checkout or a failing make update is logged and skipped until main moves again

## Links

### Out

- decided-by [m24-update-from-main](../../../../decision/m24-update-from-main.md)

### In

- verifies [daemon-self-update-verified](../../../../evidence/daemon-self-update-verified.md)
- verifies [daemon-self-update-qa-1](../../../../evidence/daemon-self-update-qa-1.md)
- verifies [daemon-self-update-verified-2](../../../../evidence/daemon-self-update-verified-2.md)
- verifies [pr-199-the-daemon-pulls-main-and-runs](../../../../evidence/pr-199-the-daemon-pulls-main-and-runs.md)

## History

- `01M6FEN1BRQRHP9RTJKN2YPTH6` The daemon pulls main and runs make update when main moves, and health shows the version it runs
