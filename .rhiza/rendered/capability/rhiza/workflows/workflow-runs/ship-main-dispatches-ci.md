# ship-main-dispatches-ci

capability · `01M6MAKAEANQJBCM1WFMTKWJ7P` · parent [workflow-runs](../workflow-runs.md)

## acceptance

Testscripts with a local bare origin and a fake gh, replacing run-ship-merge-tree-differs: (a) differing trees: wait-main calls gh workflow run <file> --ref main -f sha=<merge commit> exactly once across its polls, for a workflow on origin/main's tip declaring pull_request and workflow_dispatch with a sha input and a run-name referencing inputs.sha, and not for one declaring only push, keeps waiting while the fake run whose title holds that sha is in progress, and the run reaches close-out once it succeeds; (b) the fake dispatched run fails, is rerun once and fails again: main-green fails with a stop: line naming the merge commit and the failed check, and the run never reaches close-out; (c) differing trees and no dispatchable workflow: a stop: line naming the merge commit and saying no dispatchable workflow was found, with no gh workflow run call; (d) a run failed at wait-main by the old tree-differs stop, with a hand-written pre-#256 state file (pr, url, branch, head, checks, merge_commit only), dispatches and reaches close-out after rhiza run resume; (e) differing trees and a workflow declaring pull_request and workflow_dispatch but no sha input (and, separately, one whose run-name lacks inputs.sha): a stop: line naming the merge commit and that workflow file, with no gh workflow run call; (f) main advanced past the merge commit with a commit that adds the sha input to a workflow lacking it at the merge commit: the workflow is selected from origin/main's tip and dispatched, so selection follows the dispatched ref, not the merge commit.

## statement

When the merge commit's tree differs from the tested head's (main moved between the checks and the merge, or the PR was merged by a ship from before #256), ship-change's wait-main step runs CI on the merge commit instead of stopping: it selects the workflows under .github/workflows at origin/main's tip (the ref it dispatches on, whose definitions GitHub runs) that declare both a pull_request and a workflow_dispatch trigger, and dispatches each of them, once and recorded in the ship state file (`gh workflow run <file> --ref main -f sha=<merge commit>`), finds each dispatched run by the merge commit's sha in its display title (event workflow_dispatch) and waits for it to finish; a failed run's failed jobs are rerun once per merge commit as PR checks are, and a second failure fails main-green with a stop: line naming the merge commit and the failed checks, with no revert and no close-out. A selected workflow whose workflow_dispatch declares no `sha` input, or whose run-name does not reference inputs.sha, cannot pin or be found by the merge commit: the main step stops, before dispatching anything, with a stop: line naming the merge commit and that workflow file and saying it lacks the sha input or run-name. Differing trees with no selected workflow stop the main step with a stop: line naming the merge commit and saying no dispatchable workflow was found. A run stopped at, or still waiting in, wait-main whose ship state file was written before this change (pr, head, checks, merge_commit, no dispatch record) takes the same path when resumed.

## Links

### Out

- decided-by [ship-main-dispatch-ci](../../../../decision/ship-main-dispatch-ci.md)
- derived-from [plan-ship-main-dispatch](../../../../source/plan-ship-main-dispatch.md)

### In

- verifies [ship-main-dispatches-ci-verified](../../../../evidence/ship-main-dispatches-ci-verified.md)
- verifies [ship-main-dispatches-ci-qa-1](../../../../evidence/ship-main-dispatches-ci-qa-1.md)
- verifies [pr-265-fix-ship-s-new-main-tree-check](../../../../evidence/pr-265-fix-ship-s-new-main-tree-check.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
