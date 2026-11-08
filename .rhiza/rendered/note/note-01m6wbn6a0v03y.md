# note-01m6wbn6a0v03y

note · `01M6WBN6A0V03YFNW088GG64SN`

## statement

Still open at attempt 3 (2026-11-08): workflow dispatch (01M6VY7G6HPJC3V0Q18E8EF9ZA) is still on cs/workflow-dispatch (d1f24b16, in review) and not on origin/main (cf81def4) or this branch, so no internal/dispatch package maps kind bug to a repo bugfix.toml. Merging that branch here would ship another change set inside this one. Land dispatch, then rebuild this leaf on top of it.

## Links

### Out

- discusses [bugfix-needs-workflow-dispatch](../question/bugfix-needs-workflow-dispatch.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
