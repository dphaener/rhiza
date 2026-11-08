# bugfix-dispatch-still-unlanded-attempt5

note · `01M6WC1G8QEHHBJVJH4834Q657`

## statement

Round 2 attempt 2 (run attempt 5) of @bugfix-workflow-autonomous found workflow dispatch still not on main on 2026-11-08: origin/main is b696b22c, it has no internal/dispatch, and PR #367 (cs/workflow-dispatch) is open and CONFLICTING with main, so it has no CI and cannot merge as it stands. Per decision @bugfix-autonomous-waits-for-dispatch the leaf was not built and cs/workflow-dispatch was not merged here. Retrying this leaf cannot pass until #367 is brought up to date with main and lands; until then either pause this leaf or move it into a follow-up change set that depends on dispatch.

## Links

### Out

- discusses [bugfix-needs-workflow-dispatch](../question/bugfix-needs-workflow-dispatch.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
