# plan-ship-main-dispatch

source · `01M6MAEGNV4GZK5WJC2RP1VMRT`

## quotes

```json
[]
```

## raw_sha256

2c805e8fcd965bff3ecf157a59cde6a381a25be16e0310e0951b51db4581bd74

## session

ec5a304e-29cc-4f00-9ede-01a1a8a374d3

## summary

Autonomous plan run, no owner. The planner read ship-main-tree-check, shipCheckMainTree/runShipMerge/shipRerunFailed in internal/api/ship.go, ship-change.toml's wait-main and main-green steps, .github/workflows/ci.yml, the ci install template and the stopped run's error (merge commit b9229f0 tree 6b107af vs tested head ab3969c tree 0e2c59e). It found that GitHub's workflow_dispatch takes only a branch or tag ref, so `--ref <sha>` cannot pin the merge commit, that ship-change is generic across repos, and that ci.yml's concurrency group keys on the PR number (empty on a dispatch). It raised four choices to the arbiter: the intent served, which commit a dispatched run tests, which workflows to dispatch, and what a failed dispatched run does. The arbiter ruled (Decision ship-main-dispatch-ci): serve milestone-25; dispatch with --ref main and a required sha input the workflow checks out, its run-name carrying the sha so ship finds its run; dispatch every workflow declaring both pull_request and workflow_dispatch, stopping when none exists; rerun failed jobs once, then stop for a person. The planner then revised the three capabilities #256 built (ship-main-tree-check, ci-yml-cancels-stale, ci-install-cancels-stale) to carry the dispatch path, including runs resumed from before #256.

## title

Plan of Fix Ship's new main-tree check (autonomous, arbiter)

## Links

### In

- derived-from [ship-main-dispatch-ci](../decision/ship-main-dispatch-ci.md)
- derived-from [ship-main-dispatches-ci](../capability/rhiza/workflows/workflow-runs/ship-main-dispatches-ci.md)
- derived-from [ci-yml-dispatchable](../capability/rhiza/integrity/ci-yml-dispatchable.md)
- derived-from [ci-install-dispatchable](../capability/rhiza/integrity/ci-install-dispatchable.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
