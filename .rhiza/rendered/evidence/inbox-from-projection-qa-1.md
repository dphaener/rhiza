# inbox-from-projection-qa-1

evidence · `01M7H9P8A8D18666F236KSF45P` · actor [claude-code-2c502734](../actor/claude-code-2c502734.md)

## result

PASS

## statement

QA: scratch builds of HEAD 9f005fded and of base 03b972942 (go build -trimpath, capped scope). inbox and inbox --json byte-identical to base in /home/rhiza/code/rhiza and in a scratch root project with a failed run and overlapping in-flight change sets; with that run's change set folder deleted the item still lists, no error. Warm inbox in the main checkout 0.57 s (base 0.57 s; no target).

## Links

### Out

- verifies [inbox-from-projection](../capability/rhiza/views/fast-reads/inbox-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
