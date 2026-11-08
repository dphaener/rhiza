# ship-change-ci-fix-step

capability · `01M6MC56XH2VTXCYDT1G5WD55E` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A Go test in internal/workflow loads the built-in ship-change workflow and asserts pr-green and rebased-green fail to conflict-triage, conflict-triage is a command step running rhiza ship check --conflicts with on.fail ci-triage, ci-triage is a command step running rhiza ship check --test-failures with on.pass ci-fix and no on.fail, ci-fix is an agent step passing to render-pr with max_visits 2, ci-triage max_visits 3, max_visits 10 on render-pr, open-pr, wait-pr and pr-green, max_visits 8 on signoff, max_visits 6 on rebase, wait-rebased and rebased-green, and landed-green, rebase-merge and rebased-merge-green fail to conflict-triage-landed, which has no on.fail; a Go test in internal/api shows ship check --test-failures passing with a recorded test failure and failing with none or an infra-only failure; a testscript run-ship-ci-fix with a fake gh and a fake fixer agent goes red with a parsed test failure, through conflict-triage, ci-triage and one ci-fix commit carrying a Review: trailer, back to green and signoff, a case red at rebased-green goes through ci-fix, signoff and rebase again to land, a case red after two fixes fails the run with the failing tests in the step output, and an infra-only failure fails at ci-triage with no ci-fix visit; run with go test -trimpath ./internal/workflow/... ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship' only

## statement

The built-in ship-change workflow sends a failure at pr-green or rebased-green (on.fail) to conflict-triage, which runs `rhiza ship check --conflicts` and falls back (its on.fail) to a command step ci-triage on anything but a merge conflict, so red CI reaches ci-triage through conflict-triage. ci-triage runs `rhiza ship check --test-failures`: it passes only when the head's latest ship check recorded at least one parsed test failure, and fails otherwise, and ci-triage has no on.fail, so an infra-only or unparsable failure ends the run there as today. ci-triage passes to ci-fix, an agent step (kind = "agent", one fixer session on the run worktree) whose prompt is the review fixer's prompt told to read the failing tests from `rhiza ship check --json` (its test failures and their no-leaf findings), re-run only those tests/scripts and their packages, and make one new commit with a Review: <run> trailer before `rhiza session done`; ci-fix passes to render-pr, then open-pr, wait-pr, pr-green and signoff. ci-fix runs at most twice per run (max_visits 2, so a third red head with test failures goes over the cap and fails the run, pr-green's or rebased-green's output naming the failing tests); ci-triage has max_visits 3; render-pr, open-pr, wait-pr and pr-green have max_visits 10 and signoff max_visits 8, so two CI fixes, conflict fixes and the owner's two request-changes fit; rebase, wait-rebased and rebased-green have max_visits 6 so a run re-enters them after each CI or conflict fix routed from rebased-green. After land, a failure at landed-green, rebase-merge or rebased-merge-green goes to conflict-triage-landed, which has no on.fail, so a red CI or any other non-conflict failure still ends the run there

## Links

### Out

- depends-on [ship-check-ci-test-findings](ship-check-ci-test-findings.md)
- decided-by [ship-change-ci-fix-cycle](../../../../../decision/ship-change-ci-fix-cycle.md)

### In

- verifies [ship-change-ci-fix-step-verified](../../../../../evidence/ship-change-ci-fix-step-verified.md)
- verifies [ship-change-ci-fix-step-qa-1](../../../../../evidence/ship-change-ci-fix-step-qa-1.md)
- verifies [pr-280-ship-change-fixes-red-pr-ci](../../../../../evidence/pr-280-ship-change-fixes-red-pr-ci.md)
- verifies [pr-383-fix-ship-change-stops-on-a](../../../../../evidence/pr-383-fix-ship-change-stops-on-a.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
- `01M6WMCCRC0KSCWEV17VY8KVQP` Fix Ship-change stops on a rebase conflict and waits for the owner instead of resolving it: run 01M6W0Z45JXE (change set 01M6VY7G6HPJ, Workflow dispatch, PR #367) stopped at wait-pr on 2026-11-08 ~01:20Z with 'rebasing cs/workflow-dispatch onto origin/main conflicts in describe-golden.txtar, internal/api/change.go, internal/api/run.go, internal/changeset/body.go, body_test.go, internal/daemon/requests.go, internal/store/read_test.go; the rebase was aborted; answer question rebase-conflict-workflow-dispatch'. With several runs in flight, main always moves under a long run, so this will recur. Instead, ship-change should merge origin/main into the branch (never rebase or force-push) and give the conflicts to a fixer agent session with both sides' commits and change sets as context. The fixer resolves them, regenerates golden files, runs the touched packages' tests, commits with trailers and pushes; the run records the push and goes on. It asks the owner (or the arbiter in an autonomous run) only if the fixer reports a real design clash. Unstuck by hand: an agent merged main in the run's worktree, answered the question and resumed the run.
