# bugfix-needs-workflow-dispatch

question · `01M6WB3CJ5R7561BXEX3X6HP2N`

## statement

Workflow dispatch (change set 01M6VY7G6HPJC3V0Q18E8EF9ZA, branch cs/workflow-dispatch) is not on this branch or main, so nothing maps a request kind to a repo workflow: planRunArgs always picks plan or plan-milestone. The leaf's acceptance, a bugfix.toml mapped to kind bug, cannot pass until dispatch lands. Land dispatch, then rebuild this leaf on top of it.

## Links

### Out

- blocks [bugfix-workflow-autonomous](../capability/rhiza/orchestration/accept-policy/bugfix-workflow-autonomous.md)

### In

- answers [bugfix-builds-on-landed-dispatch](../decision/bugfix-builds-on-landed-dispatch.md)
- discusses [note-01m6wbn6a0v03y](../note/note-01m6wbn6a0v03y.md)
- discusses [bugfix-dispatch-still-unlanded-attempt5](../note/bugfix-dispatch-still-unlanded-attempt5.md)
- discusses [note-01m6wcd1b7smjz](../note/note-01m6wcd1b7smjz.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
