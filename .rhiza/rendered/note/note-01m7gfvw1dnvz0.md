# note-01m7gfvw1dnvz0

note · `01M7GFVW1DNVZ003XNJCDK57VS`

## statement

Parked for after milestone 39 (owner: "Yes" to these coming later): move the application code out of internal/api. The daemon (daemon*.go, 6.4k lines; runDaemon 366 lines at daemon.go:94, daemonPoll1 267 at :1005), the run engine (run*.go, 10.6k), shipping (ship*.go, 5.1k) and the web handlers (web_*.go, which make internal/api import internal/web and tui/launcher) go into domain packages that depend on the projection interface this milestone defines, so api is the transport-independent registry again. It waits for that interface. Survey: ~/orchestrator/architecture-survey-2026-11-15.md, sections 1 and 4.

## Links

### Out

- discusses [milestone-39](../intent/milestone-39.md)

## History

- `01M7GFRBHB4ZXKD4882FQ480Q3` Reads come from the projection. Owner, 2026-11-15, after the orchestrator measured this repo's log (24,560 entries, 127 MB: rhiza status 2.0 s, inbox 2.3 s, check 15.3 s) and ran an architecture survey (~/orchestrator/architecture-survey-2026-11-15.md): reads and writes come from the incremental index instead of folding the whole log, under one set of model rules with the strict engine as the spec; four change sets (parity harness, check without repeated what-if rebuilds, finishing the projection, rhiza in-process in scripts), the application-layer move and the log-size investigation later; milestone 38's cold and CI targets carried here.
