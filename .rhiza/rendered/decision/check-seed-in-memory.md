# check-seed-in-memory

decision · `01M7GGEBTPYMQHBY4882KF2V4X`

## reasoning

Only internal/index builds a model.Seed today (loadSeed, from its stored rows), and check reads the log files through store.ReadLog, runs on a fresh clone in CI and needs no daemon or index, so it cannot borrow the index's seed. AsIsState already returns the State half of a Seed, so an AsIsSeed beside it in seed.go is the small step. The strict engine (internal/model) is the spec, so WhatIfFrom over this seed must give exactly WhatIf's model and problems; a model test holds that.

## statement

rhiza check gets its as-is seed from a new model.AsIsSeed(sets) in internal/model/seed.go, returning the as-is Model, its problems and a model.Seed built in memory from the folded change sets, and never from internal/index

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [check-whatif-profile-plan](../source/check-whatif-profile-plan.md)

### In

- decided-by [check-whatif-once](../capability/rhiza/integrity/check/check-whatif-once.md)
- decided-by [model-seed-views](../capability/rhiza/integrity/check/model-seed-views.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
