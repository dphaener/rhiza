# plan-milestone-workflow-verified-2

evidence · `01M6WNG2WXVZZQ40PB57VQQE0B` · actor [claude-code-94ca1840](../actor/claude-code-94ca1840.md)

## result

PASS

## statement

At 20626479 (fakes only): go test -trimpath ./internal/workflow -run TestBuiltin passes, incl. TestBuiltinPlanMilestone (plan-milestone loads, Start and Steps DeepEqual plan's, every stage autonomy owner). go test -trimpath ./cmd/rhiza -run 'TestScripts/(plan-milestone-workflow|workflow-ls|daemon-milestone-request|plan-workflow-autonomous|run-autonomy-source)$' passes: with fake rex and rhiza.autonomy.plan auto, workflow ls lists plan-milestone as builtin:plan-milestone, the milestone dry-run reports workflow plan-milestone sourced from the workflow and not autonomous, the daemon detaches run start --workflow plan-milestone ... --milestone with no --autonomous, and a repo .rhiza/workflows/plan-milestone.toml declaring autonomous = true makes the daemon pass --autonomous even with plan owner. Read internal/api: daemonStartPlan/planRunArgs carry no hard-coded milestone rule (requestPick swaps plan for plan-milestone; milestonePlanWorkflow swaps it for run start --milestone and record change sets via planIsMilestone, journaling resolveWorkflow's builtin:plan-milestone path).

## Links

### Out

- verifies [plan-milestone-workflow](../capability/rhiza/workflows/built-in-workflows/plan-milestone-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
