# health-ci-stability-verified

evidence · `01M6YQVA2DWJ671YE80V9VCJ19` · actor [claude-code-06740f04](../actor/claude-code-06740f04.md)

## result

PASS

## statement

Verifier claude-code-06740f04 at 6e998eb0: go test -trimpath ./internal/api -run TestHealthCI passes (clean record ok/stable; flakes warn with each test and count, TestGarden 1 flake, TestPond 2 flakes; main red warns naming main's latest CI result red; no-local-gate warning plus next pointer to rhiza.ship.localGate only when unstable and the gate is unset, absent when stable or gated). go test -trimpath ./cmd/rhiza -run 'TestScripts/(health-ci|health-projects|health|init|describe-golden)$' passes with a clean env (env -i): health-ci shows 'ok ci stable' and the warn line, flake detail and next in both rhiza health and health --json. Read the diff (health_ci.go, health.go ci part after main, golden and health-projects counts). The scripts fail under the run's inherited session env (rhiza init sees /home/rhiza/code/rhiza) for every init-based script on this machine, unrelated to the change.

## Links

### Out

- verifies [health-ci-stability](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-ci-stability.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
