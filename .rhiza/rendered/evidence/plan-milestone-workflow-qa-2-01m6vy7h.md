# plan-milestone-workflow-qa-2-01m6vy7h

evidence · `01M6X00J681GF5KARE3P5GH09G` · actor [claude-code-bd240cbe](../actor/claude-code-bd240cbe.md)

## result

PASS

## statement

QA cycle 2, built binary: workflow ls lists 'plan-milestone builtin builtin:plan-milestone'; a record change set's plan run journals workflow plan-milestone; rhiza daemon --once on a queued milestone request with rhiza.autonomy.plan auto detaches 'run start --foreground --workflow plan-milestone --pick default ... --milestone' with no --autonomous (RHIZA_TEST_DETACH_LOG); a repo plan-milestone.toml declaring autonomous = true makes the run autonomous.

## Links

### Out

- verifies [plan-milestone-workflow](../capability/rhiza/workflows/built-in-workflows/plan-milestone-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
