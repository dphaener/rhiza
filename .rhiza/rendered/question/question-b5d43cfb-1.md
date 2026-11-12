# question-b5d43cfb-1

question · `01M76B9HFC55PASR4NZW499ATF`

## session

b5d43cfb-5283-4ba3-a002-01a1ccb2baa0

## statement

Change set 01M76B50 (ship-change sends non-test CI failures to wait-ci-rerun): when the head's failed jobs have been rerun once and fail again, how should rhiza ship check decide that a failure with no failing test is fixable (route it to ci-fix) and not infra (wait-ci-rerun)? / Change set 01M76B50: how should build and review runs keep checked-in generated files (Rhiza's docs/reference from rhiza docs) current, so stale generated reference never reaches CI? / Change set 01M76B50, part 3 (wait-ci-rerun waiting forever): PR #441 (ship-change-infra-wait-reruns, landed) already makes wait-ci-rerun rerun the failed jobs itself with backoff over its 6h window and then end the run naming the failure. What more should this change set do for part 3?

## Links

### In

- answers [arbiter-ci-nontest](../decision/arbiter-ci-nontest.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
