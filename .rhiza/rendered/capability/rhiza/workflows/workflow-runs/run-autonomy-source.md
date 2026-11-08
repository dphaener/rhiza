# run-autonomy-source

capability · `01M6VYPVTC8M1XMBKTDTGBX14K` · parent [workflow-runs](../workflow-runs.md)

## acceptance

A testscript starts a plan run whose workflow declares plan = auto and leaves accept undeclared, then shows run status --json and run status naming plan auto from workflow <name> and accept <mode> from rhiza.autonomy.accept, and the run's started entry carrying both; describe-golden.txtar is regenerated

## statement

A run's journal records, for each stage the run performs, its autonomy mode and where it came from (the workflow's declaration, naming the workflow, or the rhiza.autonomy.<stage> setting), and rhiza run status shows that line for each stage

## Links

### Out

- decided-by [autonomy-per-workflow-scope](../../../../decision/autonomy-per-workflow-scope.md)
- derived-from [autonomy-per-workflow-interview](../../../../source/autonomy-per-workflow-interview.md)

### In

- verifies [run-autonomy-source-verified](../../../../evidence/run-autonomy-source-verified.md)
- verifies [run-autonomy-source-verified-5eb92da1](../../../../evidence/run-autonomy-source-verified-5eb92da1.md)
- verifies [run-autonomy-source-qa-2-01m6vy7h](../../../../evidence/run-autonomy-source-qa-2-01m6vy7h.md)
- verifies [pr-376-autonomy-per-workflow](../../../../evidence/pr-376-autonomy-per-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
