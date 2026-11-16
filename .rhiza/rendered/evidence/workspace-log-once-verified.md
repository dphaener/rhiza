# workspace-log-once-verified

evidence · `01M7H43S1V9SNKCGNMCG37KTVK` · actor [claude-code-50aedecb](../actor/claude-code-50aedecb.md)

## result

PASS

## statement

At 1232c1a00 re-ran under the capped scope: go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestWorkspaceLog|TestChange|TestWrite)' exit 0; -run TestPropertyProjectionParity and -run '^(TestDefaultView|TestWorkspaceLog|TestPropertyProjectionParity)' exit 0 (TestWorkspaceLogReadsDiskOnce, DryRunChains, WriteDropsRead, DamagedEntry, DefaultViewSeesPending, DefaultViewFoldsNothing PASS). Diff: workspace.go keeps the first diskLog read in w.disk, clones Entries/Problems per call and lays pending on the copy; writeEntry drops w.disk before st.Write on a real write, so dry runs chain and real writes are seen. The four new tests cover each acceptance point (one disk read over model/changeSets/entriesOfSet, dry-run chaining with one read and no files written, write then read gives 2 reads and the new entry, damaged entry same problems twice and as a fresh read). The read_test.go change (w.disk = nil in TestDefaultViewSeesPending) only resets the kept read before the test's disk-read counter: the dry write already read the disk, so without it the counter is 0 though the reader still goes through the log; its behavioural assertions (the rule sees the pending change set in draft) are unchanged, so it follows the test's own setup and hides no op output change. No other code path writes entries to the workspace's log root (model_branch writes a temp worktree, session writes the session store). Code changes are workspace.go and its tests plus that one test line; no TODOs. Full internal/api not rerun per brief.

## Links

### Out

- verifies [workspace-log-once](../capability/rhiza/views/fast-reads/workspace-log-once.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
