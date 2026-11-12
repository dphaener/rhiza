# ship-head-build-step

capability · `01M7642MDWX73CXGPSKBXQCJK4` · parent [shipping](../shipping.md)

## acceptance

internal/api ship tests show a repo whose ship-change names a head build renders with it and one that names none renders in-process, with isRhizaRepo gone from ship; TestShipRebaseBuildsMergedHead and its related tests show ship rebase merges origin/main, then runs SHIP_HEAD_BUILD, then renders, and that a failed build pushes nothing and resets the branch; the ship-land-branch-build testscript shows ship land landing and rendering with the SHIP_HEAD_RHIZA build and failing, with nothing landed or pushed, when that build is missing; Rhiza's workflows, bugfix included, set SHIP_HEAD_BUILD and pass rhiza workflow check.

## statement

Ship no longer checks whether the repo is Rhiza's own. A workflow names a build of the PR head in its [env]: SHIP_HEAD_RHIZA, the binary ship's render and land run with, and SHIP_HEAD_BUILD, the shell command that builds it. Rhiza's own workflows (ship-change, docs and bugfix, which includes ship-change) set both, and their head-build steps run SHIP_HEAD_BUILD before rebase, land and rebase-merge, failing when it is empty. Ship rebase changes the head it renders, so it runs SHIP_HEAD_BUILD itself (sh -c, in the worktree) after merging origin/main and before rendering, both when settling a rendered-only conflict and before its final render; a failed build pushes nothing and resets the branch. A named binary that is missing fails ship before it lands, renders or pushes anything. A workflow that names no head build renders and lands in-process.

## Links

### In

- verifies [ship-head-build-step-verified](../../../../evidence/ship-head-build-step-verified.md)
- verifies [ship-head-build-step-verified-2](../../../../evidence/ship-head-build-step-verified-2.md)
- verifies [ship-head-build-step-qa-2](../../../../evidence/ship-head-build-step-qa-2.md)
- verifies [ship-head-build-step-verified-3](../../../../evidence/ship-head-build-step-verified-3.md)
- verifies [ship-head-build-step-qa-3](../../../../evidence/ship-head-build-step-qa-3.md)
- verifies [pr-475-rhiza-s-go-behaviour-moves](../../../../evidence/pr-475-rhiza-s-go-behaviour-moves.md)

## History

- `01M763YX47S5XY0W3P8FVJ54S8` Rhiza's Go behaviour moves into its own workflows: workflows can set environment variables (Rhiza's set GOTMPDIR to the run's temp folder), Rhiza's cleanup phase trims the Go build cache and stale go-build and rhiza test temp folders, Rhiza's ship-change workflow builds rhiza from the PR head in a step; trimGoCache, gosweep's daemon use, the product's GOTMPDIR and ship's is-Rhiza-repo check leave the product.
