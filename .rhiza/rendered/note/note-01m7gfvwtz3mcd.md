# note-01m7gfvwtz3mcd

note · `01M7GFVWTZ3MCD9JM3CAXMYZZN`

## statement

Parked for after milestone 39 (owner: "Yes" to these coming later): investigate what fills the log. This repo's .rhiza/log is 127 MB for 24,560 entries, about 5 KB an entry; which entry types and fields (findings, sources, transcripts, run records) fill it decides whether an index's first build and rhiza check can ever be fast, and whether entries should be smaller. Investigation first, measured with the sizes per entry type and field; no change set until it is done.

## Links

### Out

- discusses [milestone-39](../intent/milestone-39.md)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
