# split-check-whatif-once

decision · `01M7GHM0P1KCRS7SMH71QMT3VE`

## reasoning

The finding was confirmed major by the critic in round 3 and checks out against the leaf: its statement and acceptance hold a model-engine change with its own equivalence test (TestAsIsSeedWhatIfFrom, needing no check change) and a separate rewiring of three check passes with its own test and profiles, which need only the model half's API. oversized-leaf is major, so leaving it unfixed is not an option, and the split is clean along that API seam: each half has its own test, the model half's deep copy (which also changes the index's WhatIfFrom path at internal/index/whatif.go:289) is verified on its own, and the check half's must-not-change-what-check-reports rule stays where check changes. Ordering model, then check, then read-entries-parallel keeps check-time-target's combined measurement and check-seed-in-memory's decision intact.

## statement

Split check-whatif-once into two leaves. A new model leaf (model-seed-views) carries model.AsIsSeed and model.ToBeFrom in internal/model/seed.go, Seed.builder's deep copy of each node's and link's History, Fields and aspects, and TestAsIsSeedWhatIfFrom with its acceptance as written today; it touches only internal/model. check-whatif-once keeps the check half: checker.model (check.go:378/:384), checker.inFlightLinks (:456) and checker.textRefs (textrefs.go:49) rewired onto one AsIsSeed and shared ToBeFrom/WhatIfFrom views, TestWhatIfOnceMatches, the failure case, the before/after profiles and timings (the before profile taken at its base, which the model leaf does not change for check), and depends-on model-seed-views. read-entries-parallel still depends on check-whatif-once, so the 5 s measurement still follows both

## wording

(assumed, would have escalated)

## Links

### Out

- answers [review-01m7gg7hsbsb356bvggqgh8yhb-r3-f1](../question/review-01m7gg7hsbsb356bvggqgh8yhb-r3-f1.md)

### In

- decided-by [model-seed-views](../capability/rhiza/integrity/check/model-seed-views.md)
- decided-by [check-whatif-once](../capability/rhiza/integrity/check/check-whatif-once.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
