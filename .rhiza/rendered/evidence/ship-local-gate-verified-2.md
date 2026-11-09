# ship-local-gate-verified-2

evidence · `01M6YVWVFCQN6QGXFZXCW8GEHQ` · actor [claude-code-40c7ef9f](../actor/claude-code-40c7ef9f.md)

## result

PASS

## statement

Verifier claude-code-40c7ef9f, attempt 2 at a1dc9868. go test -trimpath ./internal/api -run TestShipLocalGate -v: PASS (unstable runs the gate in the worktree, stable skips, no gate skips, failing gate fails with its output, plus OncePerHead). go test -trimpath ./internal/workflow: ok (local-gate, gate-rebase, gate-land, gate-rebase-merge before the push steps of ship-change.toml, on.fail = gate-fix). go test -trimpath ./cmd/rhiza -run 'TestScripts/(ship-local-gate|describe-golden)$': PASS; the script's own repo shows daemon settings --json with rhiza.ship.localGate project-scoped, empty default, source default, then daemon set to make check reading value make check from git config. The scripts had to run with TMPDIR set outside the repo and GOTMPDIR/RHIZA_ROOT unset: the run's GOTMPDIR sits under ~/code/rhiza/.git, so rhiza init in the script found the outer .rhiza (an environment issue, not the leaf). Diff read: removes the .rhiza/settings.json source, its test and golden wording, per Decision ship-local-gate-acceptance-from-commit; no leftover references; gofmt and vet clean on internal/api.

## Links

### Out

- verifies [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
