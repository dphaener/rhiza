# rhiza-own-workflows-go-cleanup-verified

evidence · `01M7790K1M8S78DA4W9EJ7KP56` · actor [claude-code-c32560f3](../actor/claude-code-c32560f3.md)

## result

PASS

## statement

At 74a8ef908 with a -trimpath build: `rhiza workflow check` passes for bugfix, plan-bug, docs, docs-build and docs-review. Each has [env] GOTMPDIR = "{run_tmp}" (in bugfix and plan-bug it now sits below autonomous=true) and a [[cleanup]] go-cleanup step running .rhiza/workflows/go-cleanup.sh, and `rhiza workflow describe bugfix` lists it. I ran go-cleanup.sh with GOCACHE and TMPDIR pointed at a fake tree in a scratch dir. It deleted only the stale entries: cache/0a/old-a and cache/ff/old-d (4h old), tmp/go-build111 (3h) and tmp/rhiza-old (3h). It kept cache entries 2h old or new, cache/README, trim.txt, a 4h-old file under the non-hex cache/zz, go-build dirs 90m and 30m old, rhiza-new (90m), a non-matching 3h-old dir and a rhiza-* plain file. Exit 0. make self-check in a clean worktree at the commit: rhiza check has no model fails. Its only fails are stale .rhiza/rendered files, the index and a missing actor page for the previous leaf's builder. The builder rendered with untracked logs present, so the index also lists actor 97f4cf59, which is not in the commit. The build workflow's render step regenerates and commits these after the leaves.

## Links

### Out

- verifies [rhiza-own-workflows-go-cleanup](../capability/rhiza/workflows/workflow-overrides/rhiza-own-workflows-go-cleanup.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
