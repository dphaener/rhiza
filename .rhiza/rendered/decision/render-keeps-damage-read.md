# render-keeps-damage-read

decision · `01M7H101JW71NKB5A14GCDDCWV`

## question

Does render's damaged-log check also move onto the projection's stored log problems?

## reasoning

The index does not notice an edit inside an entry file it has already listed (internal/index/index.go doc), and render's refusal of a damaged log is the safety net for exactly that case, so the damage check must keep looking at the files; the fold and the as-is rebuild are the whole-log work the intent removes.

## statement

rhiza render keeps damagedLog's own read of the log (internal/api/render.go:97-100) and takes only the change sets and the as-is model from the projection

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [reads-projection-plan-2026-11-16](../source/reads-projection-plan-2026-11-16.md)

### In

- decided-by [render-from-projection](../capability/rhiza/views/fast-reads/render-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
