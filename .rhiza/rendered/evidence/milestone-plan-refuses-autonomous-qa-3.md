# milestone-plan-refuses-autonomous-qa-3

evidence · `01M6WRKM4KTMQ6YS7M7MHFJBKA` · actor [claude-code-e8898787](../actor/claude-code-e8898787.md)

## result

FAIL

## statement

QA with this branch's build (0c176623): run start is refused inside an agent session, so origin/main's run-start-milestone-autonomous.txtar (deleted on this branch) was run against the branch: rhiza run start --foreground --workflow plan --new 'Plan the milestone' --milestone --permission-mode manual --autonomous is not refused; it runs workflow plan-milestone with 'Autonomous: yes' while run status shows every stage 'owner from workflow plan-milestone'. The amended acceptance (usage error naming plan-milestone's declaration, no journal, no draft; run resume refusal) is not met, and no script on the branch checks it.

## Links

### Out

- verifies [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
