# render-from-projection-qa-1

evidence · `01M7H9P8400K9J09VB9P87ZZC6` · actor [claude-code-2c502734](../actor/claude-code-2c502734.md)

## result

PASS

## statement

QA: scratch builds of HEAD 9f005fded and of base 03b972942 (go build -trimpath, capped scope). Scratch repo: render byte-identical to base 03b972942's and to the no-index path (diff -r), incl. evidence/logout-ev.md and 'constrains fast (retired)'; with an undecodable entry both refuse with the same 'damaged log, nothing rendered' error, exit 1; after a corrupt index, render reported 5 unchanged.

## Links

### Out

- verifies [render-from-projection](../capability/rhiza/views/fast-reads/render-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
