# index-as-is-read

decision · `01M7H101G3THQECHBD67A6C36M`

## question

Where does render's as-is model.Model come from once it stops calling model.AsIs over the folded log?

## reasoning

render.Render takes a whole internal/model Model, and the index only answers row reads, overlays and what-ifs of an in-flight change set; rebuilding as-is in internal/api from rows would be another copy of seed rules. The strict engine stays the spec because the rows are what it computed. Splitting it from render keeps each leaf one concern.

## statement

The index gets one read, (*index.Index).AsIs() (model.Model, []model.Problem, error), that assembles the stored as-is rows (nodes with fields and histories, links, the as-is problems) into the model.Model render.Render takes; it adds no rule and changes no stored row, and is its own leaf (index-as-is-model) that render-from-projection depends on

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [reads-projection-plan-2026-11-16](../source/reads-projection-plan-2026-11-16.md)

### In

- decided-by [index-as-is-model](../capability/rhiza/views/fast-reads/index-as-is-model.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
