# daemon-passes-cpuprofile-to-drivers

capability · `01M78V4WDA98WKD9YNVGJWBWA7`

## acceptance

go test -trimpath ./internal/api/... -run CPUProfile: a Go test starts a run through the daemon with a fake Start and RHIZA_CPUPROFILE=<dir>/daemon.pprof in the daemon's environment, and asserts the driver's environment carries RHIZA_CPUPROFILE=<dir>/daemon.pprof.<run id>; with RHIZA_CPUPROFILE unset the driver's environment has none

## statement

When the daemon runs with RHIZA_CPUPROFILE set, each run driver it starts gets its own RHIZA_CPUPROFILE, the daemon's path with the run's ID added, so every driver writes a start-to-exit CPU profile without overwriting the daemon's or another driver's

## Links

### Out

- decided-by [decision-idle-cpu-scope](../decision/decision-idle-cpu-scope.md)

### In

- verifies [daemon-passes-cpuprofile-to-drivers-verified](../evidence/daemon-passes-cpuprofile-to-drivers-verified.md)
- verifies [daemon-passes-cpuprofile-to-drivers-qa-1](../evidence/daemon-passes-cpuprofile-to-drivers-qa-1.md)
- verifies [pr-511-fix-the-root-daemon-still](../evidence/pr-511-fix-the-root-daemon-still.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
