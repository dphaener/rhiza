# plan-milestone-workflow-verified

evidence · `01M6WE8KMW201KZTPXZV1BPTRQ` · actor [claude-code-96dbe264](../actor/claude-code-96dbe264.md)

## result

PASS

## statement

At 9062cd7f (attempt 2): the attempt-1 finding is fixed: runRunStart checks agent steps on the request's workflow (plan-milestone when --milestone) before resolving the change set, and again only when a plan run on a record change set swaps in plan-milestone; go test -trimpath ./internal/api -run 'TestRunStartAgentRefusalsBeforeWrite|Milestone|PlanMilestone|RunStart' passes. go test ./internal/workflow passes, incl. TestBuiltinPlanMilestone (steps equal to plan's, every stage owner). Testscripts plan-milestone-workflow (workflow ls lists plan-milestone builtin; dry-run --milestone with rhiza.autonomy.plan=auto shows workflow plan-milestone, not autonomous; daemon starts the milestone request on plan-milestone without --autonomous; a repo plan-milestone.toml with autonomous = true makes the daemon pass --autonomous), plan-milestone-backlog, plan-workflow, plan-workflow-open, plan-workflow-pr, plan-requests, daemon-milestone-request, run-start-new, describe-golden and workflow-ls pass. The four plan scripts now expect the planned change set left proposed, consistent with plan-milestone's owner accept stage. The full internal/api package at load ~12 only failed TestReadBudget (timing budget) and the root usage.live.json leak guard, both environmental and outside this diff.

## Links

### Out

- verifies [plan-milestone-workflow](../capability/rhiza/workflows/built-in-workflows/plan-milestone-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
