# question-ec5a304e-1

question · `01M6MA9ZQZFFFQAY6ZSY834Y62`

## session

ec5a304e-29cc-4f00-9ede-01a1a8a374d3

## statement

Change set 01M6MA6P (Fix Ship's new main-tree check) serves no intent yet. Which intent should it serve? Its parent change set #256 (Update github runner) served milestone-25 (empty the backlog). / When the merge commit's tree differs from the tested head's, ship dispatches CI on main. GitHub's workflow_dispatch only accepts a branch or tag as ref, not a SHA, so `gh workflow run --ref <sha>` cannot pin the merge commit. Which commit should the dispatched run test? / Which workflows should ship dispatch when the trees differ? ship-change is generic across repos (this repo has ci.yml; `rhiza ci install` writes its own file). / If the dispatched CI run on main fails, what should wait-main/main-green do?

## Links

### In

- answers [ship-main-dispatch-ci](../decision/ship-main-dispatch-ci.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
