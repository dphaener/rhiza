# milestone-39

intent · `01M7GFT007N7A6EPR12Z16GECM`

## milestone

39

## statement

Reads come from the projection: Rhiza's reads and writes come from the incremental index (the projection) instead of folding the whole log, under one set of model rules with the strict engine in internal/model as the spec, so everyday commands stay fast as the log grows and the test suite stops paying for full-log rebuilds and per-line process starts (2026-11-15, 24,560 entries, 127 MB: rhiza status 2.0 s, rhiza inbox 2.3 s, rhiza check 15.3 s; cold make test 423 s).

## Links

### Out

- decided-by [projection-scope-approved](../decision/projection-scope-approved.md)
- decided-by [strict-engine-is-spec](../decision/strict-engine-is-spec.md)
- decided-by [projection-planned-autonomously](../decision/projection-planned-autonomously.md)
- decided-by [projection-build-order](../decision/projection-build-order.md)

### In

- discusses [note-01m7gfvw1dnvz0](../note/note-01m7gfvw1dnvz0.md)
- discusses [note-01m7gfvwtz3mcd](../note/note-01m7gfvwtz3mcd.md)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
