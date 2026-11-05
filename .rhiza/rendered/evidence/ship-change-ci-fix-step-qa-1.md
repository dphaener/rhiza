# ship-change-ci-fix-step-qa-1

evidence · `01M6MS4EEVQ47NDX7622XWKK0R` · actor [claude-code-7f0b608d](../actor/claude-code-7f0b608d.md)

## result

PASS

## statement

QA at dc3c5d5c with a built binary: rhiza workflow check ship-change says valid; rhiza workflow describe ship-change shows pr-green and rebased-green on.fail -> ci-triage, ci-triage (command) on.pass -> ci-fix (agent), and landed-green and rebased-merge-green with no on.fail. In a scratch repo with a fake gh, rhiza ship check --test-failures exited 1 ('no failing test was recorded for the head: an infra or unparsable failure') before any check, after only the rerun, after an infra-only log and after green; it exited 0 after a red check whose log named failing tests, with --json listing tests and no-leaf findings and the same bytes as rhiza api call ship.check '{"test_failures":true}'; gh.log shows it ran no gh. --test-failures with --settled is refused with exit 2. go test -trimpath ./internal/workflow/... ok and go test -trimpath ./cmd/rhiza -run 'TestScripts/run-ship' ok (run-ship-ci-fix: fix at pr-green, fix at rebased-green then land, failure after two fixes, infra failure at ci-triage with no ci-fix). The agent step itself ran only against the testscript's fake agent, not a real one

## Links

### Out

- verifies [ship-change-ci-fix-step](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-change-ci-fix-step.md)

## History

- `01M6MBPT20RGKTT334GQK997B3` Ship-change fixes red PR CI instead of failing: now that agents run only scoped tests (#259), the PR's GitHub CI is where the rest of the suite runs, so a failing test there must go back to a fixer, not end the run. When pr-green (or rebased-green / landed-green) finds the PR's checks failed with test failures (not infra or flakes already retried), read the failed job log (gh run view --log-failed), turn each failing test or script into a finding, run one fixer cycle on the run worktree with those findings (scoped re-tests of what failed), push, and wait for CI again, within a cap; only then fail. Seen 2026-11-04: actor-dedupe PR #260 failed CI on TestScripts/review-check-worktree and its run failed at pr-green
