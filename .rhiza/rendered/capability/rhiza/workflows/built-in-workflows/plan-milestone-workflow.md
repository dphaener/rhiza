# plan-milestone-workflow

capability · `01M6VYPVZJ3PB7QGSN9XZVP9GN` · parent [built-in-workflows](../built-in-workflows.md)

## acceptance

Testscripts with a fake rex and agents: a milestone request's plan run journals workflow builtin:plan-milestone and starts without --autonomous while rhiza.autonomy.plan is auto; a repo .rhiza/workflows/plan-milestone.toml declaring autonomous = true makes the daemon start it with --autonomous; workflow ls lists plan-milestone as built-in; internal/workflow tests load the built-in and find its steps equal to plan's

## statement

The binary ships a plan-milestone workflow with the plan workflow's steps declaring autonomous = false, and a milestone request's plan run (and run start --milestone or on a record change set) runs it; daemonStartPlan's hard-coded non-autonomous milestone rule goes, so a repo changes milestone planning's autonomy by overriding plan-milestone like any other workflow

## Links

### Out

- decided-by [autonomy-per-workflow-scope](../../../../decision/autonomy-per-workflow-scope.md)
- derived-from [autonomy-per-workflow-interview](../../../../source/autonomy-per-workflow-interview.md)

### In

- verifies [plan-milestone-workflow-verified](../../../../evidence/plan-milestone-workflow-verified.md)
- verifies [plan-milestone-workflow-verified-2](../../../../evidence/plan-milestone-workflow-verified-2.md)
- verifies [plan-milestone-workflow-verified-3](../../../../evidence/plan-milestone-workflow-verified-3.md)
- verifies [plan-milestone-workflow-qa-2-01m6vy7h](../../../../evidence/plan-milestone-workflow-qa-2-01m6vy7h.md)
- verifies [pr-376-autonomy-per-workflow](../../../../evidence/pr-376-autonomy-per-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
