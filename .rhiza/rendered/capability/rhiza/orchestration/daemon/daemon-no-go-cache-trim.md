# daemon-no-go-cache-trim

capability · `01M7642MDEKRSQERY6KR37594J` · parent [daemon](../daemon.md)

## acceptance

The daemon-settings testscript lists no trimGoCache; internal/api daemon tests show no Go cache sweep; the generated reference is regenerated.

## statement

The daemon no longer trims the Go build cache: rhiza.daemon.trimGoCache is removed from daemon settings and the reference, and the daemon no longer calls gosweep (daemon-trims-go-cache and go-build-cache-sweep are retired).

## Links

### In

- verifies [daemon-no-go-cache-trim-verified](../../../../evidence/daemon-no-go-cache-trim-verified.md)
- verifies [daemon-no-go-cache-trim-qa-1](../../../../evidence/daemon-no-go-cache-trim-qa-1.md)
- verifies [daemon-no-go-cache-trim-verified-2](../../../../evidence/daemon-no-go-cache-trim-verified-2.md)
- verifies [pr-475-rhiza-s-go-behaviour-moves](../../../../evidence/pr-475-rhiza-s-go-behaviour-moves.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
