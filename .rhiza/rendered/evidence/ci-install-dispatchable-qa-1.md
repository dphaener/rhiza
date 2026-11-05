# ci-install-dispatchable-qa-1

evidence · `01M6MDYSKAKMTRWRPAYJRGT840` · actor [claude-code-9aaf5286](../actor/claude-code-9aaf5286.md)

## result

PASS

## statement

QA at 5188b19f: a scratch build ran rhiza ci install in a fresh repo, and PyYAML parsed the .github/workflows/rhiza.yml it wrote. Triggers: pull_request, push (branches [main]), workflow_dispatch with inputs.sha required: true. run-name 'Rhiza ${{ inputs.sha || github.event.pull_request.title || github.sha }}'. Concurrency group rhiza-${{ github.event.pull_request.number || inputs.sha || github.sha }}, cancel-in-progress ${{ github.event_name == 'pull_request' }}. Checkout ref ${{ inputs.sha || github.sha }}, fetch-depth 0. The base-branch, before-push and first-push steps are unchanged, plus a step gated on workflow_dispatch running rhiza check --base HEAD~1. That command, run by the scratch build on this repo, exits 0. Ship selected this installed file as dispatchable and dispatched it on the merge commit by hand. go test -trimpath ./cmd/rhiza -run TestScripts/ci-install$ passed.

## Links

### Out

- verifies [ci-install-dispatchable](../capability/rhiza/integrity/ci-install-dispatchable.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
