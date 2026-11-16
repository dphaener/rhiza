# workspace-log-once-qa-1

evidence · `01M7H9P7XMVJ0KVJAKN3W5VXX7` · actor [claude-code-2c502734](../actor/claude-code-2c502734.md)

## result

PASS

## statement

QA: scratch builds of HEAD 9f005fded and of base 03b972942 (go build -trimpath, capped scope). Every read op (status, inbox, change ls, slice, change show, overlaps) with an empty IndexPath (no HOME/XDG_CACHE_HOME) gave the same bytes as with the index, in a clean log, with a damaged entry plus junk.json, and after a folder vanished; the chained dry run saw its earlier pending writes.

## Links

### Out

- verifies [workspace-log-once](../capability/rhiza/views/fast-reads/workspace-log-once.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
