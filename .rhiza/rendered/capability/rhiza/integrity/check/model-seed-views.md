# model-seed-views

capability · `01M7GXQ4ZKWZZKMTP9JQPEPA3D` · parent [check](../check.md)

## acceptance

Tests: go test -trimpath -p 2 -count=1 ./internal/model passes. The existing TestAsIsSeedWhatIfFrom in internal/model/seed_test.go, over a log with a landed base, two in-flight change sets where one depends on the other, one in-flight change set whose ops no longer apply and a duplicate handle, also shows that ToBeFrom over AsIsSeed's seed equals ToBe (model and problems). ToBeFrom and every WhatIfFrom call run over one shared seed, the results are compared only after every call has run, and the seed is then compared with a fresh AsIsSeed: each must still be equal (no aliasing between views or with the seed). Failure case: a ToBeFrom whose result shared a node's History or Fields with the seed, or with a WhatIfFrom view, fails that comparison. Touches internal/model/seed.go and internal/model/seed_test.go only; adds no rule and changes no existing function's result.

## statement

model.ToBeFrom(seed Seed, sets []changeset.ChangeSet) (Model, []Problem) in internal/model/seed.go builds the to-be view from a seed (the accepted change sets layered on a copy of it), beside the model.AsIsSeed and the copying Seed.builder that the parity harness (#599, decision parity-seed-bridge-in-model) already added and this leaf reuses, adding no second constructor; TestAsIsSeedWhatIfFrom covers ToBeFrom too

## Links

### Out

- decided-by [split-check-whatif-once](../../../../decision/split-check-whatif-once.md)
- decided-by [check-seed-in-memory](../../../../decision/check-seed-in-memory.md)
- decided-by [parity-seed-bridge-in-model](../../../../decision/parity-seed-bridge-in-model.md)

### In

- depends-on [check-whatif-once](check-whatif-once.md)
- verifies [model-seed-views-verified](../../../../evidence/model-seed-views-verified.md)
- verifies [model-seed-views-qa-1](../../../../evidence/model-seed-views-qa-1.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
