# note-01m6wcd1b7smjz

note · `01M6WCD1B7SMJZN34CTXPKVJST`

## statement

Round 2 attempt 3 (run attempt 6), 2026-11-08: still blocked. origin/main is b696b22c with no internal/dispatch; PR #367 (cs/workflow-dispatch, 3d0ab412) is open, DIRTY against main, with no checks. Per decision @bugfix-autonomous-waits-for-dispatch the leaf was not built, cs/workflow-dispatch was not merged here and the acceptance was not weakened. Further build attempts cannot pass: the run should pause this leaf until #367 is brought up to date and lands, or move the leaf into a follow-up change set that depends on dispatch so this change set can ship.

## Links

### Out

- discusses [bugfix-needs-workflow-dispatch](../question/bugfix-needs-workflow-dispatch.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
