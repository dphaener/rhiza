# m39-check-fast

intent · `01M7GFT039RM7ZP4FVPP57TB46` · parent [milestone-39](../milestone-39.md)

## statement

rhiza check finishes within 5 s on this repo (2026-11-15: 15.3 s). Measure, in the main checkout with nothing else running: run it once, then /usr/bin/time -f '%e s' rhiza check.

## Links

### Out

- decided-by [projection-targets](../../decision/projection-targets.md)

### In

- serves change set `01M7GFT8ADJN38YJCZ3ZGDJ6YH` (rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
