# ship-check-names-unrelated-failure-verified

evidence · `01M6XFCP155K214BAGZHV49CBR` · actor [claude-code-2f5edecf](../actor/claude-code-2f5edecf.md)

## result

PASS

## statement

Verifier at fce29b5a: read the diff (internal/api/ship.go adds shipRerunsSpent, inserting ' after <n> reruns' after the failed checks when reruns were spent on the head, and shipDiffOnlyRhiza, appending '; the PR changes only .rhiza/' when git merge-base origin/main..head diff --name-only lists only .rhiza/ paths, any git error or empty diff giving no note). go test -trimpath -count=1 ./internal/api/ -run TestShipCheckFailureDetailNotes -v: PASS for all 3 cases (fake gh with git answered through the same fake Exec, as every ship test does): .rhiza-only gives 'failed: test after 2 reruns; failing tests: TestGarden; the PR changes only .rhiza/'; with main.go too it gives the same without the note; a lint check whose log reports a stale .rhiza/rendered/ file gives 'failed: lint after 2 reruns; the PR changes only .rhiza/'. Each case asserts no 'unrelated', 'flaky' or 'on main'. go test -run Ship ./internal/api/: ok, and ship_test.go:254's 'PR #7 at abc1234: failed: test, ci/x' is unchanged (the diff only adds to ship_test.go). go test -trimpath ./cmd/rhiza -run TestScripts/ship-check-ci-test-findings: ok, and lines 39 and 51 still read stdout 'failed: test'. The whole ./internal/api/... package went over its timeout at load average ~15 in TestReadBudget/TestInboxGateQuestionsBudget and the default 10m limit: those are timing tests this diff does not touch, and CI's make check covers them

## Links

### Out

- verifies [ship-check-failure-names-reruns-and-rhiza-only-diff](../capability/rhiza/orchestration/shipping/wait-green-pr/ship-check-failure-names-reruns-and-rhiza-only-diff.md)

## History

- `01M6WHTZC20SXN9ETH9ENEZ1PN` Fix The plan workflow ends a run on a CI flake: its pr-green and rebased-green steps fail the run as soon as a check is red, with no rerun. Evidence: nebari plan run 01M6VY5APJP9 (PR #78) passed review, gate, pr-green and one rebase, then rebased-green failed at 18230d2 on 'test (4/5); failing tests: rspec ./spec/system/garden_location_spec.rb[1:1:11]', a spec that fails at random and passes on rerun (nebari runs 42427557598, 42427449690). Unstuck by gh run rerun --failed plus rhiza run resume --retry. A plan PR changes only .rhiza/, so a red app test can't be its fault: on red, the plan workflow should rerun the failed jobs once (budgeted) and only then fail, naming the test as unrelated to the plan. Ship-change has ci-triage for this (#350); plan has none.
