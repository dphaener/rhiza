# pr-265-fix-ship-s-new-main-tree-check

evidence · `01M6MGAXP9B01SAZK3M510KNBH` · actor [dphaener](../actor/dphaener.md)

## head

59fd6fc3d6b325d5a87d8c2d22f9f133927abe3f

## pr

https://github.com/dphaener/rhiza/pull/265

## result

1 check passed: check SUCCESS

## statement

Pull request #265 (https://github.com/dphaener/rhiza/pull/265) passed CI at 59fd6fc

## status

verified

## Links

### Out

- verifies [ship-main-tree-check](../capability/rhiza/workflows/workflow-runs/ship-main-tree-check.md)
- verifies [ci-yml-cancels-stale](../capability/rhiza/integrity/ci-yml-cancels-stale.md)
- verifies [ci-install-cancels-stale](../capability/rhiza/integrity/ci-install-cancels-stale.md)
- verifies [ship-main-dispatches-ci](../capability/rhiza/workflows/workflow-runs/ship-main-dispatches-ci.md)
- verifies [ci-yml-dispatchable](../capability/rhiza/integrity/ci-yml-dispatchable.md)
- verifies [ci-install-dispatchable](../capability/rhiza/integrity/ci-install-dispatchable.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
