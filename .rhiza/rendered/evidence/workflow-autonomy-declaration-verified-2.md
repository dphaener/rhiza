# workflow-autonomy-declaration-verified-2

evidence · `01M6WN28ZE5MGBTYRY9G8VEH4X` · actor [claude-code-29bf50db](../actor/claude-code-29bf50db.md)

## result

PASS

## statement

Verifier claude-code-29bf50db at 13ac9726: go test -trimpath ./internal/workflow/... passes; TestAutonomyDeclaration parses autonomous = true (all four stages auto), autonomous = false (all owner), a partial [autonomy] table (plan owner, ship auto, accept/signoff undeclared) and no declaration (none), and Parse and Validate both refuse both forms ('autonomous and [autonomy] are both set'), an unknown stage ('autonomy.build: unknown stage') and a bad value ('autonomy.ship: maybe is not a mode') plus a non-boolean autonomous. go test ./cmd/rhiza -run TestScripts/(workflow-check|workflow-stage-autonomy|plan-milestone-workflow) passes; workflow-check.txtar shows rhiza workflow check exiting non-zero with the bad-value and both-forms messages. By hand with a scratch build, files in a scratch repo's .rhiza/workflows checked by name gave the unknown-stage, both-forms and bad-value refusals naming the key (exit 1) and a partial table valid (exit 0). Read parseAutonomy in internal/workflow/workflow.go and its call from Parse and check.go; the leaf commit itself adds only model records, the code having landed earlier in this change set.

## Links

### Out

- verifies [workflow-autonomy-declaration](../capability/rhiza/workflows/workflow-definitions/workflow-autonomy-declaration.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
