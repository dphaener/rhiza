# ship-check-failure-names-reruns-and-rhiza-only-diff

capability · `01M6WJ6ZEEPYS69MVFGK79HHKN` · parent [wait-green-pr](../wait-green-pr.md)

## acceptance

A Go test in internal/api with a fake gh and a git repo shows a red head after 2 reruns failing with the detail 'failed: test after 2 reruns; failing tests: <parsed test>; the PR changes only .rhiza/' when the branch's diff touches only .rhiza/ paths, and the same detail without the .rhiza/ note when it also touches a file outside .rhiza/; a .rhiza-only head whose failing check reads .rhiza/ (a check named lint whose log reports a rhiza check problem on .rhiza/rendered/) gets 'failed: lint after 2 reruns; the PR changes only .rhiza/' and no detail in any case contains 'unrelated', 'flaky' or 'on main'; the existing assertions pass unchanged, including internal/api/ship_test.go's human line 'PR #7 at abc1234: failed: test, ci/x' (no rerun spent) and the testscript ship-check-ci-test-findings's stdout 'failed: test' (lines 39 and 51); run with go test -trimpath ./internal/api/... and go test -trimpath ./cmd/rhiza -run 'TestScripts/ship-check-ci-test-findings' only

## statement

When a head's checks stay red, rhiza ship check's failure detail keeps today's form, 'failed: <checks>' then '; failing tests: <tests>' when it parsed any from the failed runs' logs (the human line already prefixes the PR and short head, 'PR #7 at abc1234: '), and extends it in two ways only: when one or more reruns were spent on the head it inserts ' after <n> reruns' right after the failed checks ('failed: test after 2 reruns; failing tests: ...'), and when the head's diff from its merge base with origin/main touches only paths under .rhiza/ it appends '; the PR changes only .rhiza/'. The note states only that fact and never claims the failure's cause: a check that reads .rhiza/ itself (Rhiza's own make check runs rhiza check and fails on stale .rhiza/rendered/ files) can go red on a .rhiza-only head because of the plan, so ship check never says a failure is unrelated to the plan or a flake on main. A head with no rerun spent (for example a check with no Actions run) gets no rerun count, and a diff touching any other path, or one that cannot be computed, gets no .rhiza/ note, so those details are byte-identical to today's

## Links

### Out

- decided-by [plan-ci-flake-reruns](../../../../../decision/plan-ci-flake-reruns.md)

### In

- depends-on [plan-run-reruns-ci-flakes](../../../workflows/built-in-workflows/plan-workflow/plan-run-reruns-ci-flakes.md)
- blocks [review-01m6whtv-r3-f2](../../../../../question/review-01m6whtv-r3-f2.md)
- verifies [ship-check-names-unrelated-failure-verified](../../../../../evidence/ship-check-names-unrelated-failure-verified.md)
- verifies [ship-check-names-unrelated-failure-qa-1](../../../../../evidence/ship-check-names-unrelated-failure-qa-1.md)
- verifies [ship-check-failure-names-reruns-and-rhiza-only-diff-verified-2](../../../../../evidence/ship-check-failure-names-reruns-and-rhiza-only-diff-verified-2.md)
- verifies [pr-386-fix-the-plan-workflow-ends-a](../../../../../evidence/pr-386-fix-the-plan-workflow-ends-a.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
