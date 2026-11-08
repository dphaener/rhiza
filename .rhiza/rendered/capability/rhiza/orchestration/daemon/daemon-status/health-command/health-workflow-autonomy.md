# health-workflow-autonomy

capability · `01M6VYPVWZYB3NKSKBKACYXEXZ` · parent [health-command](../health-command.md)

## acceptance

health.txtar: with a repo workflow declaring autonomous = true the autonomy part's detail names that workflow, its file and plan, accept, ship and signoff auto, and the built-in plan-milestone with plan, accept, ship and signoff owner; with no repo declaration only plan-milestone is listed; health --json and api call health print the same bytes

## statement

rhiza health's autonomy part lists, after the rhiza.autonomy.* settings, every workflow the repo resolves (personal, repo or built-in) that declares autonomy, with each stage it declares and the file it comes from, so the owner sees where a run's autonomy will come from

## Links

### Out

- decided-by [autonomy-per-workflow-scope](../../../../../../decision/autonomy-per-workflow-scope.md)
- derived-from [autonomy-per-workflow-interview](../../../../../../source/autonomy-per-workflow-interview.md)

### In

- verifies [health-workflow-autonomy-verified](../../../../../../evidence/health-workflow-autonomy-verified.md)
- verifies [health-workflow-autonomy-verified-2](../../../../../../evidence/health-workflow-autonomy-verified-2.md)
- verifies [health-workflow-autonomy-qa-1](../../../../../../evidence/health-workflow-autonomy-qa-1.md)
- verifies [health-workflow-autonomy-qa-2-01m6vy7h](../../../../../../evidence/health-workflow-autonomy-qa-2-01m6vy7h.md)
- verifies [pr-376-autonomy-per-workflow](../../../../../../evidence/pr-376-autonomy-per-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
