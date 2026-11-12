# ci-fix-fixes-ci-findings

capability · `01M76BDZC7N8TC2CC7WA2RARN9` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

go test -trimpath ./internal/workflow -run TestShipChangeCIFixPrompt asserts the ci-fix prompt speaks of CI findings, not only failing tests, and keeps the reproduce-locally rules; a testscript run-ship-ci-lint-fix replays the 2026-11-11 case with a fake gh whose lint job fails again after its rerun with 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and a make Error line, and a fake fixer: the run goes conflict-triage, ci-triage, ci-fix, back to green and signoff, with no wait-ci-rerun visit in the journal; run with go test -trimpath ./internal/workflow/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship' only

## statement

ship-change's ci-fix prompt tells the fixer to read the CI findings (failing tests and recognised non-test causes such as self-check, fmt-check, vet or make target failures) from rhiza ship check --json, reproduce each locally with the repo's scoped command for it (for Rhiza: make self-check, make fmt-check, or go vet on the package), fix it or regenerate the stale files, and commit once that command passes, so a deterministic lint or self-check red goes ci-triage -> ci-fix and never to wait-ci-rerun

## Links

### Out

- depends-on [ship-check-ci-cause-findings](ship-check-ci-cause-findings.md)
- decided-by [arbiter-ci-nontest](../../../../../decision/arbiter-ci-nontest.md)
- derived-from [ci-nontest-failures-plan](../../../../../source/ci-nontest-failures-plan.md)

### In

- verifies [ci-fix-fixes-ci-findings-verified](../../../../../evidence/ci-fix-fixes-ci-findings-verified.md)
- verifies [ci-fix-fixes-ci-findings-qa-1](../../../../../evidence/ci-fix-fixes-ci-findings-qa-1.md)
- verifies [pr-455-ship-change-treats-a](../../../../../evidence/pr-455-ship-change-treats-a.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
