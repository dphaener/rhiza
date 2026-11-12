# ci-wait-backs-off-verified-2

evidence · `01M7907QMJXYC28RNY4RESXBSY` · actor [claude-code-7e9d8b54](../actor/claude-code-7e9d8b54.md)

## result

PASS

## statement

At feb1c57a3 (review fix adding max_interval's KeyDoc to the KindWait spec and docs/reference/workflow-file.md): go test -trimpath ./internal/run/... ./internal/workflow/... -run 'WaitBackoff|MaxInterval' passed TestWaitBackoff (fake clock, scripted unchanged-then-changed output, sleeps 30s,60s,120s,120s,120s then 30s after the change), TestWaitBackoffStopsAtTimeout (sleeps sum to exactly the 5m timeout), TestWaitMaxInterval (parses 2m, rejects below interval and on non-wait steps) and TestMaxIntervalShipChangeCIWaits (built-in ship-change ship-check waits set 30s/2m). Full ./internal/workflow/... and ./internal/run/... pass. rhiza workflow check exits 0 on every .rhiza/workflows/*.toml; all seven rhiza ship check waits in both builtin and repo ship-change.toml (which bugfix includes) set interval 30s, max_interval 2m. Diff read: engine.go doubles the sleep only when output equals the previous check's and caps at the remaining timeout.

## Links

### Out

- verifies [ci-wait-backs-off](../capability/ci-wait-backs-off.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
