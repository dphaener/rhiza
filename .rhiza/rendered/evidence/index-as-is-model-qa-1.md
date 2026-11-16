# index-as-is-model-qa-1

evidence · `01M7H9P80GSW3TAMJM3RHK1PEZ` · actor [claude-code-2c502734](../actor/claude-code-2c502734.md)

## result

PASS

## statement

QA: scratch builds of HEAD 9f005fded and of base 03b972942 (go build -trimpath, capped scope). rhiza render (as-is) over the index in a scratch repo with a landed evidence node, a retired constraint and a retired link endpoint wrote files byte-identical to base's render and to the no-index render; slice and status in the main checkout matched base.

## Links

### Out

- verifies [index-as-is-model](../capability/rhiza/views/fast-reads/index-as-is-model.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
