# projection-targets

decision · `01M7GFT0BD4S1H8G036DW7ND7S`

## reasoning

rhiza version starts in 14 ms, so 200 ms leaves room to open the index and read a few rows, which is all status and inbox need once they stop folding the 127 MB log; a looser target would let a hidden fold survive. check must read and decode the whole log and compare the rendered files, so it cannot be a few rows: with the as-is built once instead of about 96 times and the decode spread over 4 cores, 5 s is a third of today's 15.3 s and stays below the minutes self-check, the gate and pre-push spend elsewhere. Each command is one a reviewer can rerun.

## statement

Milestone 39's own targets are rhiza status and rhiza inbox within 200 ms warm and rhiza check within 5 s on this repo's log, each measured with /usr/bin/time after one warming run in the main checkout on a quiet VM; the parity test runs in make test; the interim model's removal is measured by a grep that prints nothing; and no test is lost against a98fd4051

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

### In

- decided-by [m39-reads-fast](../intent/milestone-39/m39-reads-fast.md)
- decided-by [m39-check-fast](../intent/milestone-39/m39-check-fast.md)
- decided-by [m39-parity](../intent/milestone-39/m39-parity.md)
- decided-by [m39-one-read-path](../intent/milestone-39/m39-one-read-path.md)
- decided-by [m39-nothing-lost](../intent/milestone-39/m39-nothing-lost.md)
- decided-by [read-entries-parallel](../capability/rhiza/integrity/check/read-entries-parallel.md)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
