# ci-install-dispatchable-verified

evidence · `01M6MDA8QF6PRS1DAMRE4TB7K7` · actor [claude-code-2a4a1ece](../actor/claude-code-2a4a1ece.md)

## result

PASS

## statement

At 7f5581c4, verifier claude-code-2a4a1ece read the diff to internal/api/ci.go workflowTemplate and to cmd/rhiza/testdata/script/ci-install.txtar, then parsed the expected rhiza.yml with PyYAML. It still has the push (branches: [main]) and pull_request triggers and their check steps: --base origin/${{ github.base_ref }}, --base ${{ github.event.before }}, and the first-push check. It adds workflow_dispatch with inputs.sha required: true, and checkout ref ${{ inputs.sha || github.sha }} (fetch-depth 0). run-name is 'Rhiza ${{ inputs.sha || github.event.pull_request.title || github.sha }}'. The concurrency group is rhiza-${{ github.event.pull_request.number || inputs.sha || github.sha }}, with cancel-in-progress ${{ github.event_name == 'pull_request' }}. A step with if github.event_name == 'workflow_dispatch' runs rhiza check --base HEAD~1. go test -trimpath ./cmd/rhiza -run TestScripts/ci-install$ passed: the installed file cmp-matches the expected workflow. go test -trimpath ./internal/api with RHIZA_ROOT/CHANGE/RUN unset passed (TestCIInstall* included). A first internal/api run with the run's RHIZA_* env inherited had failed, which comes from that environment and not from this change.

## Links

### Out

- verifies [ci-install-dispatchable](../capability/rhiza/integrity/ci-install-dispatchable.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
