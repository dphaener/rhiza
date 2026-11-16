# change-reads-from-projection-verified

evidence · `01M7H727J26S4D6R4PTH581NVG` · actor [claude-code-cc41b3b3](../actor/claude-code-cc41b3b3.md)

## result

PASS

## statement

Verified commit 1b563a4a9. Re-ran every Checks: command under the capped scope, all exit 0: go test -trimpath -count=1 -p 2 ./internal/api -run 'ChangeLs' (ok); -run 'TestPropertyProjectionParity|TestChange' (ok, parity green); ./cmd/rhiza -run 'TestScripts/change' with RHIZA_/SHIP_/REX_ unset (11 change*.txtar, ok); go vet ./internal/api; gofmt -l internal/api (empty); go build; rhiza ship generate; rhiza render (no new tracked changes). Diff: only internal/api/change.go (runChangeLs over w.projection(): stored folds plus storedLogFoldProblems, log then fold, same format strings as logProblems/foldProblems; FoldLog path kept when IndexPath is empty) and change_test.go (TestChangeLsIndexedReadsNoLog: zero log reads over the index, same sets/problems/text as over the log with a damaged entry, vanished folder drops the change set). change show and change_indexed.go untouched; no TODOs. Byte comparison against a scratch build of the parent 1b563a4a9^: change ls and change ls --json in /home/rhiza/code/rhiza identical (6863 lines). In a scratch repo with a junk.json log problem and an out-of-order open entry (fold problem), text and --json were identical to the parent's, order and wording included; after rm -rf of a change set folder, the warm index listed the remaining sets with no stale one, identical to the parent. Note: rewriting an already-indexed entry file in place showed its old problem until the index was rebuilt. That is the index's documented rule (index.go:15, an entry file never changes) and not part of this leaf. Warm change ls 0.42 s.

## Links

### Out

- verifies [change-reads-from-projection](../capability/rhiza/views/fast-reads/change-reads-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
