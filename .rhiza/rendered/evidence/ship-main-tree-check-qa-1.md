# ship-main-tree-check-qa-1

evidence · `01M6M052CN10V9Z15F5JJAQJ99` · actor [claude-code-ca676fb4](../actor/claude-code-ca676fb4.md)

## result

PASS

## statement

At 83445164, with the built binary in a scratch repo with a bare local origin, real git and a logging fake gh, and a hand-written .git/rhiza/ship state: rhiza ship check --main --settled on a squash merge commit whose tree equals the recorded head's printed 'main at <sha>: the merge commit's tree is the tested head's' and exited 0; with a merge commit of a different tree it printed 'stop: merge commit 7cf4b0b has tree b9f5290, but the tested head 5b43c36 has tree 8e3bca3; main is not what the checks passed on' and exited 75; with no head recorded it exited 2 naming rhiza ship check. The fake gh logged no call for any --main check. rhiza ship merge with origin/main moved past the head printed 'Not merged: origin/main has moved past PR #7's checked head 5b43c36; rebase and check again.', exited 1, made no gh call and recorded no merge commit; with a head holding origin/main it called gh pr merge 7 --squash --match-head-commit <head> and exited 0. go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-change|run-ship-merge-main-moved|run-ship-merge-tree-differs)$' passed (the full ship-change run back through rebase-merge on a moved main, and the stop on a differing tree).

## Links

### Out

- verifies [ship-main-tree-check](../capability/rhiza/workflows/workflow-runs/ship-main-tree-check.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
