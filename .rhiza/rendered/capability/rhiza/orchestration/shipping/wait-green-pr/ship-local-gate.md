# ship-local-gate

capability · `01M6YCRKT22SQQH70HYX1GPFZW` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/api -run TestShipLocalGate shows the gate run in the run worktree when unstable, skipped when stable, skipped with no gate set, and a failing gate failing the step with its output; go test -trimpath ./internal/workflow passes with the gate step before each push step of ship-change.toml and its on.fail to a fixer; go test -trimpath ./cmd/rhiza -run TestScripts/ship-local-gate passes, showing in its own repo rhiza daemon settings --json listing rhiza.ship.localGate as a project setting with an empty default (source default), and rhiza daemon set rhiza.ship.localGate 'make check' then listing it with value make check from git config

## statement

While CI is unstable, the built-in ship-change workflow runs the project's local gate command (rhiza.ship.localGate; make gate for Rhiza itself) in the run's worktree before every push and pushes only when it exits 0, sending a failing gate to a fixer like a red CI; while CI is stable, or with no gate set, the step passes at once and the run keeps its scoped tests

## Links

### Out

- depends-on [ci-stability-state](ci-stability-state.md)
- decided-by [ci-unstable-gate-rules](../../../../../decision/ci-unstable-gate-rules.md)
- decided-by [no-ci-as-debugger](../../../../../decision/no-ci-as-debugger.md)
- decided-by [local-gate-trigger-identified-flakes](../../../../../decision/local-gate-trigger-identified-flakes.md)
- decided-by [local-gate-repeated-named-flakes](../../../../../decision/local-gate-repeated-named-flakes.md)

### In

- depends-on [health-ci-stability](../../daemon/daemon-status/health-command/health-ci-stability.md)
- depends-on ship-gate-full-log (retired)
- depends-on [ship-gate-run-log](ship-gate-run-log.md)
- depends-on [ship-gate-reuse-code-tree](ship-gate-reuse-code-tree.md)
- depends-on [ship-gate-base-env](ship-gate-base-env.md)
- blocks [review-01m6ycdrqqthsj4fv3sassw3v3-r3-f2](../../../../../question/review-01m6ycdrqqthsj4fv3sassw3v3-r3-f2.md)
- blocks [leaf-failed-01m6ygg6-1](../../../../../question/leaf-failed-01m6ygg6-1.md)
- verifies [ship-local-gate-verified](../../../../../evidence/ship-local-gate-verified.md)
- verifies [ship-local-gate-qa-1](../../../../../evidence/ship-local-gate-qa-1.md)
- verifies [ship-local-gate-verified-2](../../../../../evidence/ship-local-gate-verified-2.md)
- verifies [ship-local-gate-verified-5](../../../../../evidence/ship-local-gate-verified-5.md)
- verifies [ship-local-gate-verified-8](../../../../../evidence/ship-local-gate-verified-8.md)
- verifies [pr-403-stop-using-ci-as-a-debugger](../../../../../evidence/pr-403-stop-using-ci-as-a-debugger.md)
- verifies [pr-513-a-passing-local-gate-is-reused](../../../../../evidence/pr-513-a-passing-local-gate-is-reused.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
- `01M78T0G869AE4NQ9YGNB2DDGE` A passing local gate is reused while the code it checked hasn't changed, and the gate checks only what the change touches. Owner, 2026-11-12, agreed to file after the performance report. Today run 01M787SV ran the full make check 4 times (local-gate, gate-rebase-merge, after conflict-fix, after a rendered-only merge), ~20 min each serialized on the test lock; most reruns followed changes that touched no Go code (rendered model files, .rhiza/log entries). Wanted: (1) record the tree hash (or the hash of the non-.rhiza/ tree) a gate passed on and skip the gate when the next push has the same code tree; (2) the gate runs the change's scoped checks (go test for the packages the branch changes vs origin/main, their reverse dependencies, the affected testscripts, make self-check) instead of the whole suite, leaving the full suite to CI; (3) the project setting stays language-neutral (the gate command and its scoping come from project/workflow config, per the cleanup milestone). Estimated saving ~20 min per run on days the gate is on. Evidence: ~/orchestrator/perf-2026-11-12/report-prospective.md.
