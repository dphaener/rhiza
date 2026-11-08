# workflow-autonomy-declaration-qa-1

evidence · `01M6WRKKP5S6XE2D203GDGM78B` · actor [claude-code-e8898787](../actor/claude-code-e8898787.md)

## result

PASS

## statement

QA by hand with this branch's build (0c176623) in a scratch repo: rhiza workflow check passes .rhiza/workflows/bugfix.toml with autonomous = true and partial.toml with [autonomy] plan = "auto"; it exits 1 naming the key for autonomous plus [autonomy] together ('autonomous and [autonomy] are both set'), [autonomy] deploy ('autonomy.deploy: unknown stage; want plan, accept, ship, signoff'), plan = "yes" ('autonomy.plan: yes is not a mode') and autonomous = "yes" ('autonomous: yes is not a boolean'). go test ./internal/workflow passes.

## Links

### Out

- verifies [workflow-autonomy-declaration](../capability/rhiza/workflows/workflow-definitions/workflow-autonomy-declaration.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
