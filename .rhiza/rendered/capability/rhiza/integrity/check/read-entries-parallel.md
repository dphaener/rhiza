# read-entries-parallel

capability · `01M7GFXTNFE5TRP176PWG9TATG` · parent [check](../check.md)

## acceptance

go test -trimpath -p 2 -count=1 ./internal/store passes. A new TestReadEntriesParallelSame, over a log of a few hundred entries with an unreadable file, an entry that does not decode and an entry whose ID does not match its path, shows the parallel result deep-equals the serial one. go test -race -trimpath -p 2 -count=1 ./internal/store -run ReadEntries is clean. Failure case: with GOMAXPROCS=1 it decodes on one worker and still returns the same Log. Measure, with @check-whatif-once already in: build $TMPDIR/rhiza-x as in @check-whatif-once and record in the commit's Checks: block /usr/bin/time -f '%e s' $TMPDIR/rhiza-x check in this repo's main checkout, after one warming run, before and after this leaf, and store.ReadLog's cumulative pprof share. After this leaf, check takes at most 5 s on the quiet 4-core VM, milestone 39's target (@projection-targets). Touches internal/store/read.go and its tests. Must not add a dependency or change the order or content of Entries or Problems.

## statement

store.ReadEntries (internal/store/read.go:97) decodes entry files on up to four workers (bounded by GOMAXPROCS), the way the index's listing already stats folders in parallel, and returns the same Log as the serial decode: the same entries and problems in the same order

## Links

### Out

- depends-on [check-whatif-once](check-whatif-once.md)
- decided-by [projection-scope-approved](../../../../decision/projection-scope-approved.md)
- decided-by [projection-targets](../../../../decision/projection-targets.md)
- decided-by [check-time-target](../../../../decision/check-time-target.md)

### In

- verifies [read-entries-parallel-verified](../../../../evidence/read-entries-parallel-verified.md)
- verifies [read-entries-parallel-qa-1](../../../../evidence/read-entries-parallel-qa-1.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
