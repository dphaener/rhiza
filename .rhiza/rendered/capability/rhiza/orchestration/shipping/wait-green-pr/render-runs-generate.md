# render-runs-generate

capability · `01M76BDZEWD3YNK3ZSSAZ35FRZ` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/workflow passes with the build and review render steps running the generate setting before rhiza render and committing its changes; a testscript render-generate shows rhiza daemon settings --json listing rhiza.ship.generate with an empty default, a build run's render step with it set to a command that rewrites a tracked file committing that file with the Change trailer, an empty setting skipping, and a failing command failing the step; make docs on this repo leaves docs/reference unchanged and make self-check passes; run with go test -trimpath ./internal/workflow/... ./internal/api/..., go test -trimpath ./cmd/rhiza -run 'TestScripts/(render-generate|describe-golden)' and make self-check only

## statement

A language-neutral project setting rhiza.ship.generate holds a shell command that regenerates the project's checked-in generated files with the branch's own code (empty by default, and empty skips); the render steps of the built-in build and review workflows run it in the run's worktree before rhiza render, stage and commit whatever it changed in the same commit with the same Change trailer, and fail the step when it fails; rhiza daemon settings lists it as a project setting, and this repo sets it to a new Makefile target (make docs) that builds the branch's binary in a temp folder and runs its rhiza docs, so stale docs/reference never reaches CI

## Links

### Out

- decided-by [arbiter-ci-nontest](../../../../../decision/arbiter-ci-nontest.md)
- derived-from [ci-nontest-failures-plan](../../../../../source/ci-nontest-failures-plan.md)

### In

- verifies [render-runs-generate-verified](../../../../../evidence/render-runs-generate-verified.md)
- verifies [render-runs-generate-qa-1](../../../../../evidence/render-runs-generate-qa-1.md)
- verifies [pr-455-ship-change-treats-a](../../../../../evidence/pr-455-ship-change-treats-a.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
