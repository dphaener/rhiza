# stage-autonomy-from-workflow-verified-7e0092cc

evidence · `01M6WNGY1GCJ5A20QQT50BPH2E` · actor [claude-code-785308d6](../actor/claude-code-785308d6.md)

## result

PASS

## statement

Verifier at 7e0092cc (commit adds only the builder's actor record; no code). Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(workflow-plan-autonomy|workflow-stage-autonomy|stage-autonomy|plan-workflow-autonomous|plan-workflow-autonomous-left|autonomy-accept-rule)$': all PASS, with fake rex and agents. Read the scripts against the acceptance: workflow-plan-autonomy covers plan.toml autonomous=true starting the request's plan run --autonomous under rhiza.autonomy.plan=owner, plan=owner starting it without under auto, and a ship-change plan=auto leaving plan runs unchanged; workflow-stage-autonomy covers ship=auto starting a ready change set under ship=owner, signoff=owner starting run start without --autonomous under signoff=auto, plan=auto not affecting ship, and accept=owner (both a custom workflow and a plan.toml override) leaving drafts proposed under accept=auto; stage-autonomy passes unchanged for the no-declaration fallback.

## Links

### Out

- verifies [stage-autonomy-from-workflow](../capability/rhiza/orchestration/accept-policy/stage-autonomy-from-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
