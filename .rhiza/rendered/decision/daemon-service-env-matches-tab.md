# daemon-service-env-matches-tab

decision · `01M6JY7KFRW3AH998GN8GP7AEF`

## reasoning

I read the question, the change set's slice and intent, and internal/api/daemon_service.go and daemon_startstop.go. daemonArgs [daemon --autonomous] already feeds both ExecStart and the tab (since #200), and both scripts grep for it, so building the title as worded would be an empty change. The intent is that choosing service or tab never changes whether runs go on without the owner. serviceUnit writes only WorkingDirectory, ExecStart and Restart, so a fresh install runs under systemd's bare PATH, where agents cannot find claude, go, gh or make, and runs then fail without the owner. That is the one real way service and tab differ today, and the VM only works because of a hand-written env.conf drop-in. Capturing the installing shell's PATH and TMPDIR is the smallest change that closes it, matches what the tab inherits, and is testable in testscript without systemd. Leaving the existing drop-in alone is safe: systemd drop-ins override the unit's Environment= lines.

## statement

The change set's leaves cover the environment, not the arguments, which already match. rhiza daemon install writes the PATH and TMPDIR of the shell that runs it into the unit as Environment= lines, quoted the way systemd reads Environment= (the whole NAME=value in double quotes, with each backslash and double quote backslash-escaped and % doubled, and $ left as is, unlike ExecStart's word quoting), so the systemd service starts agents with the same tools on PATH and the same temp dir as the Rex tab and a fresh install needs no hand-written drop-in. daemon-install.txtar checks the unit carries both lines next to ExecStart with 'daemon --autonomous'. The shared daemonArgs and its existing greps stay as they are and get no leaf of their own.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-db198e1f-1](../question/question-db198e1f-1.md)
- derived-from [plan-daemon-service-env](../source/plan-daemon-service-env.md)

### In

- decided-by [daemon-service-env](../capability/rhiza/orchestration/daemon/daemon-service-env.md)

## History

- `01M6JY2NZESHV22Y931D7XMBEE` Fix The daemon behaves the same however its installation runs it: the systemd unit and the Rex tab start rhiza daemon with the same arguments, so choosing service or tab never changes whether runs go on without the owner (replaces draft change set 01M6FN21KA287N, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
