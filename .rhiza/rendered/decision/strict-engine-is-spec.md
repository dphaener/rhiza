# strict-engine-is-spec

decision · `01M7GFT08KZFZSZTD7G35JHPTJ`

## question

When the parity test finds a disagreement, the strict engine (internal/model, which the index uses) is the spec?

## quote

Ok

## quote_where

owner, 2026-11-15 chat, answer to question 2

## statement

Where the parity test finds the model rules disagree, the strict engine in internal/model (the rules the index stores) is the spec: the other side (the interim model in internal/api/model.go or the hybrid in write_indexed.go) changes to match, and each case is recorded

## Links

### Out

- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

### In

- decided-by [milestone-39](../intent/milestone-39.md)
- decided-by [m39-parity](../intent/milestone-39/m39-parity.md)
- decided-by [projection-parity-property](../capability/rhiza/views/fast-reads/projection-parity-property.md)
- decided-by [projection-interface](../capability/rhiza/views/fast-reads/projection-interface.md)
- decided-by [parity-interim-from-strict](../capability/rhiza/views/fast-reads/parity-interim-from-strict.md)
- decided-by [parity-write-index-from-strict](../capability/rhiza/views/fast-reads/parity-write-index-from-strict.md)
- decided-by [index-as-is-model](../capability/rhiza/views/fast-reads/index-as-is-model.md)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
