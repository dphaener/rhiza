# idle-cpu-plan-interview

source · `01M78V40HNZXEHTCAZ6AN9G789`

## quotes

```json
[]
```

## raw_sha256

348affd721a6bc946e9f556df916b28f1549fcd4e9f0694756a345aae8d904ac

## session

73c618f6-89e6-4e1f-881c-01a1d1add530

## summary

Autonomous plan-bug run for change set 01M78TV0 (no owner present; the arbiter answered). Facts gathered: the live root daemon had used 48,960 CPU-s over 115,135 s uptime (about 0.43 core) and one run driver about 0.15 core, matching the 2026-11-12 sampler report. The earlier fixes 01M6H6SG (daemon poll fold cache) and 01M6NE0X (driver idle poll from the index) have both landed, so the remaining cost lies elsewhere. Candidates found in the code: the usage governor (usageWindowsNow -> claudecode UsageRecords) walks ~4,000 transcript files under ~/.claude/projects (2.6 GB) on every 5 s poll, once per registered project; ship-change's CI wait steps run a fresh 'rhiza ship check --settled' process every 30 s; RHIZA_CPUPROFILE only profiles a process from start to exit, so the running root service cannot be profiled without a restart. One question (intent, profiling method, budgets and cadence) was raised to the arbiter, whose Decision decision-idle-cpu-scope settled: serve m4-orchestrated-by-rhiza; add 'rhiza daemon profile [--seconds N] [--out file]' for a live daemon capture, pass RHIZA_CPUPROFILE down to drivers, and add in-process benchmarks of idle daemon and driver polls; an unchanged daemon poll does no fold and no pick, with a budget test; targets of at most 0.05 core idle for the daemon and 0.01 core per driver; transcripts scanned once per poll for all projects, only changed files re-read, full walk at most once a minute; CI waits back off from 30 s to at most 2 minutes while CI is pending and return to 30 s on any state change.

## title

Plan interview: idle CPU of the root daemon and run drivers (autonomous)

## Links

### In

- derived-from [decision-idle-cpu-scope](../decision/decision-idle-cpu-scope.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
