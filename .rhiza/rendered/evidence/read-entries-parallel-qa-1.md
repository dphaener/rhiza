# read-entries-parallel-qa-1

evidence · `01M7H08YTCFWCW4MM3FSPT7NWV` · actor [claude-code-b257b28a](../actor/claude-code-b257b28a.md)

## result

PASS

## statement

QA ran the built rhiza check (3bcbde2fe) on a scratch log with an unreadable entry (chmod 000), an entry that does not decode and an entry copied under a non-matching file name: each was one fail log finding in path order, with the fold findings after, and human and --json output under default GOMAXPROCS and GOMAXPROCS=1 were byte-identical to main 5180b25eb. On ~/code/rhiza, check --json matched main byte for byte; main took 14.14 s, the branch 7.46 s (6.80 s human, 11.09 s at GOMAXPROCS=1); the 5 s miss is ruled (check-time-gap-carried).

## Links

### Out

- verifies [read-entries-parallel](../capability/rhiza/integrity/check/read-entries-parallel.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
