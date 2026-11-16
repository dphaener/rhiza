# check-whatif-once-verified

evidence · `01M7GYHS1ESBYB698ZAZDF9GZ5` · actor [claude-code-ea1ae9dd](../actor/claude-code-ea1ae9dd.md)

## result

PASS

## statement

At 89905bae9 (capped scope): go test -trimpath -p 2 -count=1 ./internal/check exit 0; TestWhatIfOnceMatches/Listed (landed base, dependent A/B clean, stale op, duplicate handle, bad in-flight link, broken text ref, listed exactly) and /Once (stale op, bad link, broken ref each reported once on A, never on dependents B, C) pass; go test ./cmd/rhiza -run 'TestScripts/(check|check-.*|text-refs-.*)$' with RHIZA_/SHIP_/REX_ unset exit 0. Built base 6f8de2300 and this commit; in ~/code/rhiza, rhiza check text and --json output byte-identical (cmp), both exit 0. pprof -top -cum: base model.WhatIf 45.5% and model.asIs 24.8% under check.Run; after, model.WhatIf and model.ToBe are absent under check.Run (focus/show filter empty), model.WhatIfFrom 21.7%, model.asIs 0.47% called only from model.AsIsSeed (-peek), which check.go calls once in checker.model. Diff touches only internal/check/check.go, textrefs.go, whatif_test.go plus the builder's actor record and render; no TODOs.

## Links

### Out

- verifies [check-whatif-once](../capability/rhiza/integrity/check/check-whatif-once.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
