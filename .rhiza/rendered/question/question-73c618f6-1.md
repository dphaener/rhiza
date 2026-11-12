# question-73c618f6-1

question · `01M78V019YZ361DP3Y8XQVZAMD`

## session

73c618f6-89e6-4e1f-881c-01a1d1add530

## statement

Bug 01M78TV0 (root daemon ~0.43 core idle, drivers up to ~0.15 core): which intent should it serve? Both earlier CPU fixes (01M6H6SG daemon fold cache, 01M6NE0X driver idle poll) landed and served m24-on-the-vm / milestone-25, neither on the roadmap now; open milestone 4's m4-orchestrated-by-rhiza is what the daemon load-hold fix (01M75T74) served. / How should the live daemon and drivers be profiled at idle? Today RHIZA_CPUPROFILE only profiles a process from start to exit, so the root service can't be profiled without a restart. / What should the idle CPU budgets and the CI-wait cadence be? The usage governor re-walks ~4,000 Claude transcript files (2.6 GB dir) each 5 s poll, once per project; ship-change's CI waits run a fresh 'rhiza ship check' process every 30 s.

## Links

### In

- answers [decision-idle-cpu-scope](../decision/decision-idle-cpu-scope.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
