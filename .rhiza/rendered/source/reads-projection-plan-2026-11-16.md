# reads-projection-plan-2026-11-16

source · `01M7H0XYE27BMPZ66KZRVX6K8H`

## quotes

```json
[]
```

## raw_sha256

87a2065fbf7dd62cd9dd8e255aa8b60e342b0268964dc9e2b4de4c81e985835d

## session

da57a0d6-2dc0-4cb8-95be-01a1e20b013f

## summary

Autonomous plan run for change set 01M7GFTDY50TAHVDP688ENP6TX, re-planning the leaves the milestone plan (#596) wrote, after the parity harness (#599) and check without rebuilds (#600) landed. The orchestrator (arbiter) asked that design choices be recorded as assumed rulings, that each leaf stay one concern, that the planner stay out of 01M7GFTS0N's leaves and keep the interim model's deletion in 01M7GFTKCA. Research found: newModelSets and the daemon fold already build as-is once (newModelBase) since #599; inboxGateRead already reads the worktree's own index; inbox with an index reads no log, and one probe on the main checkout at 6ebf119de measured rhiza inbox at 0.62 s (40% of its CPU in run.Store.List parsing run journals) and rhiza status at 1.93 s; writeIndex in write_indexed.go is already an index with the dry run's pending entries laid over it; the index has no read that returns the whole as-is model.Model that render.Render takes. Leaves were rescoped to that: the projection interface grows from writeIndex and is opened lazily, the log-once memo is its own leaf, the inbox leaf only moves its reads onto the interface, the inbox's remaining time goes to a discovered-from change set, render gets an index as-is read as its own leaf and keeps its damage read of the log.

## title

Planning Reads from the projection, autonomous (2026-11-16)

## Links

### In

- derived-from [projection-grows-from-write-index](../decision/projection-grows-from-write-index.md)
- derived-from [log-once-is-own-leaf](../decision/log-once-is-own-leaf.md)
- derived-from [inbox-rescoped-after-parity](../decision/inbox-rescoped-after-parity.md)
- derived-from [index-as-is-read](../decision/index-as-is-read.md)
- derived-from [render-keeps-damage-read](../decision/render-keeps-damage-read.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
