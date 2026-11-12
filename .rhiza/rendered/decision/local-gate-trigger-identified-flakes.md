# local-gate-trigger-identified-flakes

decision · `01M78FKX3BQT0ABF59FCKH9B6Z`

## reasoning

Read the question, this change set's intent (three runs 01M787SV, 01M788HB, 01M7859E all failed the ~20 min serialized make check gate on an unrelated root-guard leak, with the gate on only because of one 'scripts (4/4) / UNKNOWN STEP' flake on #498 while 29 of the last 30 CI runs were green), decision ci-unstable-gate-rules (any flake in the window, or main red, makes CI unstable), and the code: ciVerdict in internal/api/ci_stability.go returns unstable on any non-empty Flakes in the window, and shipFailedTests in internal/api/ship.go falls back to 'CI run <id>' when a log names no job/step. The gate exists to predict CI green before a push; a flake nobody can name cannot be reproduced or targeted locally, is as likely an infra or runner error as a test race, and AGENTS.md itself treats an infra error as a red not explained by the change. Letting one such event serialize every run behind a 20 min suite for the next 20 CI runs costs far more than it buys and, as the evidence shows, exposes runs to unrelated local-only failures. Keeping identified flakes and main-red as triggers preserves the original owner intent ('while main's CI is not reliably green', 'red runs that pass on rerun, per test') with the smallest change to the verdict; a count threshold was weighed but rejected as a larger departure from ci-unstable-gate-rules that the evidence does not require. Still recording, surfacing and filing the unidentified flake keeps it visible so a recurring infra problem is not hidden.

## statement

This decision supersedes ci-unstable-gate-rules, restating its rules with clause (2) changed so the full local gate switches on only for instability a local run can predict. (1) A repo's CI counts as unstable while main's latest completed CI result in the stability record is red, or while its last completed CI runs (main and run PR heads; 20 by default, the project setting rhiza.ship.ciWindow) hold an identified flake: a run that went red and then passed on a rerun of the same head, attributed to a named failing test (a '--- FAIL:' test name) or a 'job / step' whose step is known. It counts as stable again once neither holds. Main's latest result is the one the latest rhiza ship check --main recorded, since CI does not run on a push to main: a merge whose tree equals its tested PR head's takes that head's green result, and a tree-differing merge takes its dispatched CI run's result. (2) An unidentified flake (a step read as 'UNKNOWN STEP', or 'CI run <id>' whose log names no job or step) is still recorded in the CI stability record, still shown by rhiza health with its count and still filed as a bug by the watchdog, but it does not count toward the unstable verdict and so does not switch on the gate; the CI verdict's reason and ship.gate's skip reason say when one was seen but ignored. rhiza health shows the state with the flaky tests and their counts. (3) The 'full local gate' is one command the repo configures as a project setting (for Rhiza itself: make check). While CI is unstable, ship-change runs it in the run's worktree before every push and pushes only if it exits 0; on failure it fixes locally like any red. If no gate is configured, the run falls back to the scoped tests and rhiza health warns that CI is unstable with no local gate set. (4) A red PR that passes on its one rerun is a flake: the run is not held, it proceeds on the green head and lands as usual. The flake is recorded against the failing test in the CI stability record, and, with the watchdog's fileBugs on, a bug is filed once per test (deduplicated) to reproduce it locally and fix or quarantine it. A head that is still red after its one rerun is not rerun again: the run reproduces the failing test locally (repeatedly), fixes or quarantines it with a bug filed, and pushes again only after that test and the touched packages pass locally.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-92aa21cf-1](../question/question-92aa21cf-1.md)
- supersedes [ci-unstable-gate-rules](ci-unstable-gate-rules.md)
- derived-from [local-ship-gate-plan](../source/local-ship-gate-plan.md)

### In

- decided-by [ci-stability-state](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-state.md)
- decided-by [ci-stability-record](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-record.md)
- decided-by [health-ci-stability](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-ci-stability.md)
- decided-by [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)
- decided-by [ci-fix-reproduces-locally](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-fix-reproduces-locally.md)
- decided-by [flaky-test-bugs](../capability/rhiza/orchestration/daemon/daemon-watchdog/flaky-test-bugs.md)
- supersedes [local-gate-repeated-named-flakes](local-gate-repeated-named-flakes.md)

## History

- `01M78FD039D0F1QVN2BWEAEYCC` Fix The local ship gate (rhiza.ship.localGate = make check) fails in every run under the VM root: internal/api passes all tests but its leak guard (internal/rootguard) ends 'FAIL: the test run wrote into the state folder of the root enclosing it: .rhiza-root/usage.live.json: file changed', because every live Claude session's status line rewrites ~/code/.rhiza-root/usage.live.json during the ~10-20 min suite. Evidence 2026-11-12 16:30-17:10: runs 01M787SV, 01M788HB and 01M7859E all failed local-gate/gate-rebase-merge at 'make: *** [Makefile:64: test] Error 1' and went to gate-fix; the gate was on only because one 'scripts (4/4) / UNKNOWN STEP' flake on #498 marked CI unstable (29 of the last 30 CI runs green). The gate-fix for 01M7859E fixed the guard inside an unrelated change set (commit e075aefa6 'Judge the live usage file by captured_at in the root guard', PR #498 rhiza release check); the other two fixers were working the same failure in parallel, so expect duplicate fixes or conflicts. Also: the gate output stored in the journal is only the last ~2 KB, which cut off the failing test's name, so every fixer had to rerun the full suite (~20 min each, serialized on the test lock). Fix ideas: land the guard fix on its own; store or point to the full gate log; consider whether one unidentified flake should switch on the full local gate for every run.
