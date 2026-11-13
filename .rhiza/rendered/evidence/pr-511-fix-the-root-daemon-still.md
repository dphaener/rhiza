# pr-511-fix-the-root-daemon-still

evidence · `01M791BX2S33P3KE1BA86XK523` · actor [dphaener](../actor/dphaener.md)

## head

d164a0de0898c479a45a4c5470d1ad39724a0e98

## pr

https://github.com/dphaener/rhiza/pull/511

## result

13 checks passed: changes SUCCESS, go-install SUCCESS, model SKIPPED, lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #511 (https://github.com/dphaener/rhiza/pull/511) passed CI at d164a0d

## status

verified

## Links

### Out

- verifies [daemon-live-cpu-profile](../capability/daemon-live-cpu-profile.md)
- verifies [daemon-passes-cpuprofile-to-drivers](../capability/daemon-passes-cpuprofile-to-drivers.md)
- verifies [idle-poll-benchmarks](../capability/idle-poll-benchmarks.md)
- verifies [root-daemon-idle-poll-no-work](../capability/root-daemon-idle-poll-no-work.md)
- verifies [usage-scan-once-per-poll](../capability/usage-scan-once-per-poll.md)
- verifies [ci-wait-backs-off](../capability/ci-wait-backs-off.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
