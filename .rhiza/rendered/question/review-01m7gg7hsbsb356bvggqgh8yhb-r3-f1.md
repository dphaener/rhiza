# review-01m7gg7hsbsb356bvggqgh8yhb-r3-f1

question · `01M7GHJ4GE0JYYYEPT5AH36X8G`

## statement

Plan review round 3 left finding f1 (oversized-leaf on check-whatif-once): The leaf holds two concerns: a model-engine change (new model.AsIsSeed and model.ToBeFrom in internal/model/seed.go, and Seed.builder deep-copying each node's and link's History, Fields and aspects, which also changes the path internal/index/whatif.go:289 takes through model.WhatIfFrom) with its own TestAsIsSeedWhatIfFrom equivalence test; and rewiring three check passes (checker.model at internal/check/check.go:378/:384, checker.inFlightLinks at :456, checker.textRefs at internal/check/textrefs.go:49) onto shared what-if views with TestWhatIfOnceMatches, plus two profiled timings. The model half stands on its own (its test needs no check change) and the check half needs only its API, so they split into a model leaf and a check leaf that depends on it. The critic confirmed it: Verified against the slice and code: after two revisions check-whatif-once now carries a model-engine change (new AsIsSeed and ToBeFrom in internal/model/seed.go, plus a deep copy in Seed.builder that also changes the path internal/index/whatif.go:289 takes through WhatIfFrom) with its own model test, and a separate rewiring of three check passes (check.go:378/:384/:456, textrefs.go:49) with its own check test and two profiled timings. The model half is testable without any check change and the check half needs only its API, so the leaf holds two concerns and splits cleanly into a model leaf and a dependent check leaf. oversized-leaf is major per the plan-review severities; major stands. How should it be settled?

## Links

### Out

- blocks [check-whatif-once](../capability/rhiza/integrity/check/check-whatif-once.md)

### In

- answers [split-check-whatif-once](../decision/split-check-whatif-once.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
