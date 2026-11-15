# projection-planned-autonomously

decision · `01M7GFT0A4PH0CG1P11JQGQEP8`

## question

Plan in chat or autonomous?

## quote

Autonomous

## quote_where

owner, 2026-11-15 chat, answer to question 3

## statement

Milestone 39 is planned autonomously: the orchestrator writes its change sets and leaves from the architecture survey and the owner's answers, deciding what is left as arbiter, without a further interview

## Links

### Out

- derived-from [projection-milestone-interview](../source/projection-milestone-interview.md)

### In

- decided-by [milestone-39](../intent/milestone-39.md)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
