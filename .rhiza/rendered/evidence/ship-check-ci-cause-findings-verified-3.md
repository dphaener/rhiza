# ship-check-ci-cause-findings-verified-3

evidence · `01M76NQ0NZ349JN4ZTXGJY7GXQ` · actor [claude-code-2df313ff](../actor/claude-code-2df313ff.md)

## result

PASS

## statement

Verifier, attempt 3, at db270e16a. Read the diff of 49149bb7a, abed08418 and db270e16a: parseCICauses and logBody in internal/api/ship.go, plus ship_test.go. Ran go test -trimpath -count=1 -v ./internal/api/ -run 'CICause|ShipInfraWait|Ship.*Cause|TestShipCheck': all pass, ok. That covers TestParseCICauses (gh-prefixed self-check log with fail render lines and 'make: *** [Makefile:186: self-check] Error 1' → a lint job/step finding quoting them; a fmt-check listing; a vet 'x.go:12:3:' diagnostic; a bare 'make: *** [build] Error 2' with the lines before it; apt 'Temporary failure resolving' and billing logs → none), TestShipCheckCIFindingsNonTest (fake gh: ci_findings recorded, --test-failures exit 0 on a cause and 1 on the infra log), TestShipInfraWaitSettlesOnCause (--settled --changed --rerun-infra settles on an attempt with a cause) and TestParseCICausesSkipsTestOutput (indented testscript output quoting 'fail render', a vet diagnostic or 'unformatted files:' yields no cause; the test is still found). This confirms cycle-2 quality f1 is fixed. go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden: ok. Every internal/api test passed (PASS) in the full go test -trimpath ./internal/api/..., but the package exited FAIL because rootguard reported '.rhiza-root/usage.live.json: file changed'. That file is rewritten continuously by live Claude Code status lines (internal/adapter/claudecode/statusline.go) and is not in rootguard's owned list. This leaf does not touch it, so it is a pre-existing guard gap unrelated to this leaf. Model f1 fixed: the ship-change-ci-fix-step, ship-check-ci-test-findings and ci-triage-infra-waits statements now say a recognised non-test cause passes ci-triage and goes to ci-fix.

## Links

### Out

- verifies [ship-check-ci-cause-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-ci-cause-findings.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
