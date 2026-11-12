# ship-check-ci-cause-findings

capability · `01M76BDZ8W61MZ0NAAEAQWF6NG` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

Go tests in internal/api parse sample --log-failed output carrying gh's job/step/timestamp prefix: a self-check log with two 'fail render docs/reference/...' lines and 'make: *** [Makefile:186: self-check] Error 1' into findings naming the lint job and step and quoting the fail lines, a fmt-check listing, a go vet 'x.go:12:3: ...' diagnostic, a bare 'make: *** [build] Error 2' with the lines before it, and an apt 'Temporary failure resolving' log and a billing log into none; ship check with a fake gh records those findings in ci_findings and --test-failures exits 0 on them and 1 on the infra log; --settled --changed --rerun-infra settles on an attempt whose log has a recognised cause; describe-golden regenerated; run with go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/describe-golden' only

## statement

When a PR head's failed jobs fail again after their one same-head rerun, rhiza ship check reads each failed Actions run's gh run view --log-failed output with a cause recogniser beside parseFailedTests and records each recognised cause as a CI finding in ship check's findings and the state file's ci_findings, exactly as a failing test is, so rhiza ship check --test-failures (ci-triage) passes on it: rhiza check or self-check 'fail ...' lines (such as 'fail render docs/reference/cli.md: differs from what rhiza docs writes'), gofmt/fmt-check file listings, go vet and compiler diagnostics ('<file>.go:<line>:<col>: ...'), and, as the fallback, a 'make: *** [<target>] Error N' line naming the target; each finding names what failed (the job and step, or the command) and its reason quotes the failing lines (for the make fallback, the lines just before it); a red with none of these lines or a failing test, or whose log gh cannot give, stays an infra failure as today; wait-ci-rerun's attempt check (shipInfraWait, --settled --changed --rerun-infra) applies the same recogniser, so an attempt red with a recognised cause settles and goes on like one red with a failing test; the --test-failures and findings descriptions name these causes

## Links

### Out

- decided-by [arbiter-ci-nontest](../../../../../decision/arbiter-ci-nontest.md)
- derived-from [ci-nontest-failures-plan](../../../../../source/ci-nontest-failures-plan.md)

### In

- depends-on [ci-fix-fixes-ci-findings](ci-fix-fixes-ci-findings.md)
- verifies [ship-check-ci-cause-findings-verified](../../../../../evidence/ship-check-ci-cause-findings-verified.md)
- verifies [ship-check-ci-cause-findings-qa-1](../../../../../evidence/ship-check-ci-cause-findings-qa-1.md)
- verifies [ship-check-ci-cause-findings-verified-3](../../../../../evidence/ship-check-ci-cause-findings-verified-3.md)
- verifies [pr-455-ship-change-treats-a](../../../../../evidence/pr-455-ship-change-treats-a.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
