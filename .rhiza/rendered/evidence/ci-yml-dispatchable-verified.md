# ci-yml-dispatchable-verified

evidence · `01M6MBWTE2FF6WA2WBTFDHZ9PH` · actor [claude-code-637fc130](../actor/claude-code-637fc130.md)

## result

PASS

## statement

At 6c7b06c5 I read the diff: .github/workflows/ci.yml adds workflow_dispatch with a required sha input, run-name 'CI ${{ inputs.sha || github.event.pull_request.title }}', concurrency group ci-${{ github.event.pull_request.number || inputs.sha }}, and checkout ref ${{ inputs.sha || github.sha }}. Ran go test -trimpath -count=1 -v ./internal/ciworkflow/...: TestCIRunsOnPullRequestsAndDispatch, TestCICheckoutAndRunName and TestCICancelsStaleRuns all PASS. In a scratch copy I mutated ci.yml six ways (workflow_dispatch renamed, required: false, ref without inputs.sha, run-name without inputs.sha, group without the inputs.sha fallback, pull_request replaced by push) and each made the test fail. The package sits under go test ./..., which make test and so make check run.

## Links

### Out

- verifies [ci-yml-dispatchable](../capability/rhiza/integrity/ci-yml-dispatchable.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
