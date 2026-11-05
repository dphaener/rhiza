# ship-main-dispatch-ci

decision · `01M6MABYJY48N0S4D6DM9GJE8Q`

## reasoning

Read the change set's intent, ship-main-tree-check, shipCheckMainTree and shipRerunFailed in internal/api/ship.go, the wait-main/main-green steps in ship-change.toml, .github/workflows/ci.yml and the ci install template in internal/api/ci.go. (1) The parent #256 served milestone-25 and this is a fix that lets daemon-built change sets ship, which is exactly that milestone. (2) workflow_dispatch takes only a branch or tag, so --ref <sha> fails; dispatching on main alone would test main's tip, which may already hold later merges and so would not prove the merge commit green. A sha input plus checkout ref pins the tested commit, and a run-name carrying the sha gives ship a reliable way to find the run it dispatched without racing other dispatches. (3) ship-change is generic across repos, so hard-coding ci.yml is wrong; the workflows that matter are the ones that gate PRs, and only ones declaring workflow_dispatch can be dispatched. Stopping loudly when none exist keeps the no-untested-main guarantee instead of silently passing. (4) Matching the PR checks' existing rerun-once policy absorbs flaky runners consistently; beyond that a red main needs a human, as the old main-green did on push CI failure (it had no on.fail), and an automatic revert would be a new, unasked-for behaviour.

## statement

(1) Change set 01M6MA6P serves milestone-25 (empty the backlog), like its parent #256: it unblocks the daemon's own shipping. (2) The dispatched run tests the merge commit itself, pinned by input, not main's tip: ship runs `gh workflow run <file> --ref main -f sha=<merge commit>`; the workflow's workflow_dispatch trigger declares a required `sha` input, checks out `ref: ${{ inputs.sha || github.sha }}`, and sets `run-name` to include the sha, so ship finds its run by the sha in the run's display title (event workflow_dispatch) and waits on it, never on whatever main has advanced to. (3) Ship dispatches every workflow file under .github/workflows that declares both a pull_request and a workflow_dispatch trigger (the workflows that gate PRs and can be dispatched); in this repo that is ci.yml, and the `rhiza ci install` template gains the same workflow_dispatch trigger with the sha input. When the trees differ and no such workflow exists, the main step stops with a stop: line naming the merge commit and saying no dispatchable workflow was found. Trees equal: nothing is dispatched, as now. (4) A dispatched run that fails gets the same one rerun of its failed jobs that PR checks get (shipRerunFailed, once per merge commit); if it still fails, main-green fails with a stop: line naming the merge commit and the failed checks, and the run stops for a person: no revert, no automatic fix, no close-out.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-ec5a304e-1](../question/question-ec5a304e-1.md)
- derived-from [plan-ship-main-dispatch](../source/plan-ship-main-dispatch.md)

### In

- decided-by [ship-main-tree-check](../capability/rhiza/workflows/workflow-runs/ship-main-tree-check.md)
- decided-by [ship-main-dispatches-ci](../capability/rhiza/workflows/workflow-runs/ship-main-dispatches-ci.md)
- decided-by [ci-yml-dispatchable](../capability/rhiza/integrity/ci-yml-dispatchable.md)
- decided-by [ci-install-dispatchable](../capability/rhiza/integrity/ci-install-dispatchable.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
