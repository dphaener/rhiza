# slice-and-show-from-projection-verified

evidence · `01M7H8DBRAKV1T8AJ5QJXR2NN9` · actor [claude-code-408563d1](../actor/claude-code-408563d1.md)

## result

PASS

## statement

Verified commit 9f005fded under the capped scope. go test -trimpath -count=1 -p 2 ./internal/api -run '^(TestSessionBrief|TestSlice|TestChangeShow|TestRunSpawn|TestRunInvoke|TestPropertyProjectionParity)' exit 0 (22 top-level PASS, incl. TestPropertyProjectionParity, new TestSliceProjectionOpensOnce and TestChangeShowProjectionOpensOnce, which count opens through ctx.openedProjection (fired in openWriteIndex before index.Open) and require exactly 1 per slice/change show with answers equal to the default-index or folded read). go test -trimpath -count=1 -p 2 ./cmd/rhiza -run 'TestScripts/(slice|change-show)' with RHIZA_/SHIP_/REX_ unset: exit 0, all 6 scripts (slice, slice-budget, slice-change, slice-errors, change-show-canonical-text, change-show-serves). go vet ./internal/api exit 0; grep -n 'index.Open' internal/api/change_indexed.go prints nothing; openIndex kept (slice.go:110) and called only in openSliceSource's ctx.IndexPath == "" branch; rhiza ship generate and rhiza render on a clean copy of HEAD leave no diff. Scratch builds of HEAD and parent d802ae4a6 run in /home/rhiza/code/rhiza give byte-identical output and exit codes for slice, slice --json, slice fast-reads, slice fast-reads --json --view as-is, slice --change <cs> fast-reads --json, slice nope (exit 4), change show [<cs>] [--json]; outside a repo slice is still not_found with next rhiza init (exit 4); in a fresh repo change show is still usage 'no change set: name one or make one current', next rhiza change ls (exit 2), text and --json identical. run_brief.go: needed by this leaf, since runSlice now honours ctx.IndexPath and the brief's worktree pass would otherwise read the caller's checkout's index against the worktree's root; it sets IndexPath to index.DefaultPath(worktree root), the same file openIndex opened before (and on error falls back to openIndex, which fails the same way), so a spawned child's brief slice is unchanged; TestSessionBrief/TestRunSpawn pass. projection.go/write_indexed.go add projection.store(), needed for the reads the projection did not cover. No TODOs or stubs.

## Links

### Out

- verifies [slice-and-show-from-projection](../capability/rhiza/views/fast-reads/slice-and-show-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
