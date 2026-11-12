# ship-check-ci-cause-findings-verified

evidence · `01M76J5FCYDTY0QWKK9MZE162W` · actor [claude-code-ee6ff33b](../actor/claude-code-ee6ff33b.md)

## result

PASS

## statement

Verifier at 49149bb7a: read the diff (parseCICauses beside parseFailedTests in internal/api/ship.go, wired into shipReadFailures and shipInfraWait; --test-failures and findings descriptions name the causes; fmt header matches the Makefile's 'unformatted files:'). go test -trimpath ./internal/api/ -run 'TestParseCICauses|TestShipCheckCIFindingsNonTest|TestShipInfraWaitSettlesOnCause' -v -count=1: all three PASS. They cover gh job/step/timestamp-prefixed logs: a self-check log with two 'fail render docs/reference/...' lines plus 'make: *** [Makefile:186: self-check] Error 1' becomes one finding named 'lint / Run make lint: self-check' quoting the fail lines; a fmt-check listing; a go vet 'x.go:12:3:' diagnostic; a bare 'make: *** [build] Error 2' with the lines before it; apt 'Temporary failure resolving', billing and failing-test logs give no cause. Ship check with a fake gh records the finding in ci_findings; --test-failures exits 0 on it and 1 on the infra log; --settled --changed --rerun-infra settles on an attempt with a recognised cause. go test -trimpath ./cmd/rhiza -run TestScripts/describe-golden -count=1: ok (golden regenerated). go test -trimpath -timeout 45m ./internal/api/...: every test passed (no '--- FAIL'); the package exited 1 only on the rootguard check, because .rhiza-root/usage.live.json changed, which the live root daemon writes and rootguard does not list as owned; this predates the leaf, is unrelated to its diff, and is filed as draft change set 01M76J4VQ6PMFR6T23Y957E8HK (discovered-from). The first run without -timeout hit go test's 10m default under host load ~12.

## Links

### Out

- verifies [ship-check-ci-cause-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-ci-cause-findings.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
