# common-dir-once-verified

evidence · `01M7G01EZ860CWN1F2SJN16MSR` · actor [claude-code-25875785](../actor/claude-code-25875785.md)

## result

PASS

## statement

Verifier at leaf commit 26a9b80ce (all under flock+systemd-run, RHIZA/SHIP/REX unset): go test -trimpath -p 2 -count=1 ./internal/api -run '^TestGitCommonDirOnce$' PASS (subtests repo and repo.wt: one git call rev-parse --path-format=absolute --git-common-dir, same absolute cleaned path, orphanLogEntries makes no further call); -run '^TestOrphanLogSymlinkRoot$' PASS (real repo via symlinked root, common dir equals --absolute-git-dir, orphan entry listed); the combined TestGitCommonDirOnce|TestOrphanLogSymlinkRoot|TestShareCommonDir|PrePush|Session|Orphan|Untracked run, -run 'Session|Orphan|Untracked', -run 'Daemon|PrePush|Ship|Check' and ./internal/web all exit 0; cmd/rhiza TestScripts session-run, session-start, apply, run-amend, check-orphan-log, model-commit-untracked each PASS (IDs unchanged, scripts match), and describe-golden, plan-workflow-pr, ship-prepush-check, ship-land-branch-build, run-ship-change-rebase, run-status PASS; go vet ./internal/api ./cmd/rhiza exit 0; rhiza ship generate exit 0 with 0 changed files; rhiza render exit 0. Measure: strace -f -s 300 -e trace=execve,clone,clone3,fork,vfork on rhiza.test -test.run 'TestScripts/^run-ship-conflict-fix-more$' PASS, 313 rev-parse --git-common-dir across about 300 rhiza children (was 757 for 287). Every child asks at most once except 4 'run start --foreground' processes: 1 call in repo, then 3 in the build worktree repo.ship-it from run_leaves.go openCheckout through worktreeGit (unchanged code), which decision common-dir-asked-once keeps as it is. Diff touches only the named files plus ship_prepush.go (shareCommonDir) and fake gits in the daemon tests, both explained in the commit. No TODOs or stubs. Before this leaf (HEAD~1), a clean-checkout render already left .rhiza/rendered/index.md with stale actor-index lines; this leaf's own actor line is present.

## Links

### Out

- verifies [common-dir-once](../capability/common-dir-once.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
