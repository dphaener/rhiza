# rebase-keeps-branch-merges-verified

evidence · `01M68CYQMD13NN1A4V4S5F1XM4` · actor [claude-code-cc6ff980](../actor/claude-code-cc6ff980.md)

## result

PASS

## statement

At cfc2ae4 I ran go test ./cmd/rhiza -run TestScripts/run-ship-change-rebase -count=1 and it passed. Its case 3 has a builder that merges origin's main, which had added a conflicting feature.txt, and resolves the conflict; the run pushes that merge, then main moves with an unrelated commit. The script checks that the run succeeds after the rebase and rebase-merge steps, that no question is raised and the ls output has no rebase-conflict, that the landed head holds main's new commit and has the builder's merge as an ancestor, and that merge passes. Case 4 is the same except main's new commit changes feature.txt again: rebase-merge fails and its status reads exit 1, question rebase-conflict-ship-it is raised naming feature.txt, gh pr merge is never called, and the checkout is left with no merge in progress. I read the diff: runShipRebase counts merges in onto..HEAD and uses the new shipMergeOnto when there are any. shipMergeOnto settles rendered-only conflicts by rendering again and aborts on any other conflict. go test ./internal/api passed. One make check failed in the test step, a full go test ./... then passed, and a second make check exited 0.

## Links

### Out

- verifies [rebase-keeps-branch-merges](../capability/rhiza/workflows/workflow-runs/rebase-when-behind/rebase-keeps-branch-merges.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing
