# ship-main-tree-check-verified-2

evidence · `01M6M2HWN31Q0WZ4GH180C4JZY` · actor [claude-code-a57e89bf](../actor/claude-code-a57e89bf.md)

## result

PASS

## statement

Verifier, attempt 2 at 906a4a57 (review fixer cycle 2). Read the fixer commit: it revises close-out (main confirmed by the merge commit's tree, no CI wait) and repo-ship-command's acceptance (ship command runs after main is confirmed by tree), resolving model findings f1/f2. Re-ran the acceptance testscripts with a local bare origin and fake gh: go test ./cmd/rhiza -run 'TestScripts/(run-ship-change|run-ship-change-rebase|run-ship-merge-main-moved|run-ship-merge-tree-differs)$' -count=1 -> ok. run-ship-change asserts merge then main-green -> close-out with no gh run call; run-ship-merge-main-moved asserts merge visit 1 fails back to rebase-merge, rebased head re-checked, exactly one gh pr merge pinned to the new head, no '^run ' in gh.log; run-ship-merge-tree-differs asserts a stop: line naming merge commit and tested head with their trees, no close-out, no gh run call. Built-in ship-change.toml's wait-main/main-green run rhiza ship check --main (tree check), no CI wait on main. make test: every package ok. make check's self-check fails only on stale rendered files and untracked actor records, which the run's render-pr step commits.

## Links

### Out

- verifies [ship-main-tree-check](../capability/rhiza/workflows/workflow-runs/ship-main-tree-check.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
