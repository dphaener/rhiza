# rebase-conflict-merged-main

decision · `01M6X37K4EAHNWGCWM74C9YF38`

## reasoning

The conflicts were additive: this branch's autonomy tests and fields and main's scheduled-trigger, watchdog and stall additions touched neighbouring lines, with no clash in behaviour. A merge (commit 246e8a33) keeps both without rewriting pushed history, the rendered index is regenerated from the merged log, and the touched packages' tests and affected scripts pass.

## statement

Resolve the conflicts by merging origin/main into cs/autonomy-per-workflow, keeping both sides: workflow_test.go keeps the autonomy declaration and plan-milestone tests beside main's step stall test, run info and the journal keep stage autonomy beside scheduled triggers and warnings, and .rhiza/rendered/index.md is re-rendered.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [rebase-conflict-autonomy-per-workflow](../question/rebase-conflict-autonomy-per-workflow.md)

## History

- `01M6VY7H9PYX3N129TV6Z8G2Z4` Autonomy per workflow: a workflow file declares its autonomy per stage (plan, accept, ship, signoff as auto or owner), with autonomous = true|false as shorthand, or omits it to follow rhiza.autonomy.*; the declaration overrides rhiza.autonomy.* for that workflow's runs (today autonomyStage in internal/api/request.go:294 and daemon.go:834-843 read only the settings), and rhiza health and the run record show where a run's autonomy came from. The built-in milestone planning workflow declares owner, as a default a repo can change (supersedes failed request 01M6MY89 that hard-codes milestone planning needing a human). Acceptance: a bugfix workflow declaring autonomous = true takes a bug request from plan through ship with no owner step while rhiza.autonomy.plan stays owner. Follows decision workflows-are-the-rules (PR #353): every rule is workflow policy the owner defines, Rhiza ships defaults and enforces only references between nodes.
