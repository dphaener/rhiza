# ship-main-tree-check-verified

evidence · `01M6KZXX6HN71P6D5YCTM3NQ6G` · actor [claude-code-4c99bec9](../actor/claude-code-4c99bec9.md)

## result

PASS

## statement

Verifier claude-code-4c99bec9 at c043fd87 read the diff (internal/api/ship.go, ship-change.toml, tests) and ran go test ./internal/api -run Ship, the testscripts run-ship-change, run-ship-merge-main-moved, run-ship-merge-tree-differs and describe-golden, and make test (exit 0). run-ship-change (local bare origin, fake gh): merge passes to wait-main/main-green and close-out with no gh 'run ' call in gh.log. run-ship-merge-main-moved: origin/main pushed after the head's checks; run status shows merge visit 1 fail -> rebase-merge, rebase-merge visit 2, rebased-merge-green visit 2, merge visit 2 pass; exactly one 'pr merge' call, pinned to the rebased head, and the merge commit holds the other change set's commit. run-ship-merge-tree-differs: fake gh squashes onto a different tree; the run stops with 'stop: merge commit <sha> has tree <t>, but the tested head <sha> has tree <t>' and never reaches close-out. ship check --main now compares merge commit and head trees with git rev-parse ^{tree} after fetch, and no step waits on CI for main. make check's self-check failed only on stale rendered actor files from the earlier leaf commit 6a7e233b, which the run's render-pr step renders, so this leaf did not cause it.

## Links

### Out

- verifies [ship-main-tree-check](../capability/rhiza/workflows/workflow-runs/ship-main-tree-check.md)

## History

- `01M6HW7MG0ZVWWHYP903FXV9YF` Update github runner to: Cancel stale runs, and stop re-testing on main
