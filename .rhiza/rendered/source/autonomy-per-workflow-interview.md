# autonomy-per-workflow-interview

source · `01M6VYN0TMXZMDF553PH6QCDFH`

## quotes

```json
[]
```

## raw_sha256

ce6acdbb38fc12c9273540a1432c45cd91de5d12d8b8c318a18964418a884fc5

## session

e4df472e-8f7a-4104-8aa1-01a1b7e438c0

## summary

Autonomous plan run for change set 01M6VY7H (Autonomy per workflow). No owner took part. The planner read the autonomy code (autonomyStage in internal/api/autonomy.go, the daemon's plan, ship and signoff reads in request.go and daemon.go, the accept read at the plan gate, the milestone hard-codes in daemonStartPlan and run start's --autonomous refusal) and raised three questions: what the change set serves with no open milestone, which workflow's declaration governs each stage, and how milestone planning gets its owner-only default. The arbiter answered all three in Decision autonomy-per-workflow-scope: serve no milestone and wait on scheduling as policy (01M6VY7GRV) (revised after the interview: the orchestrator linked the change set to milestone-25, the backlog milestone, so it serves that and milestone-check passes); each stage takes its autonomy from the workflow of the run that performs it, falling back to rhiza.autonomy.<stage>; milestone planning gets a built-in plan-milestone workflow declaring autonomous = false and the hard-code goes. The planner then wrote six leaves: the declaration in the workflow file, stage resolution in the daemon and runs, the run record and run status showing each stage's source, health listing declaring workflows, the plan-milestone built-in replacing the milestone hard-codes, and the end-to-end bugfix check, which needs workflow dispatch (01M6VY7G6H).

## title

Autonomy per workflow: autonomous plan interview, answered by the arbiter

## Links

### In

- derived-from [autonomy-per-workflow-scope](../decision/autonomy-per-workflow-scope.md)
- derived-from [workflow-autonomy-declaration](../capability/rhiza/workflows/workflow-definitions/workflow-autonomy-declaration.md)
- derived-from [stage-autonomy-from-workflow](../capability/rhiza/orchestration/accept-policy/stage-autonomy-from-workflow.md)
- derived-from [run-autonomy-source](../capability/rhiza/workflows/workflow-runs/run-autonomy-source.md)
- derived-from [health-workflow-autonomy](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-workflow-autonomy.md)
- derived-from [plan-milestone-workflow](../capability/rhiza/workflows/built-in-workflows/plan-milestone-workflow.md)
- derived-from [bugfix-workflow-autonomous](../capability/rhiza/orchestration/accept-policy/bugfix-workflow-autonomous.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
