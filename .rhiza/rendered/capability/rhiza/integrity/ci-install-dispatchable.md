# ci-install-dispatchable

capability · `01M6MAKF848SVSTKWG667GY87T` · parent [integrity](../integrity.md)

## acceptance

The ci-install testscript's expected workflow keeps the push-to-main and pull_request triggers and their check steps, adds a workflow_dispatch trigger with a required sha input, checks out ${{ inputs.sha || github.sha }}, has a run-name holding inputs.sha on dispatch, keys its concurrency group on the pull request number, then inputs.sha, then github.sha, with cancel-in-progress ${{ github.event_name == 'pull_request' }}, and has a step run only on workflow_dispatch that runs rhiza check --base HEAD~1.

## statement

The workflow rhiza ci install writes can be dispatched on a given commit, so ship can run it on a merge commit: it gains a workflow_dispatch trigger with a required sha input; a dispatched run checks out that sha, carries it in its run-name, runs in a concurrency group keyed on the sha, and checks the commit against its first parent with rhiza check --base HEAD~1, while pull request and push-to-main runs behave as before.

## Links

### Out

- decided-by [ship-main-dispatch-ci](../../../decision/ship-main-dispatch-ci.md)
- derived-from [plan-ship-main-dispatch](../../../source/plan-ship-main-dispatch.md)

### In

- verifies [ci-install-dispatchable-verified](../../../evidence/ci-install-dispatchable-verified.md)
- verifies [ci-install-dispatchable-qa-1](../../../evidence/ci-install-dispatchable-qa-1.md)
- verifies [pr-265-fix-ship-s-new-main-tree-check](../../../evidence/pr-265-fix-ship-s-new-main-tree-check.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
