# ci-nontest-failures-plan

source · `01M76BDAQD86NK9DS654CM8YCM`

## quotes

```json
[]
```

## raw_sha256

89180d2d075cb3d02085ce233472d16be3f20b8c281021faf37ed6ed4632da75

## session

b5d43cfb-5283-4ba3-a002-01a1ccb2baa0

## summary

Autonomous plan run (no owner present). The planner read the change set's intent (2026-11-11, run 01M75WRZ / PR #442: CI's lint job failed in self-check on stale docs/reference, ci-triage called it infra and the run sat in wait-ci-rerun 45+ minutes), the wait-green-pr leaves (ship-check-ci-test-findings, ship-check-runner-test-failures, ship-change-ci-fix-step, ci-triage-infra-waits, ship-change-infra-wait-reruns from PR #441) and internal/api/ship.go (parseFailedTests, shipReadFailures, runShipTestFailures, shipInfraWait), the built-in build/review/ship-change workflows and the Makefile. It raised three questions to the arbiter: how to classify a non-test red as fixable, how runs keep checked-in generated files current, and what part 3 (wait-ci-rerun waiting forever) still needs given #441. The arbiter answered with Decision arbiter-ci-nontest: a cause recogniser beside parseFailedTests (rhiza check 'fail ...' lines, gofmt/fmt-check listings, vet and compiler 'file.go:line:col:' diagnostics, and as a fallback 'make: *** [target] Error N' with the lines before it) records CI findings so ci-triage passes to ci-fix, also on wait-ci-rerun's attempts; ci-fix's prompt widens to CI findings; a language-neutral project setting rhiza.ship.generate run by the build and review render steps, set here to a make docs target; nothing new for the wait itself, met by the recogniser, with one testscript replaying the 2026-11-11 case. The change set serves m4-runs-hold-up, as the sibling CI-infra change set 01M75R7F does. No owner words were quoted: there are none in an autonomous run.

## title

Plan of change set 01M76B50: ship-change routes deterministic non-test CI failures to ci-fix

## Links

### In

- derived-from [arbiter-ci-nontest](../decision/arbiter-ci-nontest.md)
- derived-from [ship-check-ci-cause-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-ci-cause-findings.md)
- derived-from [ci-fix-fixes-ci-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-fix-fixes-ci-findings.md)
- derived-from [render-runs-generate](../capability/rhiza/orchestration/shipping/wait-green-pr/render-runs-generate.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
