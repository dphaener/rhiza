# model-seed-views-qa-1

evidence · `01M7H08YRAM8NN4W5DH9WKEHXE` · actor [claude-code-b257b28a](../actor/claude-code-b257b28a.md)

## result

PASS

## statement

QA exercised AsIsSeed, ToBeFrom and WhatIfFrom through the built rhiza check (3bcbde2fe): on the scratch repo with dependent, stale and bad-link in-flight change sets and on ~/code/rhiza, every in-flight and to-be finding and the whole --json report matched main 5180b25eb (which uses ToBe and WhatIf) byte for byte, so the seed views gave the same models and problems.

## Links

### Out

- verifies [model-seed-views](../capability/rhiza/integrity/check/model-seed-views.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
