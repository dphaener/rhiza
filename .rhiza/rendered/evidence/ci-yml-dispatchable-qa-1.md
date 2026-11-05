# ci-yml-dispatchable-qa-1

evidence · `01M6MDYYJ0NRTPQXJ00STNJH0B` · actor [claude-code-9aaf5286](../actor/claude-code-9aaf5286.md)

## result

PASS

## statement

QA at 5188b19f: read .github/workflows/ci.yml. Triggers are pull_request (branches [main]) and workflow_dispatch with a required sha input. run-name is 'CI ${{ inputs.sha || github.event.pull_request.title }}', the group is ci-${{ github.event.pull_request.number || inputs.sha }}, and checkout ref is ${{ inputs.sha || github.sha }}. A copy of this file on a scratch origin/main was selected and dispatched by the scratch build's ship check --main ('workflow run ci.yml --ref main -f sha=<merge commit>'), and found by the sha in its run title. go test -trimpath ./internal/ciworkflow/... passed. GitHub's real gh run list shows displayTitle and name fields for this repo's CI runs, the fields ship matches on. Not dispatched on GitHub itself, because this branch is not pushed yet.

## Links

### Out

- verifies [ci-yml-dispatchable](../capability/rhiza/integrity/ci-yml-dispatchable.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
