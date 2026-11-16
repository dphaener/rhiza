# inbox-rescoped-after-parity

decision · `01M7H101D82M3QXB920R8D196P`

## question

The inbox no longer folds the log with an index and takes 0.62 s, mostly parsing run journals: does the inbox leaf keep the 200 ms target?

## reasoning

The milestone plan's 2.3 s came before #599: with an index the inbox already reads states and in-flight change sets from it, and inboxGateRead already opens the worktree's own index, so only the no-index fold and the hand-opened indexes are left to move. A probe on the main checkout showed the rest of the time in run journal parsing, not in any log fold; making run.Store.List cheaper is work on the run store, not on reads from the projection, so holding this leaf to 200 ms would make it unfinishable inside its intent.

## statement

The inbox leaf only moves inboxFold (internal/api/inbox.go:602-619 and the overlaps at :685) and inboxGateRead (:1405) onto the projection interface; it carries no timing target. The inbox's remaining time (0.62 s on 2026-11-16 at 6ebf119de, 40% of its CPU in run.Store.List parsing run journal files) is outside this change set's intent and goes to a discovered-from draft change set, which serves m39-reads-fast for the inbox

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [reads-projection-plan-2026-11-16](../source/reads-projection-plan-2026-11-16.md)

### In

- decided-by [inbox-from-projection](../capability/rhiza/views/fast-reads/inbox-from-projection.md)

## History

- `01M7GFTDY50TAHVDP688ENP6TX` Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.
