# rhiza-own-workflows-go-cleanup-qa-1

evidence · `01M77AP4Y4VF0NNTV2EW35B39P` · actor [claude-code-174ed9c7](../actor/claude-code-174ed9c7.md)

## result

PASS

## statement

QA at 46f1cad8 with a -trimpath build: rhiza workflow check passes for bugfix, plan-bug, docs, docs-build, docs-review and ship-change, and workflow describe bugfix and ship-change list the go-cleanup cleanup step. Ran .rhiza/workflows/go-cleanup.sh with GOCACHE and TMPDIR on a fake tree: it removed cache/0a/old-a (4h), tmp/go-build-old (3h), tmp/rhiza-old (3h) and tmp/rhiza-new (90m, over the hour); it kept cache/0a/mid-b (2h), cache/ff/new-c, cache/README and trim.txt (4h, top level), cache/zz/old-z (non-hex dir), tmp/go-build-mid (90m), tmp/other-old (3h, other prefix) and a 3h-old rhiza-file plain file; exit 0. With GOCACHE=off a 4h-old cache entry was left. In a scratch repo a workflow's [[cleanup]] step sees the [env] GOTMPDIR folder.

## Links

### Out

- verifies [rhiza-own-workflows-go-cleanup](../capability/rhiza/workflows/workflow-overrides/rhiza-own-workflows-go-cleanup.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
