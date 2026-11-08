# workflow-autonomy-declaration-verified

evidence · `01M6W8RN0SZBRV40YS1P0NA6M6` · actor [claude-code-ca571187](../actor/claude-code-ca571187.md)

## result

PASS

## statement

Verifier at 6e3f6571: go test -trimpath ./internal/workflow/... passes; TestAutonomyDeclaration parses autonomous = true (all four stages auto), autonomous = false (all owner), a partial [autonomy] table (plan owner, ship auto, accept/signoff undeclared) and no declaration (every stage undeclared, so it follows rhiza.autonomy.<stage>) through Definition.Autonomy, and both Parse and Validate refuse both forms in one file ("autonomous and [autonomy] are both set"), an unknown stage ("autonomy.build: unknown stage"), a bad value ("autonomy.ship: maybe is not a mode") and a non-boolean autonomous, each naming the key. go test -trimpath ./cmd/rhiza -run TestScripts/workflow-check$ passes: rhiza workflow check accepts a file with autonomous = true and exits non-zero with the same messages for a bad [autonomy] value and for both forms. Diff read: parseAutonomy runs on every load (parse) and in Validate; [autonomy] keys are excluded from the unknown-key check and checked by parseAutonomy instead.

## Links

### Out

- verifies [workflow-autonomy-declaration](../capability/rhiza/workflows/workflow-definitions/workflow-autonomy-declaration.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
