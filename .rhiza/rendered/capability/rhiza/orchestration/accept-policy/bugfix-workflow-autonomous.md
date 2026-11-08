# bugfix-workflow-autonomous

capability · `01M6VYPW4D8366CNFT3EH72465` · parent [accept-policy](../accept-policy.md)

## acceptance

A testscript with a fake rex, agents and gh: a repo .rhiza/workflows/bugfix.toml declaring autonomous = true and mapped to kind bug, rhiza.autonomy.plan owner and rhiza.autonomy.signoff owner; a queued bug request's plan run starts autonomous, its gate accepts the change set, the daemon starts the ship run autonomous, and the run lands with no inbox item for the owner at any point

## statement

A repo workflow declaring autonomous = true, picked for bug requests by workflow dispatch, takes a bug request from its plan run through its ship run with no owner step while rhiza.autonomy.plan stays owner

## Links

### Out

- decided-by [autonomy-per-workflow-scope](../../../../decision/autonomy-per-workflow-scope.md)
- derived-from [autonomy-per-workflow-interview](../../../../source/autonomy-per-workflow-interview.md)

### In

- blocks [bugfix-needs-workflow-dispatch](../../../../question/bugfix-needs-workflow-dispatch.md)
- blocks [leaf-failed-01m6w86n-1](../../../../question/leaf-failed-01m6w86n-1.md)
- blocks [leaf-failed-01m6w86n-2](../../../../question/leaf-failed-01m6w86n-2.md)
- verifies [bugfix-workflow-autonomous-verified](../../../../evidence/bugfix-workflow-autonomous-verified.md)
- verifies [bugfix-workflow-autonomous-verified-2](../../../../evidence/bugfix-workflow-autonomous-verified-2.md)
- verifies [pr-376-autonomy-per-workflow](../../../../evidence/pr-376-autonomy-per-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
