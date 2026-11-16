# check-time-gap-carried

decision · `01M7GZ7H98QNRXCHTGS72ASM1D`

## reasoning

check-time-target says a miss fails the leaf and the gap is raised, not the target loosened, so the target stays 5 s and the miss is not called a pass. But the leaf cannot close the gap by itself: it touches only internal/store/read.go, the decode it parallelised was about 1.3 s of CPU (7% in the 2026-11-15 profile), and the measured 16.25 s user + 4.14 s sys is at least 5.1 s even spread perfectly over 4 cores, so no retry inside the leaf's touch list could reach 5 s. What remains is WhatIfFrom (~19%), checker.model/inFlightLinks and render, all outside this leaf's scope. Failing the whole change set would throw away verified work that roughly halves check's time and matches every correctness point (model-seed-views and check-whatif-once both verified PASS). Raising the gap as its own change set keeps the target binding on milestone 39 while letting this change set's correct work land.

## statement

read-entries-parallel is settled on the rest of its acceptance (TestReadEntriesParallelSame, -race, GOMAXPROCS=1, same order and content, no dependency, the Checks: block timings), and its 5 s clause is recorded as unmet, not passed: the change set lands with the measured gain (20.6 s at base to 7.3–11.7 s), the 5 s target in projection-targets and check-time-target stays exactly as written, and the conductor drafts a discovered-from change set in milestone 39 that carries the remaining gap (model.WhatIfFrom, checker.model/inFlightLinks and the render comparison, profiled at e4d148ed5) to 5 s. No retry of read-entries-parallel is spent on the target

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-903fe4d8-1](../question/question-903fe4d8-1.md)

### In

- decided-by [read-entries-parallel-verified](../evidence/read-entries-parallel-verified.md)

## History

- `01M7GFT8ADJN38YJCZ3ZGDJ6YH` rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.
