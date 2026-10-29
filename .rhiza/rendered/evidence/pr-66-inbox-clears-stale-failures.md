# pr-66-inbox-clears-stale-failures

evidence · `01M63X2WGVAKZ08F7W46WF1T91` · actor [claude-code-orchestrator](../actor/claude-code-orchestrator.md)

## head

cb511160e5ac536a178c3d109e95e541eb357894

## pr

https://github.com/dphaener/rhiza/pull/66

## result

1 check passed: check SUCCESS

## statement

Pull request #66 (https://github.com/dphaener/rhiza/pull/66) passed CI at cb51116

## status

verified

## Links

### Out

- verifies [inbox-clears-stale](../capability/rhiza/orchestration/inbox/inbox-clears-stale.md)
- verifies [abandon-failed-run](../capability/rhiza/workflows/workflow-runs/abandon-failed-run.md)
- verifies [run-abandon](../capability/rhiza/workflows/workflow-runs/run-abandon.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
