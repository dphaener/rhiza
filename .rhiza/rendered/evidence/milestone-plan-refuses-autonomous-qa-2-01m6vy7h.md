# milestone-plan-refuses-autonomous-qa-2-01m6vy7h

evidence · `01M6X00HYQN255FA1B57EHVCD3` · actor [claude-code-bd240cbe](../actor/claude-code-bd240cbe.md)

## result

PASS

## statement

QA cycle 2, built binary in a scratch root with a fake rex that runs commands: run start --workflow plan --new 'Ship v3' --milestone --autonomous (with and without --dry-run) and run start --autonomous on a record change set each exit 2 with 'workflow plan-milestone (builtin:plan-milestone) declares the plan stage owner', no draft and no .git/rhiza/runs; the record start without --autonomous runs plan-milestone; an intention plan with --autonomous succeeds with autonomous true; a repo plan-milestone.toml with autonomous = false refuses naming the repo file; run resume of that run, edited to be journaled autonomous, exits 2 naming 'as run ... journaled it'; a repo plan-milestone.toml with autonomous = true lets the milestone --autonomous start succeed.

## Links

### Out

- verifies [milestone-plan-refuses-autonomous](../capability/rhiza/workflows/autonomous-runs/milestone-plan-refuses-autonomous.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
