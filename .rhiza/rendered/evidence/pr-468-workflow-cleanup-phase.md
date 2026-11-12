# pr-468-workflow-cleanup-phase

evidence · `01M776PGSNQCT7C6VN7WDF84C4` · actor [dphaener](../actor/dphaener.md)

## head

f83142ffd34bc42cecd9d0bd0e48ea63c30670dc

## pr

https://github.com/dphaener/rhiza/pull/468

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #468 (https://github.com/dphaener/rhiza/pull/468) passed CI at f83142f

## status

verified

## Links

### Out

- verifies [workflow-cleanup-steps](../capability/rhiza/workflows/workflow-definitions/workflow-file/workflow-cleanup-steps.md)
- verifies [workflow-check-cleanup](../capability/rhiza/workflows/workflow-authoring/workflow-check-cleanup.md)

## History

- `01M763YS08ZM13JZ065Q95YA2E` Workflow cleanup phase: a workflow file declares command steps that run once when a run of it ends, done, failed or abandoned; a failing cleanup step is logged and never changes the run's result; workflow check and describe cover it.
