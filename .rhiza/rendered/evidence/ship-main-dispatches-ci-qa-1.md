# ship-main-dispatches-ci-qa-1

evidence · `01M6MDYMTMHQP1XS1WXV88EBJH` · actor [claude-code-9aaf5286](../actor/claude-code-9aaf5286.md)

## result

PASS

## statement

QA at 5188b19f with a scratch build (go build ./cmd/rhiza), a scratch repo, a local bare origin and a fake gh, run by hand. Tested head on cs/ship-it; main advanced by another commit, then a squash commit, so the trees differ; state file in the pre-#256 shape (pr, url, branch, head, checks, merge_commit). Origin/main held this repo's real .github/workflows/ci.yml, the rhiza.yml that rhiza ci install writes, and a push-only push.yml. ship check --main --settled: poll 1 'not settled: no run of ci.yml, rhiza.yml on 54c422b yet' (exit 1); poll 2 in_progress (exit 1); poll 3 '4 checks passed' (exit 0); then ship check --main wrote main checks (written true). gh.log: exactly one 'workflow run ci.yml --ref main -f sha=<merge commit>' and one for rhiza.yml across 4 polls, push.yml never; a run titled with another sha was ignored. With the run failed: one 'run rerun 42 --failed', then --settled exits 0 with 'failed: CI', and ship check --main exits 75 with 'stop: merge commit 54c422b: CI failed on main'. When main's ci.yml lost its sha input: exit 75 'workflow ci.yml lacks the sha input', no workflow run. With no dispatchable workflow: exit 75 'no dispatchable workflow was found at origin/main', no workflow run. Trees equal: 'the merge commit's tree is the tested head's', exit 0, no gh call. Real gh against GitHub accepted the same 'run list --workflow ci.yml --event workflow_dispatch --json databaseId,name,displayTitle,status,conclusion' call. The ci-install and four run-ship-main-dispatch/ship-main-dispatch-select testscripts passed.

## Links

### Out

- verifies [ship-main-dispatches-ci](../capability/rhiza/workflows/workflow-runs/ship-main-dispatches-ci.md)

## History

- `01M6MA6P7KKMQPFBM6ZD6EMDM3` Fix Ship's new main-tree check (ship-main-tree-check, #256) stops every run whose PR merged after another PR: GitHub's squash merge onto an advanced main gives a merge commit whose tree differs from the tested head's tree (run 01M6KXGZYGRY stopped at wait-main: 'merge commit b9229f0 has tree 6b107af, but the tested head ab3969c has tree 0e2c59e'), and with CI no longer running on push to main nothing ever turns it green. Keep the no-retest rule for the common case (trees equal: no CI on main), and when they differ dispatch CI on that main commit (add workflow_dispatch to .github/workflows/ci.yml and the ci install template; gh workflow run ci.yml --ref <sha>) and wait for it, so the owner's ask holds: "Update github runner to: Cancel stale runs, and stop re-testing on main". Runs started before #256 (old wait-main on push CI) must also finish: on resume they should take the new path
