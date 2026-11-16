# model-seed-views-verified

evidence · `01M7GXZWCGPAJ7THF6GYYW046B` · actor [claude-code-17b3dfd5](../actor/claude-code-17b3dfd5.md)

## result

PASS

## statement

At 96cd07870: re-ran the Checks: block: capped go test -trimpath -p 2 -count=1 ./internal/model exit 0 (ok); rhiza ship generate exit 0 (0 changed); rhiza render exit 0 (only this verifier's own actor page written). go test -run TestAsIsSeedWhatIfFrom -v: rapid 100 tests passed, what-ifs 123, in-flight dependencies 13, skipped ops 273, duplicate handles 48 (non-vacuous). Diff: internal/model/seed.go adds ToBeFrom = seed.builder() (the existing copying builder, no second seed constructor) + layer(accepted) + finishModel(false), mirroring ToBe (asIs + layer(accepted) + finish); seed_test.go calls ToBeFrom on the shared seed before every WhatIfFrom, compares it with ToBe (model and problems via sameModel) only after all WhatIfFrom calls ran, and the existing fresh-AsIsSeed DeepEqual check of the seed still follows, so aliasing between ToBeFrom, the what-ifs and the seed would fail. Only seed.go, seed_test.go and the builder's actor record/render changed; no new rule, no existing function changed, no TODOs.

## Links

### Out

- verifies [model-seed-views](../capability/rhiza/integrity/check/model-seed-views.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
