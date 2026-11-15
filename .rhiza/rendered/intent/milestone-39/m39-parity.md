# m39-parity

intent · `01M7GFT044X4B4XPAZEFYDY26G` · parent [milestone-39](../milestone-39.md)

## statement

A parity property test shows that the projection and a fresh strict fold of the log (changeset.FoldLog, then internal/model's views) agree on random logs, including two landed heads on one target, a rebase that drops entry files, and dry-run pending entries; it runs in make test, not only in the opt-in budgets. Measure: go test -trimpath -count=1 ./internal/api -run '^TestPropertyProjectionParity$' -v passes, and make test's output lists it.

## Links

### Out

- decided-by [projection-targets](../../decision/projection-targets.md)
- decided-by [strict-engine-is-spec](../../decision/strict-engine-is-spec.md)

### In

- serves change set `01M7GFV2ZMHRJR78H0YTMKWXB7` (A parity property test checks that the interim model, the index's write model and a fresh strict fold of the log agree on random logs (two landed heads, rebases dropping entry files, dry-run pending entries), each disagreement is fixed on the non-strict side with the strict engine as the spec, and internal/api's tests run on the index, as production does. Evidence: the model rules are implemented three times (api/model.go, internal/model, write_indexed.go); testCtx sets no IndexPath, so all ~770 api tests run the interim path.)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
