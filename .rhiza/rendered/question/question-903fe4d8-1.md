# question-903fe4d8-1

question · `01M7GZ5BM1JN6WYTY7VCC2CJFZ`

## session

903fe4d8-83ee-438a-b6fd-01a1e1d9fa1c

## statement

Change set 'check without what-if rebuilds': all three leaves are built. The read-entries-parallel verifier fails that leaf only on the 5 s check target (check-time-target). At e4d148ed5, rhiza check in ~/code/rhiza took 11.67 s wall at load ~3 (16.25 s user + 4.14 s sys, which is ≥5.1 s even on 4 perfect cores). The builder measured 7.3–9.5 s at load ~2, down from 20.6 s at base. The decode it parallelised was only ~1.3 s of CPU. The rest is WhatIfFrom (~19%), checker.model/inFlightLinks and render, all outside this leaf. Everything else (tests, -race, same output order, no dependency) passes. check-time-target says the leaf fails and the gap is raised, not the target loosened. How should it be settled?

## Links

### In

- answers [check-time-gap-carried](../decision/check-time-gap-carried.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
