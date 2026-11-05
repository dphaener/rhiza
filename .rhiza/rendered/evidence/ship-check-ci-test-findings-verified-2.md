# ship-check-ci-test-findings-verified-2

evidence · `01M6MTW6GCYW9JT9EQSA6MC8GT` · actor [claude-code-9e42de90](../actor/claude-code-9e42de90.md)

## result

PASS

## statement

At 4c73aede (attempt 2, fixing quality f1) read the diff: parseFailedTests now drops any name another '--- FAIL:' line extends with name+'/', keeping only leaf tests; TestParseFailedTests' nested case now expects only TestScripts/review-check-worktree (two-tests and infra cases unchanged), and run-ship-ci-fix's fake gh log names TestScripts/x and TestScripts/x/sub and its fixer asserts the ship check JSON holds TestScripts/x/sub and not the bare parent. Ran go test -trimpath ./internal/api/... (ok) and go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship|ship-check-ci-test-findings)' (ok, 204s), which covers the run-ship* scripts and the leaf's ship-check-ci-test-findings script that the acceptance's 'TestScripts/run-ship' pattern alone would skip (compliance nit f1).

## Links

### Out

- verifies [ship-check-ci-test-findings](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-ci-test-findings.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
