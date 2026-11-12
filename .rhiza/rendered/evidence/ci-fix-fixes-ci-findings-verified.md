# ci-fix-fixes-ci-findings-verified

evidence · `01M76K0S4NQFJN57TKPM57XC4Y` · actor [claude-code-e8f760a9](../actor/claude-code-e8f760a9.md)

## result

PASS

## statement

At ce2b6bdc2: go test -trimpath ./internal/workflow -run TestShipChangeCIFixPrompt PASS (asserts 'Fix the CI findings', 'failing tests, and recognised non-test causes', 'self-check, fmt-check, vet or build-target failure', 'Reproduce each finding locally', 'regenerate them', plus the kept quarantine/bug-filed/commit-only-once-passing/never-push-to-see rules); go test -trimpath ./internal/workflow/... ok; go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship' ok (255s), run-ship-ci-lint-fix PASS: fake gh lint check stays FAILURE after its one rerun with 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and 'make: *** [Makefile:186: self-check] Error 1' in --log-failed, fake fixer greps those lines from rhiza ship check --json; journal goes pr-green, conflict-triage, gate-triage, ci-triage pass -> ci-fix pass -> render-pr -> pr-green pass -> signoff, with no wait-ci-rerun. Diff read: the ship-change.toml ci-fix prompt and its .rhiza/workflows/docs.toml copy are identical, and run-ship-ci-fix.txtar's match string was updated. make self-check fails only on stale .rhiza/rendered actor/index files from run actor records, which render-pr regenerates, not on this leaf.

## Links

### Out

- verifies [ci-fix-fixes-ci-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-fix-fixes-ci-findings.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
