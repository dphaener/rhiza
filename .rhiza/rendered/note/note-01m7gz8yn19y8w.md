# note-01m7gz8yn19y8w

note · `01M7GZ8YN19Y8WFB22FPDXPY0H`

## statement

read-entries-parallel's 5 s clause is unmet, not passed (Decision check-time-gap-carried): rhiza check went from 20.6 s at base to 7.3-11.7 s on the shared VM at load 2-3 (verifier 11.67 s wall, 16.25 s user + 4.14 s sys at e4d148ed5). The gap to 5 s is carried by draft change set 01M7GZ8K2B (Check within 5 s); the target stays as written.

## Links

### Out

- discusses change set `01M7GFT8ADJN38YJCZ3ZGDJ6YH` (rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
