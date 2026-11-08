# run-autonomy-source-qa-2-01m6vy7h

evidence · `01M6X00K5TD50DRYAEVF1T73K8` · actor [claude-code-bd240cbe](../actor/claude-code-bd240cbe.md)

## result

PASS

## statement

QA cycle 2, built binary: a plan run whose repo plan.toml declares only plan = auto, with rhiza.autonomy.accept owner, shows in run status 'Autonomy plan: auto from workflow plan', 'Autonomy accept: owner from rhiza.autonomy.accept', and ship/signoff from their settings; the journal's started entry carries an autonomy array of stage, mode, source and workflow; a plan-milestone run shows all four owner from workflow plan-milestone.

## Links

### Out

- verifies [run-autonomy-source](../capability/rhiza/workflows/workflow-runs/run-autonomy-source.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
