# projection-scope-approved

decision · `01M7GFT06MX3HQ7M70CHYJ0Y85`

## question

Scope: 1–4 in this milestone, 5 and 6 later?

## quote

Yes

## quote_where

owner, 2026-11-15 chat, answer to question 1

## statement

Milestone 39 holds four change sets, in this order: a parity harness with internal/api's tests on the index; rhiza check without the repeated what-if rebuilds; finishing the projection (one projection interface, the index its only implementation with pending entries overlaid, the reads moved onto it and the interim model deleted); and rhiza run in-process in cmd/rhiza's scripts. Moving the daemon, run, ship and web application code out of internal/api, and finding what fills the 127 MB log, come later

## Links

### Out

- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

### In

- decided-by [milestone-39](../intent/milestone-39.md)
- decided-by [projection-parity-property](../capability/rhiza/views/fast-reads/projection-parity-property.md)
- decided-by [api-tests-on-index](../capability/rhiza/integrity/self-check/api-tests-on-index.md)
- decided-by [check-whatif-once](../capability/rhiza/integrity/check/check-whatif-once.md)
- decided-by [read-entries-parallel](../capability/rhiza/integrity/check/read-entries-parallel.md)
- decided-by [projection-interface](../capability/rhiza/views/fast-reads/projection-interface.md)
- decided-by [render-from-projection](../capability/rhiza/views/fast-reads/render-from-projection.md)
- decided-by [change-reads-from-projection](../capability/rhiza/views/fast-reads/change-reads-from-projection.md)
- decided-by [daemon-fold-from-projection](../capability/rhiza/views/fast-reads/daemon-fold-from-projection.md)
- decided-by [cli-no-process-globals](../capability/rhiza/integrity/self-check/cli-no-process-globals.md)
- decided-by [cli-child-env](../capability/rhiza/integrity/self-check/cli-child-env.md)
- decided-by [run-clock-from-ctx](../capability/rhiza/integrity/self-check/run-clock-from-ctx.md)
- decided-by [run-concurrent-guard](../capability/rhiza/integrity/self-check/run-concurrent-guard.md)
- decided-by [workspace-log-once](../capability/rhiza/views/fast-reads/workspace-log-once.md)
- decided-by [slice-and-show-from-projection](../capability/rhiza/views/fast-reads/slice-and-show-from-projection.md)
- decided-by [scripts-in-process-harness](../capability/rhiza/integrity/self-check/scripts-in-process-harness.md)
- decided-by [scripts-in-process-a-d](../capability/rhiza/integrity/self-check/scripts-in-process-a-d.md)
- decided-by [scripts-in-process-e-q](../capability/rhiza/integrity/self-check/scripts-in-process-e-q.md)
- decided-by [scripts-in-process-r](../capability/rhiza/integrity/self-check/scripts-in-process-r.md)
- decided-by [scripts-in-process-s-w](../capability/rhiza/integrity/self-check/scripts-in-process-s-w.md)
- decided-by [scripts-in-process-docs](../capability/rhiza/integrity/self-check/scripts-in-process-docs.md)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
