# daemon-service-env

capability · `01M6JYBE226RWZVG1N6CCTS81F` · parent [daemon](../daemon.md)

## acceptance

daemon-install.txtar runs install with a known PATH and TMPDIR (one holding a space) and checks the unit carries Environment=PATH=... and Environment=TMPDIR=... with those values, quoted the way systemd reads them, next to ExecStart ending in daemon --autonomous; a second install with the same environment writes nothing, and one with a changed PATH rewrites the unit and reloads systemd

## statement

rhiza daemon install writes the PATH and TMPDIR of the shell that runs it into the systemd unit as Environment= lines, quoted the way systemd reads Environment= (the whole NAME=value in double quotes, with each backslash and double quote backslash-escaped and % doubled, and $ left as is, unlike ExecStart's word quoting), so the service starts agents with the same tools on PATH and the same temp dir as the Rex tab, with no hand-written drop-in

## Links

### Out

- decided-by [daemon-service-env-matches-tab](../../../../decision/daemon-service-env-matches-tab.md)

### In

- verifies [daemon-service-env-verified](../../../../evidence/daemon-service-env-verified.md)
- verifies [daemon-service-env-qa-1](../../../../evidence/daemon-service-env-qa-1.md)
- verifies [pr-269-fix-the-daemon-behaves-the](../../../../evidence/pr-269-fix-the-daemon-behaves-the.md)

## History

- `01M6JY2NZESHV22Y931D7XMBEE` Fix The daemon behaves the same however its installation runs it: the systemd unit and the Rex tab start rhiza daemon with the same arguments, so choosing service or tab never changes whether runs go on without the owner (replaces draft change set 01M6FN21KA287N, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
