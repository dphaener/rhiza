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

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
