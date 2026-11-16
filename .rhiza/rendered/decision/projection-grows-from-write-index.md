# projection-grows-from-write-index

decision · `01M7H1017W72NYN680ZGD9P2H8`

## question

Does the projection interface wrap the existing writeIndex or start as a new type, and what does a read do when ctx.IndexPath is empty?

## reasoning

writeIndex already answers change sets, a change set's entries and the strict model of a view over the index with pending entries overlaid, which is the interface the intent names; a second implementation beside it would be a fourth copy the next change set must delete. Opening it lazily keeps ops that never read the model from paying an index open, and leaving the empty-IndexPath branch in place keeps this change set from doing 01M7GFTKCA's deletion, which the orchestrator kept out of it.

## statement

The projection interface's one implementation grows out of writeIndex (internal/api/write_indexed.go: the index with the dry run's pending entries laid over it, already used by every write, run_model_cache.go and inboxGateRead), moved behind an interface in a new internal/api/projection.go, not a new type beside it; read.go's reader stays as it is for 01M7GFTKCA's writes-from-projection. workspace opens it lazily, at most once per call, when ctx.IndexPath is set and an op first asks for it; with ctx.IndexPath empty it is nil and each moved read keeps its log path, which 01M7GFTKCA's throwaway-index-fallback and interim-model-deleted remove

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [reads-projection-plan-2026-11-16](../source/reads-projection-plan-2026-11-16.md)

### In

- decided-by [projection-interface](../capability/rhiza/views/fast-reads/projection-interface.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
