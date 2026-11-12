# render-runs-generate-verified

evidence · `01M76KB7GPF1A21S12B5TJBF2J` · actor [claude-code-161be154](../actor/claude-code-161be154.md)

## result

PASS

## statement

Verified as claude-code-161be154 at 7e7d2ba60. go test -trimpath ./internal/workflow/... ok. go test -trimpath ./internal/api/... ok (379s) from a git-archive copy outside ~/code; run inside ~/code the tests all passed, but the root-state guard flagged .rhiza-root/usage.live.json, which the live root daemon writes, not this change. go test -trimpath ./cmd/rhiza -run 'TestScripts/(render-generate|describe-golden)' PASS: render-generate shows daemon settings listing rhiza.ship.generate with an empty default, empty skipping with nothing staged, a build run's render step committing the rewritten tracked file with the Change trailer and a clean tree, and a failing command failing (exit 1, output shown); the render steps chain it with && before rhiza render, so a failure fails the step. This repo's git config sets rhiza.ship.generate=make docs; make docs leaves docs/reference unchanged. make self-check at the leaf commit fails only on stale .rhiza/rendered (the run's render step refreshes it); in a scratch clone after rhiza render, make self-check is ok. Diff read: build.toml/review.toml run ship generate first and commit all staged changes with the Change trailer; ship.generate stages only paths whose status or content changed since before the command.

## Links

### Out

- verifies [render-runs-generate](../capability/rhiza/orchestration/shipping/wait-green-pr/render-runs-generate.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
