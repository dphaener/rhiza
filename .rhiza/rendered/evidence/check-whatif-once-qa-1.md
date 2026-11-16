# check-whatif-once-qa-1

evidence · `01M7H08YP6DS068H4FX97XCH0D` · actor [claude-code-b257b28a](../actor/claude-code-b257b28a.md)

## result

PASS

## statement

QA built 3bcbde2fe to $TMPDIR/rhiza-x and main 5180b25eb to $TMPDIR/rhiza-main. On a scratch repo with a landed base, in-flight A (forged verifies link, forged dangling @@{01ZZ...} text ref, ref to a since-retired node), B depends-on A and edits A's node, and C whose set on store is stale after a record retire, rhiza check reported the bad in-flight link, the broken text ref, the retired-ref warning and C's what-if conflict once each, exit 1; human and --json output were byte-identical to main, also under GOMAXPROCS=1. rhiza check --json on ~/code/rhiza (5180b25eb), on this worktree and on a fresh init repo (unrendered and rendered) were byte-identical to main.

## Links

### Out

- verifies [check-whatif-once](../capability/rhiza/integrity/check/check-whatif-once.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
