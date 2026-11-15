# m39-reads-fast

intent · `01M7GFT02J9RW50MCR273H3X2P` · parent [milestone-39](../milestone-39.md)

## statement

rhiza status and rhiza inbox each answer within 200 ms warm on this repo's own log (2026-11-15 at 24,560 entries: 2.0 s and 2.3 s). Measure, in the main checkout with nothing else running: run each once, then /usr/bin/time -f '%e s' rhiza status and /usr/bin/time -f '%e s' rhiza inbox.

## Links

### Out

- decided-by [projection-targets](../../decision/projection-targets.md)

### In

- serves change set `01M7GFTDY50TAHVDP688ENP6TX` (Every read goes through one projection interface behind workspace, whose only implementation is the index with dry-run pending entries laid over it: status and its hook, the inbox, render, change ls and the daemon fold stop folding the whole log, and workspace reads the log at most once per call. Evidence: status.go:123, render.go:99, change.go:669, inbox.go:619 and daemon_fold.go:197 call changeset.FoldLog; inboxGateRead clears IndexPath; workspace.log() rereads the disk on every call; status 2.0 s and inbox 2.3 s on this repo.)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
