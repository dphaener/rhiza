# ship-check-failure-names-reruns-and-rhiza-only-diff-verified-2

evidence · `01M6XMNCZ2YTFPM6ZPXY15RE4H` · actor [claude-code-3ccf796d](../actor/claude-code-3ccf796d.md)

## result

PASS

## statement

Verifier, attempt 2 at 648f0eeb (fix for model finding f1): read the diff, which adds one ops entry to change set 01M6WHT and no code. rhiza show confirms that leaf 01M6WJ6ZEEPYS69MVFGK79HHKN now has the handle ship-check-failure-names-reruns-and-rhiza-only-diff with its ID, statement, acceptance and links unchanged (Decision ship-check-leaf-handle-renamed), and that Source plan-ci-flake-reruns-interview's summary now says the error appends that the PR changes only .rhiza/ (a fact, not a claim about the cause), with no 'says the failure is unrelated' (Decision plan-ci-flake-source-summary-matches). Questions review-01m6whtv-r3-f1 and f2 are answered by those Decisions. Acceptance re-run: go test -trimpath ./internal/api/ -run Ship passed (TestShipCheckFailureDetailNotes cases rhiza_only 'failed: test after 2 reruns; failing tests: TestGarden; the PR changes only .rhiza/', other_path with no note, and reads_rhiza 'failed: lint after 2 reruns; the PR changes only .rhiza/', checked for no 'unrelated', 'flaky' or 'on main'; the human line 'PR #7 at abc1234: failed: test, ci/x' assertion passed). go test -trimpath ./cmd/rhiza -run TestScripts/ship-check-ci-test-findings passed, and lines 39 and 51 still read stdout 'failed: test'. The full ./internal/api/... run failed only on TestReadBudget, a timing budget, at load average about 11, plus the known state-folder leak into the enclosing root. Neither is related to this commit. rhiza check reports the same stale-render failures (actor files and index.md from landed actor records) at HEAD~1 as at HEAD, and no problem for 01M6WHT.

## Links

### Out

- verifies [ship-check-failure-names-reruns-and-rhiza-only-diff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-failure-names-reruns-and-rhiza-only-diff.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
