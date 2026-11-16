# read-entries-parallel-verified

evidence · `01M7GZAZ456JY6VSQQZRGRX724` · actor [claude-code-4479f0f7](../actor/claude-code-4479f0f7.md)

## result

PASS

## statement

Second verifier, at builder commit e4d148ed5 (internal/store and go.mod unchanged through HEAD e7a5084d1). Ran under the capped scope: go test -trimpath -p 2 -count=1 ./internal/store: exit 0. go test -race -trimpath -p 2 -count=1 ./internal/store -run ReadEntries: exit 0, race-clean; TestReadEntriesParallelSame (300 entries over 7 change sets plus an unreadable file, an undecodable entry and an ID/path mismatch; 2, 4 and len(paths) workers and reversed paths deep-equal the serial Log) and TestReadEntriesOneProc (GOMAXPROCS=1 gives decodeWorkers 1 and the same Log) pass. Diff read: each path's result goes to its own slot, gathered in path order before the unchanged sort, so Entries and Problems keep their order and content; workers = max(1, min(4, GOMAXPROCS, n)); imports only runtime, sync and sync/atomic, no go.mod or go.sum change; only read.go, read_test.go and the builder's actor record are touched. The 5 s clause is UNMET, not passed: the first verifier measured rhiza check at 11.67 s wall (16.25 s user + 4.14 s sys, load ~3), and the builder measured 7.34-9.53 s (load ~2), against a 5 s target. Per Decision check-time-gap-carried, the leaf is settled on the rest of its acceptance and the gap is carried to draft change set 01M7GZ8K2B; check was not re-timed here.

## Links

### Out

- decided-by [check-time-gap-carried](../decision/check-time-gap-carried.md)
- verifies [read-entries-parallel](../capability/rhiza/integrity/check/read-entries-parallel.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
