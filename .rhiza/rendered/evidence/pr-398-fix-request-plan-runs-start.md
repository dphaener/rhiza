# pr-398-fix-request-plan-runs-start

evidence · `01M6YJPD7XSH6MEEFQS7HQ0K47` · actor [dphaener](../actor/dphaener.md)

## head

0a41567189787da2068657bd002671ee7928046a

## pr

https://github.com/dphaener/rhiza/pull/398

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #398 (https://github.com/dphaener/rhiza/pull/398) passed CI at 0a41567

## status

verified

## Links

### Out

- verifies [daemon-starts-within-max-runs](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-starts-within-max-runs.md)
- verifies [daemon-resume-takes-slot](../capability/rhiza/orchestration/daemon/daemon-picks/daemon-resume-takes-slot.md)
- verifies [milestone-request-starts-at-once](../capability/rhiza/orchestration/daemon/plan-requests/milestone-request-starts-at-once.md)
- verifies [request-slot-kept](../capability/rhiza/orchestration/daemon/plan-requests/request-slot-kept.md)
- verifies [health-runs-against-cap](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-runs-against-cap.md)
- verifies [daemon-status-runs-cap](../capability/rhiza/orchestration/daemon/daemon-status/daemon-status-runs-cap.md)
- verifies [request-queue](../capability/rhiza/orchestration/daemon/plan-requests/request-queue.md)

## History

- `01M6XV4WGSPC346CQEV9DGYWR3` Fix Request plan runs start outside rhiza.daemon.maxRuns, so the machine runs far more agents than its cap: on 2026-11-08 ~13:05Z, with maxRuns 4 (then 3), the root daemon started the plan runs of requests 01M6XHDV2P65 and 01M6XKT75PHW in the same burst as ship runs of 01M6XHCX4FRE and 01M6WGEPSYWV, leaving 5 running and 1 interrupted (daemon status) and load at 23 on the 4-core VM. Earlier the same day 5 ran at cap 3. Plan runs should count against the machine-wide cap (or a small separate plan cap in rhiza.toml, as policy): a waiting request may reserve the next free slot (#226) but must not exceed the cap. rhiza health should show live runs against the cap.
