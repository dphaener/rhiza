# stage-autonomy-from-workflow-verified

evidence · `01M6WEDQZKZG6XCTF84ZZE6WR0` · actor [claude-code-b52c8ced](../actor/claude-code-b52c8ced.md)

## result

PASS

## statement

Verified at deb4f262 (leaf commits 54220972, 7391c16a, deb4f262). Read the diff: daemonPoll1 and daemonDemand read ship and signoff through workflowStageMode(ship-change), which uses the workflow's declaration and falls back to rhiza.autonomy.<stage>; runRunGate and the autonomy accept rule read the accept stage from the run's recorded workflow autonomy (stateStageMode, acceptStageMode) and fall back to the setting; planRunAutonomous already reads the plan or plan-milestone workflow's plan stage. Ran go test -trimpath ./cmd/rhiza -run TestScripts/<name> with fake rex and RHIZA_TEST_DETACH_LOG, and all passed: workflow-plan-autonomy (plan.toml autonomous=true + setting owner starts the feature request's plan run with --autonomous; plan.toml plan=owner + setting auto starts it without; ship-change.toml plan=auto + setting owner starts it without), workflow-stage-autonomy (ship-change ship=auto starts a ready change set under setting owner; signoff=owner starts run start without --autonomous under setting auto; plan=auto alone leaves ship on the setting; accept=owner in both a custom workflow and a repo plan.toml leaves the gate's drafts proposed under accept auto), and unchanged stage-autonomy, plan-milestone-workflow, autonomy-accept-rule, autonomous-gate-rule, run-autonomy-source, run-autonomous, autonomous-review, review-signoff-autonomous, plan-workflow-autonomous and plan-workflow-autonomous-left. go test -trimpath ./internal/api/... printed PASS for every test, and then the package's guard reported a write to ~/code/.rhiza-root/usage.live.json. That is the known environment leak under the enclosing ~/code root and is not caused by this leaf. Both findings from attempt 1 are covered.

## Links

### Out

- verifies [stage-autonomy-from-workflow](../capability/rhiza/orchestration/accept-policy/stage-autonomy-from-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
