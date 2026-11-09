# ship-local-gate-verified

evidence · `01M6YQDCC4BM9HH988AKH67J4W` · actor [claude-code-5a166b01](../actor/claude-code-5a166b01.md)

## result

PASS

## statement

Verifier claude-code-5a166b01 at e6319c2b (attempt 2). Read the diff: ship-change.toml now has a rhiza ship gate step before every pushing step (local-gate->open-pr, gate-rebase->rebase, gate-land->land, gate-rebase-merge->rebase-merge, and merge on.fail now goes to gate-rebase-merge), each on.fail to gate-fix; ship rebase also gates the rebased head before its forced push and puts the branch back on failure; .rhiza/settings.json commits rhiza.ship.localGate = make check. Ran in a clean env (env -i; this session's RHIZA_* vars made the testscripts' rhiza init resolve the enclosing repo, at the parent commit too): go test -trimpath ./internal/api -run TestShipLocalGate PASS (unstable runs in the worktree, stable skips, no gate skips, failing gate fails with output, plus TestShipLocalGateOncePerHead, TestShipLocalGateSettingsFile); go test -trimpath ./internal/workflow PASS (TestShipChangeLocalGate checks every way into each push step goes through its gate); go test -trimpath ./cmd/rhiza -run 'TestScripts/(ship-local-gate|run-ship-change|run-ship-change-rebase|run-ship-ci-fix|describe-golden)$' ok. Full ./internal/api -v: 0 failing tests; the package's root-state guard flagged only ~/code/.rhiza-root/usage.live.json, which the live root daemon rewrites on its own (seen changing outside any test). A binary built from e6319c2b, run in this checkout, prints rhiza daemon settings --json {key: rhiza.ship.localGate, value: make check, source: .rhiza/settings.json, scope: project}; git config has no such key. Both attempt-1 findings are fixed.

## Links

### Out

- verifies [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
