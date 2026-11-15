# m39-ci

intent · `01M7GFT01P2FVWJ2J324EFT8N9` · parent [milestone-39](../milestone-39.md)

## statement

The median wall time of the last 10 completed CI runs on pull requests that changed code is within 3 minutes, carried from @m38-ci (2026-11-15: 8 runs since #588, median 181 s). Measure: gh run list --workflow CI --event pull_request --status completed --limit 10 --json createdAt,updatedAt, and the median of updatedAt minus createdAt.

## Links

### Out

- decided-by [m38-ci-target-carried](../../decision/m38-ci-target-carried.md)

### In

- serves change set `01M7GFV2ZMHRJR78H0YTMKWXB7` (A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.)
- serves change set `01M7GFT8ADJN38YJCZ3ZGDJ6YH` (rhiza check builds the as-is model once and each in-flight change set's what-if once (model.AsIsState and model.WhatIfFrom), shared between its two passes, and store.ReadEntries decodes the log on several cores, so check on this repo takes seconds, not 15. Evidence: check calls model.WhatIf twice per in-flight change set (internal/check/check.go:384 and :456), about 96 as-is rebuilds, over a serial decode of 127 MB (internal/store/read.go:97); not yet profiled.)
- serves change set `01M7GFTS0NVWBXF3CQZD916J1C` (cmd/rhiza's testscripts run rhiza in-process: run() takes its working directory and its own path as parameters instead of os.Getwd and os.Args, and only scripts whose fake bin/ commands call rhiza back get a real binary on PATH, so a script line no longer costs a process start; AGENTS.md says exactly which lines run in-process. Evidence: testscript.RunMain re-executes the test binary for each of about 5,330 rhiza lines in 398 scripts; a cold make test spends 397 of its 967 CPU-s in the kernel.)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
