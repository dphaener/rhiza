# pr-467-runs-clean-up-rhiza-s-own

evidence · `01M777HPPDF2CW2MEP7443EPDQ` · actor [dphaener](../actor/dphaener.md)

## head

ff69710fc560dfcf14f23405d1e09979a481e711

## pr

https://github.com/dphaener/rhiza/pull/467

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #467 (https://github.com/dphaener/rhiza/pull/467) passed CI at ff69710

## status

verified

## Links

### Out

- verifies [run-end-removes-leftovers](../capability/rhiza/workflows/workflow-runs/run-end-removes-leftovers.md)
- verifies [daemon-prunes-orphan-leftovers](../capability/rhiza/orchestration/daemon/daemon-prunes-orphan-leftovers.md)

## History

- `01M763YK58RWNJ35P9S8BSFKVQ` Runs clean up Rhiza's own leftovers: when a run ends, Rhiza removes the scratch folders of the agent sessions it started and its worktree's index database; the daemon prunes orphans (index files whose worktree is gone, old session scratch folders) so disk use no longer grows with runs.
