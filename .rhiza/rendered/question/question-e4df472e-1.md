# question-e4df472e-1

question · `01M6VYENVE0ZD9ZP8JYHATZQQY`

## session

e4df472e-8f7a-4104-8aa1-01a1b7e438c0

## statement

No milestone is open (rhiza health: roadmap 0 milestones; milestone 29's change sets have all landed), so change set 01M6VY7H (Autonomy per workflow) has nothing to serve, and the plan workflow's milestone-check step fails a plan whose change set serves no milestone. Its three sibling requests from decision workflows-are-the-rules (workflow dispatch 01M6VY7G6H, scheduling as policy 01M6VY7GRV, request --ref 01M6VY7JMM) are in the same spot. What should it serve? / Which workflow's autonomy declaration governs which stage of a change set? Today one plan run (plan workflow) covers plan and accept, and a later ship-change run covers ship and signoff; workflow dispatch (01M6VY7G6H) will record on the change set the workflow its request chose, so the ship pick follows it. / How does milestone planning get its owner-only default? Today milestone requests run the built-in 'plan' workflow with --milestone, and daemonStartPlan (internal/api/request.go:366-371) hard-codes that such a run is never autonomous.

## Links

### In

- answers [autonomy-per-workflow-scope](../decision/autonomy-per-workflow-scope.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
