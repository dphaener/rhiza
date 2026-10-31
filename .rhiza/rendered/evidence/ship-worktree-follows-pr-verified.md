# ship-worktree-follows-pr-verified

evidence · `01M68E2VHYW0A197EDVTH675C9` · actor [claude-code-691e9e00](../actor/claude-code-691e9e00.md)

## result

PASS

## statement

At 0c3e831, go test ./cmd/rhiza -run TestScripts/ship-follows-pr-head passes. The script uses a bare origin. After the branch is pushed with an upstream, a second clone pushes a commit to cs/ship-it. With a dirty feature.txt, rhiza ship rebase fails with 'uncommitted changes' and names feature.txt, and HEAD is unchanged. Clean, ship rebase fast-forwards HEAD to the moved head. After another outside push, ship land follows it and its push succeeds: origin's cs/ship-it equals the worktree HEAD, and HEAD contains the pushed head. A local-only commit with origin ahead is refused, naming 'Local only commit', and HEAD stays put. The script pushes the branch with git push -u instead of running ship open-pr, but that gives open-pr's resulting state. In a scratch worktree with the three shipFollowPR calls removed, the script fails twice. With the dirty step, rebase unexpectedly succeeds. With the dirty step removed, the clean rebase leaves HEAD at the old commit instead of the pushed head. So the script really guards the fix. go test ./internal/api -run TestShipFollow passes: TestShipFollowFastForwards, TestShipFollowRefusesDirtyWorktree and TestShipFollowRefusesLocalCommits. With a fake git, the last two return a conflict that names feature.txt or the local commit, and they make no merge, reset or push call. The diff calls shipFollowPR in ship rebase, land and open-pr before any commit or push. It returns early when origin has no such branch or already holds HEAD. fmt-check, vet and go test ./... pass. A run of make test right after make check failed only the timing-budget read test, and a rerun passed. make self-check fails only on this verifier session's own rendered actor file, claude-code-691e9e00, which the run recorded after the leaf commit.

## Links

### Out

- verifies [ship-worktree-follows-pr](../capability/rhiza/orchestration/shipping/ship-worktree-follows-pr.md)

## History

- `01M684XYB96R52XXCKMV80X5T1` When a run's PR branch is updated by someone other than the run, such as the orchestrator fixing a conflict, the run's wait-pr, rebase and rebase-merge steps follow the new head instead of stalling until wait-pr's timeout or failing
