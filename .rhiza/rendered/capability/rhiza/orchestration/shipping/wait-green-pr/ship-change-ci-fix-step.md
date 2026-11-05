# ship-change-ci-fix-step

capability · `01M6MC56XH2VTXCYDT1G5WD55E` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A Go test in internal/workflow loads the built-in ship-change workflow and asserts pr-green and rebased-green fail to ci-triage, ci-triage is a command step running rhiza ship check --test-failures with on.pass ci-fix and no on.fail, ci-fix is an agent step passing to render-pr with max_visits 2, ci-triage max_visits 3, max_visits 5 on render-pr, open-pr, wait-pr, pr-green and signoff, max_visits 3 on rebase, wait-rebased and rebased-green, and no on.fail on landed-green or rebased-merge-green; a Go test in internal/api shows ship check --test-failures passing with a recorded test failure and failing with none or an infra-only failure; a testscript run-ship-ci-fix with a fake gh and a fake fixer agent goes red with a parsed test failure, through ci-triage and one ci-fix commit carrying a Review: trailer, back to green and signoff, a case red at rebased-green goes through ci-fix, signoff and rebase again to land, a case red after two fixes fails the run with the failing tests in the step output, and an infra-only failure fails at ci-triage with no ci-fix visit; run with go test -trimpath ./internal/workflow/... ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship' only

## statement

The built-in ship-change workflow sends a failure at pr-green or rebased-green (on.fail) to a new command step ci-triage, which runs `rhiza ship check --test-failures`: it passes only when the head's latest ship check recorded at least one parsed test failure, and fails otherwise, and ci-triage has no on.fail, so an infra-only or unparsable failure ends the run there as today. ci-triage passes to ci-fix, an agent step (kind = "agent", one fixer session on the run worktree) whose prompt is the review fixer's prompt told to read the failing tests from `rhiza ship check --json` (its test failures and their no-leaf findings), re-run only those tests/scripts and their packages, and make one new commit with a Review: <run> trailer before `rhiza session done`; ci-fix passes to render-pr, then open-pr, wait-pr, pr-green and signoff. ci-fix runs at most twice per run (max_visits 2, so a third red head with test failures goes over the cap and fails the run, pr-green's or rebased-green's output naming the failing tests); ci-triage has max_visits 3; render-pr, open-pr, wait-pr, pr-green and signoff are raised to max_visits 5 so two CI fixes plus the owner's two request-changes fit; rebase, wait-rebased and rebased-green are raised to max_visits 3 so a run re-enters them after each of two CI fixes routed from rebased-green; landed-green and rebased-merge-green still fail the run as today

## Links

### Out

- depends-on [ship-check-ci-test-findings](ship-check-ci-test-findings.md)
- decided-by [ship-change-ci-fix-cycle](../../../../../decision/ship-change-ci-fix-cycle.md)

### In

- verifies [ship-change-ci-fix-step-verified](../../../../../evidence/ship-change-ci-fix-step-verified.md)
- verifies [ship-change-ci-fix-step-qa-1](../../../../../evidence/ship-change-ci-fix-step-qa-1.md)
- verifies [pr-280-ship-change-fixes-red-pr-ci](../../../../../evidence/pr-280-ship-change-fixes-red-pr-ci.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
