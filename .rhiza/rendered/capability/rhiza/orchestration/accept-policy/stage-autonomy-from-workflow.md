# stage-autonomy-from-workflow

capability · `01M6VYPVPYWCJVW98YSSNWHPDP` · parent [accept-policy](../accept-policy.md)

## acceptance

Testscripts with a fake rex and agents: a repo override of plan declaring autonomous = true starts a feature request's plan run with --autonomous while rhiza.autonomy.plan is owner, and one declaring plan = owner starts it without while the setting is auto; a plan workflow declaring accept = owner leaves its drafts proposed while rhiza.autonomy.accept is auto; a ship-change override declaring ship = auto starts a ready change set while rhiza.autonomy.ship is owner, and one declaring signoff = owner starts run start without --autonomous; a ship workflow's plan = auto changes nothing about plan runs; with no declaration each stage still follows its setting (stage-autonomy's scripts pass unchanged)

## statement

Each stage takes its autonomy from the workflow of the run that performs it, falling back to rhiza.autonomy.<stage> when that workflow declares nothing for it: the daemon starts a request's plan run autonomous by its plan workflow's plan stage, the plan gate accepts by that workflow's accept stage, the daemon starts or holds a ready change set's ship run by its ship workflow's ship stage and starts it autonomous by that workflow's signoff stage; a declaration for a stage the run does not perform is ignored for that run

## Links

### Out

- decided-by [autonomy-per-workflow-scope](../../../../decision/autonomy-per-workflow-scope.md)
- derived-from [autonomy-per-workflow-interview](../../../../source/autonomy-per-workflow-interview.md)

### In

- verifies [stage-autonomy-from-workflow-verified](../../../../evidence/stage-autonomy-from-workflow-verified.md)
- verifies [stage-autonomy-from-workflow-verified-7e0092cc](../../../../evidence/stage-autonomy-from-workflow-verified-7e0092cc.md)
- verifies [stage-autonomy-from-workflow-verified-59ffe16b](../../../../evidence/stage-autonomy-from-workflow-verified-59ffe16b.md)
- verifies [stage-autonomy-from-workflow-qa-2-01m6vy7h](../../../../evidence/stage-autonomy-from-workflow-qa-2-01m6vy7h.md)
- verifies [pr-376-autonomy-per-workflow](../../../../evidence/pr-376-autonomy-per-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
