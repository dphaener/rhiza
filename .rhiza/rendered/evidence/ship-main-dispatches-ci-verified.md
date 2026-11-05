# ship-main-dispatches-ci-verified

evidence · `01M6MDMV4JBPWPJBCGG9X7KYCJ` · actor [claude-code-d0a17d19](../actor/claude-code-d0a17d19.md)

## result

PASS

## statement

Verified at fc79973e by claude-code-d0a17d19 (verifier). run-ship-merge-tree-differs.txtar is gone (renamed to run-ship-main-dispatch-fail). go test -trimpath ./cmd/rhiza -run 'TestScripts/(run-ship-main-dispatch|run-ship-main-dispatch-fail|run-ship-main-dispatch-resume|ship-main-dispatch-select)$' -count=1 -v: all four PASS, plus describe-golden and every run-ship* script; go test ./internal/api/... and ./internal/workflow/... ok. (a) run-ship-main-dispatch: local bare origin and fake gh; exactly one 'workflow run ci.yml --ref main -f sha=<merge commit>' across >=3 run-list polls with the run in_progress for the first two, push.yml (push only) never dispatched, no rerun, run reaches close-out. (b) run-ship-main-dispatch-fail: one 'run rerun 42 --failed', then 'stop: merge commit <sha>: CI failed on main: CI', no close-out. (c) ship-main-dispatch-select and run-ship-main-dispatch-resume: 'no dispatchable workflow was found' stop with no workflow run call. (d) resume script: hand-written pre-#256 state file (pr,url,branch,head,checks,merge_commit), rhiza run resume --retry dispatches once and reaches close-out. (e) select script: stops naming ci.yml for a missing sha input and, separately, for a run-name without inputs.sha, no workflow run call. (f) select script covers a tip fixing the run-name; I also reproduced (f) exactly as written by hand with a scratch build and fake gh: merge commit whose ci.yml lacks the sha input stops (exit 75), then after main gains a commit adding the sha input, ship check --main --settled dispatches 'workflow run ci.yml --ref main -f sha=<merge commit>' and passes (exit 0). Read the diff: ship check --main fetches origin first, selection reads refs/remotes/origin/main, dispatch is recorded in the state file before gh runs, reruns are keyed on the merge commit.

## Links

### Out

- verifies [ship-main-dispatches-ci](../capability/rhiza/workflows/workflow-runs/ship-main-dispatches-ci.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
