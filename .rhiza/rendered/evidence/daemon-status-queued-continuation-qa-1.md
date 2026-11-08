# daemon-status-queued-continuation-qa-1

evidence · `01M6XYC0GAGN0SZHMD5NB3MN9T` · actor [claude-code-c038643d](../actor/claude-code-c038643d.md)

## result

PASS

## statement

QA at 728608f3 with the built binary, hands-on in a scratch two-project root (fake rex, RHIZA_* unset): in project b, one running run, two interrupted runs queued with rhiza run resume (QUEUEDONE01 then QUEUEDTWO02) and one interrupted run with nothing queued (PLAININT003); rhiza daemon --max-runs 1 --once, then status. Root rhiza daemon status under 'Project b:' and project b's own rhiza daemon status both printed 'QUEUEDONE01 interrupted, queued, waiting for a slot (place 1)', 'QUEUEDTWO02 ... (place 2)', 'PLAININT003 interrupted' and 'RUNNING0001 running'. Root status --json gave projects[b].runs with queued 1 and 2 on the two queued runs and no queued field on PLAININT003 or the running run. describe-golden and daemon-continue-slot testscripts PASS.

## Links

### Out

- verifies [daemon-status-queued-continuation](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-queued-continuation.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
