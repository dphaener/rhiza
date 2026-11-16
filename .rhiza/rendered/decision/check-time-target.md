# check-time-target

decision · `01M7GGEBY7QNREZ2HWET34N8K7`

## reasoning

No other change set in milestone 39 takes on check's time, so this change set carries projection-targets' 5 s figure as it stands rather than a looser one of its own. The 2026-11-15 profile (13.8 s wall) puts model.WhatIf at 46% of samples and model.asIs at 26%, while the decode (store.ReadLog) is 7%: building the as-is once (with the to-be from the same seed) and each what-if once carries most of the gain, and the parallel decode adds the rest. Ordering the decode leaf second lets one measurement show the combined result; if check still misses 5 s, the leaf fails and the gap is raised, not the target loosened.

## statement

The time target for this change set is milestone 39's target for check (projection-targets): rhiza check on this repo's main checkout within 5 s, measured with /usr/bin/time after one warming run on the quiet 4-core VM, after both leaves. read-entries-parallel is built after check-whatif-once so its measurement covers both

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [check-whatif-profile-plan](../source/check-whatif-profile-plan.md)

### In

- decided-by [read-entries-parallel](../capability/rhiza/integrity/check/read-entries-parallel.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
