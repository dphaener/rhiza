# idle-poll-benchmarks

capability · `01M78V4WH7H4JKPF1CZ9S33Z7M`

## acceptance

go test -trimpath ./internal/api/ -run '^$' -bench 'IdleRootPoll|IdleDriverPoll' -benchtime 3x -cpuprofile $TMPDIR/idle.pprof completes and writes a profile; the benchmarks build their fixture in a temp dir with no network, Rex, agent or gh

## statement

In-process Go benchmarks time one idle root-daemon poll (two registered projects, nothing changed) and one idle driver poll over a production-sized fixture (a log of 1,000+ change sets, several runs in flight, ~4,000 transcript files), so a builder profiles idle cost with go test -cpuprofile and never needs the VM

## Links

### Out

- decided-by [decision-idle-cpu-scope](../decision/decision-idle-cpu-scope.md)

### In

- depends-on [root-daemon-idle-poll-no-work](root-daemon-idle-poll-no-work.md)
- verifies [idle-poll-benchmarks-verified](../evidence/idle-poll-benchmarks-verified.md)
- verifies [idle-poll-benchmarks-qa-1](../evidence/idle-poll-benchmarks-qa-1.md)
- verifies [pr-511-fix-the-root-daemon-still](../evidence/pr-511-fix-the-root-daemon-still.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
