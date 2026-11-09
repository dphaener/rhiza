# ship-local-gate-post-land-fix-loop

decision · `01M6YYAR61GBW4PAYA09P25W1N`

## reasoning

I read the finding's full text in the run journal, internal/workflow/builtin/ship-change.toml and internal/api/ship_rebase.go and ship_openpr.go. The defect is real. gate-fix only returns through render-pr -> local-gate -> open-pr -> wait-pr -> pr-green -> signoff -> gate-rebase -> rebase -> ... -> gate-land -> land, and land has no max_visits (cap 1) and was already entered, so any gate failure after land fails the run at land's second visit. Raising land's cap would re-run rhiza ship land on a landed change set and a second owner sign-off for a fix to an already-signed-off change, so I rejected it. Reusing gate-rebase-merge -> rebase-merge as the return path would not push: runShipRebase returns before pushing when origin/main is already an ancestor of HEAD, which it is after shipKeepGateFailure's merge or once main has not moved. ship open-pr pushes any descendant of the pushed head and does not depend on model land state, so it is the right push step after land. A dedicated post-land loop keeps the leaf's statement (a failing gate goes to a fixer, like a red CI, before every push) true for the two post-land gates. It also keeps the change to the workflow file and its test, within this change set's intent.

## statement

Keep leaf @ship-local-gate as written and fix f1 in ship-change.toml with a separate fix loop for gate failures after land, which never re-enters the pre-land path or land. gate-rebase-merge's on.fail goes to a new agent step gate-fix-merge (same prompt as gate-fix, max_visits = 2). rebase-merge's on.fail goes to a new rebase-merge-gate-failed step (rhiza ship gate --failed, max_visits = 3, on.pass = gate-fix-merge) instead of the shared rebase-gate-failed. gate-fix-merge goes to render-merge (the render-pr command), then gate-push-merge (rhiza ship gate, stall 70m; on.fail = gate-fix-merge), then push-merge (rhiza ship open-pr, which pushes without forcing because the fix descends from the pushed head or the kept merge), then the existing wait-rebased-merge -> rebased-merge-green -> merge. The push goes through ship open-pr, not back through rebase-merge, because ship rebase pushes nothing when the branch already holds origin/main. Every new step sets max_visits so that a third failing gate fails the run cleanly, with the gate's output, rather than hitting land's cap of 1. No second sign-off after land, matching merge's existing on.fail to gate-rebase-merge. Tests: TestShipChangeLocalGate asserts that no step reachable from land routes back to land or to a pre-land step, and that both post-land gate failures reach gate-fix-merge.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [leaf-failed-01m6ygg6-1](../question/leaf-failed-01m6ygg6-1.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
