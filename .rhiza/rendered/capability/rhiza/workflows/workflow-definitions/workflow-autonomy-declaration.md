# workflow-autonomy-declaration

capability · `01M6VYPVKA8H196XPZQX6BHQZK` · parent [workflow-definitions](../workflow-definitions.md)

## acceptance

internal/workflow tests load autonomous = true, autonomous = false, a partial [autonomy] table and no declaration, reading each stage's declared mode or none, and refuse a file with both forms, an unknown stage and a bad value with a message naming the key; a testscript shows workflow check reporting the same refusal for a repo workflow file

## statement

A workflow file may declare its autonomy: a top-level autonomous = true|false sets all four stages (plan, accept, ship, signoff) to auto or owner, or an [autonomy] table sets any of them by name to auto or owner; a stage it leaves out follows rhiza.autonomy.<stage>, and workflow check (and every load) refuses both forms in one file, an unknown stage or a value other than auto or owner

## Links

### Out

- decided-by [autonomy-per-workflow-scope](../../../../decision/autonomy-per-workflow-scope.md)
- derived-from [autonomy-per-workflow-interview](../../../../source/autonomy-per-workflow-interview.md)

### In

- verifies [workflow-autonomy-declaration-verified](../../../../evidence/workflow-autonomy-declaration-verified.md)
- verifies [workflow-autonomy-declaration-verified-2](../../../../evidence/workflow-autonomy-declaration-verified-2.md)
- verifies [workflow-autonomy-declaration-qa-1](../../../../evidence/workflow-autonomy-declaration-qa-1.md)
- verifies [workflow-autonomy-declaration-qa-2-01m6vy7h](../../../../evidence/workflow-autonomy-declaration-qa-2-01m6vy7h.md)
- verifies [pr-376-autonomy-per-workflow](../../../../evidence/pr-376-autonomy-per-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
