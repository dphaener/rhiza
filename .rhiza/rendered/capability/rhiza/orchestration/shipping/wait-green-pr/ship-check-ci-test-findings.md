# ship-check-ci-test-findings

capability · `01M6MC4GXTJ7SWF8NV4SYXD3F9` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A Go test in internal/api parses sample --log-failed output (a nested TestScripts/review-check-worktree failure, two failing tests in one package, and an infra-only log with no '--- FAIL:' line) into the expected test names and none for the infra log; a testscript with a fake gh whose check fails again after its rerun and whose 'run view 42 --log-failed' prints a '--- FAIL: TestScripts/x' line shows ship check --json listing that test as a failure and finding; run with go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship' only

## statement

When a PR head's checks still fail after their once-per-head rerun, rhiza ship check reads each failed Actions run's log (gh run view <run> --log-failed), reports the failing Go tests and testscripts it names ('--- FAIL: TestX[/sub]' lines) as test failures, and records each as a finding with no leaf for the next fix step; a failed check with no Actions run id, or whose log names neither a failing test nor a recognised non-test cause (ship-check-ci-cause-findings), such as only an infra/setup error, stays a plain failure

## Links

### Out

- decided-by [ship-change-ci-fix-cycle](../../../../../decision/ship-change-ci-fix-cycle.md)

### In

- depends-on [ship-change-ci-fix-step](ship-change-ci-fix-step.md)
- depends-on [ship-check-runner-test-failures](ship-check-runner-test-failures.md)
- verifies [ship-check-ci-test-findings-verified](../../../../../evidence/ship-check-ci-test-findings-verified.md)
- verifies [ship-check-ci-test-findings-qa-1](../../../../../evidence/ship-check-ci-test-findings-qa-1.md)
- verifies [ship-check-ci-test-findings-verified-2](../../../../../evidence/ship-check-ci-test-findings-verified-2.md)
- verifies [pr-280-ship-change-fixes-red-pr-ci](../../../../../evidence/pr-280-ship-change-fixes-red-pr-ci.md)
- verifies [pr-455-ship-change-treats-a](../../../../../evidence/pr-455-ship-change-treats-a.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
