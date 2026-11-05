# milestone-22

intent · `01M6CS0E57Z80HJ5ZB1CQRE466`

## milestone

22

## statement

Run Rhiza from Rex: from the TUI apps alone the owner sees the whole system's health, files bugs and requests features that the daemon plans and ships, and switches each stage (planning, acceptance, shipping, sign-off) and the daemon itself between autonomous and owner-driven, with no orchestrator session

## Links

### Out

- decided-by [m22-run-from-rex](../decision/m22-run-from-rex.md)

### In

- serves change set `01M6CRZQ0YKB2QV4V5TWY71G4V` (The daemon reads its settings and each stage's autonomy from git config every poll, can be paused, and applies an acceptance policy)
- serves change set `01M6CRZQYG5PJV7X2B234DJCV4` (Feature, bug and milestone requests queue in the clone and the daemon turns each into a plan run)
- serves change set `01M6CRZRVGZBRB51YFPXHKZA75` (A requests app files features, bugs and milestones and shows each request's progress)
- serves change set `01M6CRZSRFQ94MAJ8CRTEE44P2` (rhiza health and a health app show the whole system's state on one screen)
- serves change set `01M6CRZTNDTPF1SD2WGBY4835Y` (An autonomy app edits every setting, starts or stops the daemon, and the inbox nudges stalled sessions)
- serves change set `01M6D4SB0D8F9SPGK471DF9RV6` (The request queue and the daemon turning requests into plan runs, as a leaf of plan requests)
- serves change set `01M6MY8BMPFK9ZWD3TDQHKT64J` (Milestone planning always needs a human: a plan run for a milestone (a record change set, or a request of kind milestone) never runs autonomously, even when the daemon or run is autonomous (rhiza.autonomy.plan=auto); it waits for the owner's interview instead of the arbiter answering. Owner, 2026-11-05: "milestone planning should never be autonomous even in autonomous mode. that always requires a human.")

## History

- `01M6CRZP2E8H6EEEPVDAWAKQ6H` Plan milestone 22: run Rhiza from Rex, the owner's ask of 2026-11-01
