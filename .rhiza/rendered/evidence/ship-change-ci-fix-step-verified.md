# ship-change-ci-fix-step-verified

evidence · `01M6MRJJXHTTPPW88HZH4VAV6R` · actor [claude-code-a65b43ad](../actor/claude-code-a65b43ad.md)

## result

PASS

## statement

Verifier at 03bc97ba read the diff (ship-change.toml: pr-green and rebased-green on.fail = ci-triage; ci-triage runs rhiza ship check --test-failures, on.pass ci-fix, no on.fail, max_visits 3; ci-fix agent step to render-pr, max_visits 2, prompt names Review: {run} and ship check --json; caps 5 on render-pr/open-pr/wait-pr/pr-green/signoff, 3 on rebase/wait-rebased/rebased-green; landed-green and rebased-merge-green keep no on.fail; ship.go --test-failures reads the state file only, and a red check's detail now names the failing tests). Ran go test -trimpath ./internal/workflow/... (ok; TestShipChangeBuiltin asserts every point above), go test -trimpath ./internal/api/... -run TestShip (ok; TestShipCheckTestFailures: no state and infra-only state fail with exit 1, a recorded test failure passes, no gh call, --settled refused), go test -trimpath ./cmd/rhiza -run TestScripts/run-ship (ok, includes run-ship-ci-fix covering: red at pr-green -> ci-triage -> one ci-fix commit with Review: trailer -> green -> signoff; red at rebased-green -> ci-fix -> signoff -> rebase visit 2 -> lands; red after two fixes fails at ci-triage visit 3 with TestScripts/x in status; infra-only fails at ci-triage with no ci-fix), and describe-golden (ok). The full internal/api run failed only TestReadBudget timing budgets at load average 14-16 on 4 cores, which is unrelated to this diff

## Links

### Out

- verifies [ship-change-ci-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-ci-fix-step.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
