# ship-change-ci-fix-cycle

decision · `01M6MBZWS89JM6422AY2D9Z5CP`

## reasoning

Read the change set's intent, milestone-25 and milestone-24, internal/workflow/builtin/ship-change.toml and review.toml, and the ship-check-reruns-failed-once evidence. The intent is about the daemon finishing runs without the owner after #259 moved the full suite to CI, which is milestone-25's done state; milestone-24 is about the VM host, not run outcomes. Reusing review's fixer prompt, Review: trailer and 2-cycle cap keeps one fix vocabulary across the run, and keeping the existing once-per-head rerun first means flakes never reach a fixer. Only test failures parsed from --log-failed become findings, so infra failures are not 'fixed' by an agent guessing. Routing the fix through render-pr reuses the existing push/edit/wait/check path, which already handles re-entry (open-pr edits instead of creating). pr-green and rebased-green are before land, so that path is safe; landed-green and rebased-merge-green come after rhiza ship land, and sending them back through render-pr would re-run land and touch a landed change set, so they stay failing and are left to a follow-up rather than widened here. The PR #260 case (pr-green failing on TestScripts/review-check-worktree) is covered.

## statement

(1) Intent: change set 01M6MBPT serves milestone-25 (the daemon empties the backlog on its own); a run that dies on a red PR CI needs the owner, which is what milestone-25 removes. (2) Mechanism: rhiza ship check keeps its one rerun of failed jobs per head for flakes. When the head's checks still fail and the failed Actions job's log (gh run view <run> --log-failed) names failing tests (Go '--- FAIL: TestX[/sub]' lines, testscript names included), ship check reports them as test failures and turns each failing test or script into a finding with no leaf; a failure with no Actions run id, no parsable failing test, or an infra/setup error stays a plain failure that ends the run as today. Because a command step has only pass and fail, the split is a step: pr-green and rebased-green fail to ci-triage (rhiza ship check --test-failures), which passes only when test failures were recorded and otherwise fails the run. ci-fix is an agent step: one fixer session on the run worktree with the review fixer's prompt, told to read the failing tests from rhiza ship check --json and re-run only the failing tests/scripts and their packages, making one new commit with a Review: <run> trailer (Change:/Actor: as usual). (3) Stages: ci-fix goes on pass to render-pr, so the fix is rendered, pushed by open-pr (gh pr edit, no second PR), waited on by wait-pr and checked by pr-green, then signoff again (the autonomous rule passes it; an owner sees the new commit), and from there rebase, wait-rebased and rebased-green again. landed-green and rebased-merge-green keep failing the run as today: they run after rhiza ship land, and re-entering render-pr would land twice; a post-land fix path is a separate change set, not this one. (4) Cap: at most 2 CI-fix cycles per run, counted across pr-green and rebased-green (ci-fix max_visits = 2, matching review's 2 fix cycles); a red CI with test failures after the second fix fails the run, with the parsed failing tests in the step output. Raise max_visits to 5 on render-pr, open-pr, wait-pr, pr-green and signoff so 2 CI fixes plus the owner's request-changes still fit, and to 3 on rebase, wait-rebased and rebased-green (1 today) so the rebased-green fix path re-enters them; ci-triage gets 3.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-5bc7fdce-1](../question/question-5bc7fdce-1.md)
- derived-from [ship-ci-fix-plan](../source/ship-ci-fix-plan.md)

### In

- decided-by [ship-check-ci-test-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-ci-test-findings.md)
- decided-by [ship-change-ci-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-ci-fix-step.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
