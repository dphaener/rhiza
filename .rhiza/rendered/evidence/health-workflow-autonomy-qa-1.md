# health-workflow-autonomy-qa-1

evidence · `01M6WRKKVMZ946RVX7EDM0C3BP` · actor [claude-code-e8898787](../actor/claude-code-e8898787.md)

## result

PASS

## statement

QA by hand with this branch's build (0c176623) in a scratch repo: with no repo workflow rhiza health's autonomy part lists only 'workflow plan-milestone (builtin:plan-milestone): plan owner, accept owner, ship owner, signoff owner'; with .rhiza/workflows/bugfix.toml declaring autonomous = true it adds 'workflow bugfix (.rhiza/workflows/bugfix.toml): plan auto, accept auto, ship auto, signoff auto', and a partial [autonomy] file shows only its declared stage; rhiza health --json and rhiza api call health printed identical bytes (cmp).

## Links

### Out

- verifies [health-workflow-autonomy](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-workflow-autonomy.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
