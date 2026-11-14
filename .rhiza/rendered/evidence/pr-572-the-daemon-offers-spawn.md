# pr-572-the-daemon-offers-spawn

evidence · `01M7BG874D9V53W31Z58CA3S4F` · actor [dphaener](../actor/dphaener.md)

## head

3190f8f818d503a6cd5ac7f0fdd81ebfffbde71b

## pr

https://github.com/dphaener/rhiza/pull/572

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #572 (https://github.com/dphaener/rhiza/pull/572) passed CI at 3190f8f

## status

verified

## Links

### Out

- verifies [run-spawn](../capability/rhiza/workflows/workflow-runs/run-spawn.md)
- verifies [spawn-visible](../capability/rhiza/workflows/workflow-runs/spawn-visible.md)
- verifies [gate-step](../capability/rhiza/workflows/workflow-definitions/gate-step.md)
- verifies [run-handback](../capability/rhiza/workflows/workflow-runs/run-handback.md)
- verifies [spawn-durable](../capability/rhiza/orchestration/daemon/spawn-durable.md)
- verifies [handback-route](../capability/rhiza/workflows/workflow-runs/handback-route.md)

## History

- `01M79KGRYE538GJC5NA2VNYG4K` The daemon offers spawn, gate and hand-back primitives: an agent spawns journaled, visible child sessions; a gate step runs commands the daemon owns; any red hands back to a named agent step, capped, then the inbox
