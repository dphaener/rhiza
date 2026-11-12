# ci-wait-backs-off-verified

evidence · `01M78Y66T7ZQ71EPW0226CTF5C` · actor [claude-code-f2d57b50](../actor/claude-code-f2d57b50.md)

## result

PASS

## statement

At 6629e1cb2: go test -trimpath ./internal/run/... ./internal/workflow/... -run 'WaitBackoff|MaxInterval' -v passed TestWaitBackoff (fake clock, scripted outputs a a a a a b b; sleeps 30s 60s 120s 120s 120s 30s 60s), TestWaitBackoffStopsAtTimeout (total sleep exactly the 5m timeout), TestWaitMaxInterval (parses 2m, rejects 10s below 30s interval, rejects on a command step) and TestMaxIntervalShipChangeCIWaits (every built-in ship-change ship-check wait has 30s/2m); the full internal/run and internal/workflow packages pass. A build of this commit ran rhiza workflow check on each .rhiza/workflows/*.toml (bugfix, docs*, plan-bug, ship-change): all exit 0. Read the diff: all 7 ship-check waits in both the repo and built-in ship-change set max_interval 2m (bugfix includes ship-change); the production Condition runner returns the check's output tail, so backoff resets on changed output.

## Links

### Out

- verifies [ci-wait-backs-off](../capability/ci-wait-backs-off.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
