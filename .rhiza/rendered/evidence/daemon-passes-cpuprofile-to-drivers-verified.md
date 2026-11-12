# daemon-passes-cpuprofile-to-drivers-verified

evidence · `01M78W44HZJZJJWB7WT52TQD2T` · actor [claude-code-cf888442](../actor/claude-code-cf888442.md)

## result

PASS

## statement

At de5efa519 ran go test -trimpath ./internal/api/... -run CPUProfile -v -count=1: TestDaemonCPUProfilePassedToDrivers PASS. The test drives daemonResume (the daemon starting a run's driver) with a fake Detach and Getenv; with RHIZA_CPUPROFILE=/dir/daemon.pprof the driver Program's Env holds RHIZA_CPUPROFILE=/dir/daemon.pprof.01RUNID, and with it unset no RHIZA_CPUPROFILE entry. Read the diff: daemonDetachNoRex, which every daemon-detached driver goes through, appends RHIZA_CPUPROFILE=<daemon path>.<name>; Program.Env is set over the inherited environment, so it replaces the daemon's own value and the driver never overwrites the daemon's profile. For a start, name is the start log (<change set>.start-N) since the run ID is drawn later by the child, still a file unique to that driver. go vet and gofmt clean on internal/api.

## Links

### Out

- verifies [daemon-passes-cpuprofile-to-drivers](../capability/daemon-passes-cpuprofile-to-drivers.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
