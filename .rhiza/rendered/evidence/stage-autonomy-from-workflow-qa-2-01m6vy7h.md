# stage-autonomy-from-workflow-qa-2-01m6vy7h

evidence · `01M6X00KTJQEKPWDJSH5B42FPF` · actor [claude-code-bd240cbe](../actor/claude-code-bd240cbe.md)

## result

PASS

## statement

QA cycle 2, built binary, plan stage: with rhiza.autonomy.plan owner and a repo plan.toml declaring autonomous = true, rhiza daemon --once on a feature request detaches 'run start --foreground --workflow plan --pick default --new Add exports ... --autonomous'; a milestone request with the setting auto starts plan-milestone without --autonomous. The ship and signoff stages were checked only through their testscripts (workflow-stage-autonomy, daemon-root-ship-workflow), which pass.

## Links

### Out

- verifies [stage-autonomy-from-workflow](../capability/rhiza/orchestration/accept-policy/stage-autonomy-from-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
