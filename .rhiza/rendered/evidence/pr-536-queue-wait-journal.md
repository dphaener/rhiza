# pr-536-queue-wait-journal

evidence · `01M79YFEMMS6K1RWT6262S0N97` · actor [dphaener](../actor/dphaener.md)

## head

77d460788cbcb9f7fede6dd1954c5439f3ab36d9

## pr

https://github.com/dphaener/rhiza/pull/536

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #536 (https://github.com/dphaener/rhiza/pull/536) passed CI at 77d4607

## status

verified

## Links

### Out

- verifies [queue-wait-journal-reader](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-reader.md)
- verifies [queue-wait-journal-lines](../capability/rhiza/orchestration/daemon/daemon-picks/queue-wait-journal-lines.md)

## History

- `01M76XSY0YM17EQ9K24CC4JRAK` The daemon journals why an accepted change set waits before its run starts (no free slot, load hold, usage hold, dependency, policy), and the perf analysis can read it. Evidence: accepted to ship-change start is 27m RECENT and 60m ALL median, cause not journaled (perf report rec 10).
