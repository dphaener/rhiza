# ci-yml-dispatchable

capability · `01M6MAKD15PZYXSTNV3FX2SQCN` · parent [integrity](../integrity.md)

## acceptance

The Go test that pins .github/workflows/ci.yml asserts that its triggers are exactly pull_request and workflow_dispatch, that workflow_dispatch declares a required input sha, that the checkout step's ref is ${{ inputs.sha || github.sha }}, that run-name includes inputs.sha, and that the concurrency group is keyed on github.event.pull_request.number falling back to inputs.sha; make check runs it.

## statement

Rhiza's own CI can be dispatched on a given commit: .github/workflows/ci.yml gains a workflow_dispatch trigger with a required sha input, its checkout step checks out ${{ inputs.sha || github.sha }}, its run-name carries the sha on a dispatch, and its concurrency group falls back to inputs.sha when there is no pull request, so dispatched runs never cancel each other.

## Links

### Out

- decided-by [ship-main-dispatch-ci](../../../decision/ship-main-dispatch-ci.md)
- derived-from [plan-ship-main-dispatch](../../../source/plan-ship-main-dispatch.md)

### In

- verifies [ci-yml-dispatchable-verified](../../../evidence/ci-yml-dispatchable-verified.md)
- verifies [ci-yml-dispatchable-qa-1](../../../evidence/ci-yml-dispatchable-qa-1.md)
- verifies [pr-265-fix-ship-s-new-main-tree-check](../../../evidence/pr-265-fix-ship-s-new-main-tree-check.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
