# ci-unstable-gate-rules

decision · `01M6YCNMKDZAYW484HXMT7JZJH`

## reasoning

Read the change set's intent (owner's words of 2026-11-08), milestone-25 and its served change sets, and the code: watchdog_policy.go already has ciReruns=1 and fileBugs, daemon_watchdog.go reruns failed jobs once per head, and no local-gate setting exists yet. milestone-25 is the only milestone the CI work already serves and the daemon builds only milestone change sets. A rolling window with flakes and main's red as the trigger matches the intent's 'red runs that pass on rerun, per test' and 'while main's CI is not reliably green', and clears on its own so the default (#259, scoped tests plus PR CI) returns when CI is green. The gate as a configured command keeps Rhiza repo-agnostic; make check is what CI runs here, so it is the gate that predicts green. Holding a run whose head is green over an unrelated flake would stall the backlog the milestone empties; the intent's local-reproduce rule targets reds that stay unexplained, so the flake is tracked and filed as a bug while the stay-red case gets the local fix with no second rerun.

## statement

(1) The change set serves milestone-25 (the backlog is empty: the daemon builds and ships on its own), as its sibling 'Ship-change fixes red PR CI' does. (2) A repo's CI counts as unstable while, within its last 20 completed CI runs (main and run PR heads; the window size is a project setting), any run went red and then passed on a rerun of the same head (a flake), or while main's latest completed CI run is red. It counts as stable again once the window holds no flake and main's latest run is green. Main's latest run is the result the latest rhiza ship check --main recorded, since CI does not run on a push to main: a merge whose tree equals its tested PR head's takes that head's green result, and a tree-differing merge takes its dispatched CI run's result. rhiza health shows the state with the flaky tests and their counts. (3) The 'full local gate' is one command the repo configures as a project setting (for Rhiza itself: make check). While CI is unstable, ship-change runs it in the run's worktree before every push and pushes only if it exits 0; on failure it fixes locally like any red. If no gate is configured, the run falls back to the scoped tests and rhiza health warns that CI is unstable with no local gate set. (4) A red PR that passes on its one rerun is a flake: the run is not held, it proceeds on the green head and lands as usual. The flake is recorded against the failing test in the CI stability record (which keeps CI unstable), and, with the watchdog's fileBugs on, a bug is filed once per test (deduplicated) to reproduce it locally and fix or quarantine it. A head that is still red after its one rerun is not rerun again: the run reproduces the failing test locally (repeatedly), fixes or quarantines it with a bug filed, and pushes again only after that test and the touched packages pass locally.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-25f669da-1](../question/question-25f669da-1.md)
- derived-from [ci-debugger-plan-interview](../source/ci-debugger-plan-interview.md)

### In

- decided-by [ci-stability-record](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-record.md)
- decided-by [ci-stability-state](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-stability-state.md)
- decided-by [health-ci-stability](../capability/rhiza/orchestration/daemon/daemon-status/health-command/health-ci-stability.md)
- decided-by [ship-local-gate](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-local-gate.md)
- decided-by [ci-fix-reproduces-locally](../capability/rhiza/orchestration/shipping/wait-green-pr/ci-fix-reproduces-locally.md)
- decided-by [flaky-test-bugs](../capability/rhiza/orchestration/daemon/daemon-watchdog/flaky-test-bugs.md)

## History

- `01M6YCDW1VM9R8KZPGQR330MAY` Stop using CI as a debugger while CI is unstable. Owner, 2026-11-08: "the amount of failed CI runs is way out of control. I'm not sure howthis is happening but its got to stop. if its not green here IT should never go too Ci" and then: "i think my note about ci didnt have enough nuance. the original request was meant to speed up builds. this is assuming we already have a stable, green ci pipeline. which we currently do not. when we're diagnosing flaky tests just firing away at ci A/m's hopping for the best doesn't make sense" Scoped local tests plus the PR's CI (#259) stay the default when CI is stable and green. But when a CI failure is not explained by the change (a flaky test, an infra error, a red that passes on rerun), or while main's CI is not reliably green, runs must reproduce and fix it locally instead of pushing and rerunning CI in hope: at most one CI rerun per head for infrastructure, then the failing test runs locally (repeatedly, to reproduce the flake) and is fixed or quarantined with a bug filed, and the run pushes again only after that test and the touched packages pass locally. Track CI stability (red runs that pass on rerun, per test) and show it in rhiza health; while it is unstable, ship-change runs the repo's full local gate before pushing. Record the owner's words as a decision refining #259, not replacing it.
