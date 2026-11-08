# pr-393-fix-a-queued-continuation-can

evidence · `01M6YE5ZBC9QRP90CT93DWFQJH` · actor [dphaener](../actor/dphaener.md)

## head

f80a4380f44a18821f24771f909305cda847f7ae

## pr

https://github.com/dphaener/rhiza/pull/393

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #393 (https://github.com/dphaener/rhiza/pull/393) passed CI at f80a438

## status

verified

## Links

### Out

- verifies [root-continuations-first](../capability/rhiza/orchestration/daemon/root-continuations-first.md)
- verifies [daemon-status-queued-continuation](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-queued-continuation.md)
- verifies [root-daemon-round-robin](../capability/rhiza/orchestration/daemon/root-daemon-round-robin.md)

## History

- `01M6XC4M3W7Y351WGXE8DTZACT` Fix A queued continuation can wait far longer than its slot should take: run 01M6WMX236BTT04TYQDACGKBZ4 (autonomy per workflow, PR #376 green) was queued with rhiza run resume at 2026-11-08T08:09Z and was still not started at 09:24Z; daemon status listed it as 'interrupted' among the live runs next to 3 running ones at rhiza.daemon.maxRuns 3, while two plain starts were 'waiting for a slot'. Check whether an interrupted run with a queued continuation counts against the cap itself (so it needs live <= cap-1 to start), and whether continuations take the next free slot ahead of new starts as #226/#243 intend. Unstuck by driving it with rhiza run resume --foreground.
