# arbiter-ci-nontest

decision · `01M76BC9QXWXSR3K5B2XDDG2YE`

## reasoning

Read internal/api/ship.go and the built-in ship-change, build and review workflows. Ship check already reruns a red head's failed jobs once (shipRerunFailed). After that, shipReadFailures records findings only from parseFailedTests, so a lint or self-check red gets no finding: ci-triage fails and the run waits in wait-ci-rerun. That is the 2026-11-11 run. A failure that repeats after a rerun and prints a recognisable deterministic cause is fixable. Reusing the existing ci_findings and ci-triage route keeps the workflow unchanged. The change set's intent lists exactly these patterns, and says that only failures with no recognisable cause should count as infra. The make Error fallback catches project-specific targets without per-language parsing. For (2), build and review already have a deterministic render step that commits regenerated .rhiza/rendered. Running a project-configured generate command there makes generated files as dependable as rendered output, with no agent judgement involved. Keeping it a setting follows milestone 32's language-neutral settings: Rhiza's own rhiza docs stays project config, not built-in behaviour. For (3), #441 already removed the indefinite wait: wait-ci-rerun reruns the jobs itself and ends at 6h, naming the failure. The rest of the problem was misclassification, which (1) fixes. Adding the recogniser to shipInfraWait covers a cause that first shows up only in a rerun.

## statement

(1) After the one same-head rerun fails again, rhiza ship check reads the failed jobs' logs (gh run view --log-failed) with a cause recogniser that runs alongside parseFailedTests. Any line it recognises becomes a CI finding, recorded in ci_findings just as a failing test is, so ship check --test-failures (ci-triage) passes and the run goes to ci-fix. Each finding names what failed (the job and step, or the command) and its reason quotes the failing lines. The recognised lines are: rhiza check/self-check 'fail ...' lines (e.g. 'fail render docs/reference/cli.md: differs from what rhiza docs writes'); gofmt/fmt-check file lists; go vet and compiler diagnostics ('file.go:line:col: ...'); and, as the fallback, 'make: *** [target] Error N' naming the target with the lines just before it. A red that repeats after the rerun with none of these lines, or whose log gh cannot give, stays infra and goes to wait-ci-rerun. wait-ci-rerun's attempt check (shipInfraWait) applies the same recogniser, so a recognised failure in a later attempt settles and goes to ci-fix too. ci-fix's prompt widens from 'failing tests' to 'CI findings': reproduce each finding locally with the repo's scoped command (here, make self-check, make fmt-check or go vet on the package), fix it or regenerate the files, and commit once that passes. (2) A language-neutral project setting, rhiza.ship.generate, holds a shell command that regenerates checked-in generated files. It is empty by default, and empty skips. The render steps of the built-in build and review workflows run it before rhiza render and stage and commit whatever it changes, in the same commit and with the same Change trailer. A command that fails fails the step. This repo sets it to a Makefile target (e.g. make docs) that builds the branch's binary in a temp folder and runs rhiza docs. (3) Nothing new for the wait itself: #441's self-rerun with backoff and the 6h end that names the failure stay as they are. Part 3 is met by (1): deterministic failures no longer reach wait-ci-rerun, and the recogniser also runs on its rerun attempts. One testscript replays the 2026-11-11 case: a lint job whose log has 'fail render ...' goes to ci-fix and never to wait-ci-rerun.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-b5d43cfb-1](../question/question-b5d43cfb-1.md)
- derived-from [ci-nontest-failures-plan](../source/ci-nontest-failures-plan.md)

### In

- decided-by [ship-check-ci-cause-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-ci-cause-findings.md)
- decided-by [ci-fix-fixes-ci-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-fix-fixes-ci-findings.md)
- decided-by [render-runs-generate](../capability/rhiza/orchestration/shipping/wait-green-pr/render-runs-generate.md)

## History

- `01M76B50FYAFEWQS29QZCGMTPN` Ship-change treats a deterministic non-test CI failure (lint, self-check, stale generated files) as infra and waits for a CI rerun forever instead of sending it to ci-fix. Evidence 2026-11-11: run 01M75WRZXH7V3G0CE1XCH5Y9F2 (change set 01M75T742VGKA5XMRRXV4W6ZKP, load and memory hold, PR #442 at b70e52ff): CI's lint job failed (make lint: go vet passed, then self-check printed 'failed': 'fail render docs/reference/cli.md: differs from what rhiza docs writes' and the same for config.md, because the builders added rhiza.daemon.maxLoad and minMemory and their flags but never ran rhiza docs). ci-triage ended 'no failing test was recorded for the head: an infra or unparsable failure', the run went to wait-ci-rerun at about 21:15Z and sat there 45+ minutes (a rerun only repeats a deterministic failure). The orchestrator reproduced it locally with make self-check in the run worktree, regenerated the docs with the branch's build, pushed aaafc3a5 and updated the ship state's pushed head. Wanted: (1) ci-triage recognises lint, vet, gofmt and self-check/rhiza check failures in the job log (e.g. 'fail render ...', 'gofmt', vet diagnostics, 'make: *** [...] Error') as fixable findings and routes them to ci-fix with the failing lines, the same as failing tests; only failures with no recognisable cause are treated as infra and rerun; (2) the build and review steps run rhiza docs (or the project's generated-file check, see milestone 32's language-neutral settings) when a change adds or changes settings, flags or commands, so stale generated reference never reaches CI; (3) wait-ci-rerun should not wait indefinitely when no rerun is coming (it should rerun itself within its budget or fail to the inbox). Related: 01M6MB678MD7N13K69C2KDEQ38 (ci-fix for test failures).
