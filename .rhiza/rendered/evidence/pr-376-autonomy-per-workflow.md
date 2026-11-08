# pr-376-autonomy-per-workflow

evidence · `01M6X7KAM5AKJV2YXQQYRRZT8D` · actor [dphaener](../actor/dphaener.md)

## head

f3a3250afe8ceccb6df096afea8d85aa032417a6

## pr

https://github.com/dphaener/rhiza/pull/376

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #376 (https://github.com/dphaener/rhiza/pull/376) passed CI at f3a3250

## status

verified

## Links

### Out

- verifies [workflow-autonomy-declaration](../capability/rhiza/workflows/workflow-definitions/workflow-autonomy-declaration.md)
- verifies [stage-autonomy-from-workflow](../capability/rhiza/orchestration/accept-policy/stage-autonomy-from-workflow.md)
- verifies [run-autonomy-source](../capability/rhiza/workflows/workflow-runs/run-autonomy-source.md)
- verifies [health-workflow-autonomy](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-workflow-autonomy.md)
- verifies [plan-milestone-workflow](../capability/rhiza/workflows/built-in-workflows/plan-milestone-workflow.md)
- verifies [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)
- verifies [bugfix-workflow-autonomous](../capability/rhiza/orchestration/accept-policy/bugfix-workflow-autonomous.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
