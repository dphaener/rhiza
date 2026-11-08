# plan-milestone-workflow-verified-3

evidence · `01M6WWZAQN5WFXDDRYX7W99B3W` · actor [claude-code-c84d6f6c](../actor/claude-code-c84d6f6c.md)

## result

PASS

## statement

Verifier at b3d3cab3 (attempt 2, fixer commit adding planOwnerRefuses/planOwnerJournaled in internal/api/run.go). Read the diff: run start refuses --autonomous with a usage error naming workflow plan-milestone and where its declaration was read when the resolved plan workflow declares the plan stage owner (both the --new path and the record change set path), and run resume refuses a run journaled autonomous whose plan stage is owner. Ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-milestone-workflow|run-start-milestone-autonomous|workflow-ls|daemon-milestone-request|plan-workflow-autonomous|plan-workflow-autonomous-left|workflow-plan-autonomy|workflow-stage-autonomy|plan-run-milestone-env|plan-workflow|describe-golden)$': ok, with fake rex and no real agents. plan-milestone-workflow.txtar shows workflow ls listing plan-milestone as builtin:plan-milestone; under rhiza.autonomy.plan=auto, a milestone request's daemon start runs --workflow plan-milestone without --autonomous; and a repo plan-milestone.toml with autonomous = true makes the daemon pass --autonomous. run-start-milestone-autonomous.txtar shows the exit-2 refusals that name builtin:plan-milestone and write no draft or journal, the repo-file override, resume refused 'as run X journaled it' (the journal records workflow plan-milestone; run.go:661 journals def.Name/WorkflowPath after milestonePlanWorkflow swaps in builtin:plan-milestone), and autonomous = true lifting the refusal. go test -trimpath ./internal/workflow/... ok (TestBuiltinPlanMilestone: plan-milestone steps equal plan's, every stage owner). go test -trimpath ./internal/api -run 'Milestone|RunStart|RunResume|PlanRun|Autonom' ok. A full internal/api package run hit the 600s go test timeout at VM load ~30 and was not completed; CI runs it.

## Links

### Out

- verifies [plan-milestone-workflow](../capability/rhiza/workflows/built-in-workflows/plan-milestone-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
