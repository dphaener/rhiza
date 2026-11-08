# daemon-status-queued-continuation-verified

evidence · `01M6XY0MVY0TP7XM04EX8DAZ4J` · actor [claude-code-c438e7e5](../actor/claude-code-c438e7e5.md)

## result

PASS

## statement

At 7e0959aa, with RHIZA_* unset and GOTMPDIR/TMPDIR outside the repo, ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-continue-slot|describe-golden|daemon-status.*|daemon-root.*)$': all PASS. daemon-continue-slot.txtar (fake rex) has one live running run, max-runs 1, and two interrupted runs queued with rhiza run resume. rhiza daemon status prints QUEUEDONE01 and QUEUEDTWO02 as 'interrupted, queued, waiting for a slot (place 1)/(place 2)', and --json gives queued 1 and 2 with state interrupted. An interrupted run with nothing queued (PLAININT003) shows plain 'interrupted' with no queued field. describe-golden passes with the new queued field (omitempty int) on both DaemonRun schemas. Read the diff: daemonStatusRuns numbers places over the project's continuation queue in order, counting only entries daemonContinueCanStart accepts (the same predicate the slot accounting uses), for both per-project and root-home status. go test ./internal/api -run Daemon: every test PASS. The package's root-state guard then flagged ~/code/.rhiza-root/usage.live.json, which the live root daemon writes; this change does not.

## Links

### Out

- verifies [daemon-status-queued-continuation](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-queued-continuation.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
