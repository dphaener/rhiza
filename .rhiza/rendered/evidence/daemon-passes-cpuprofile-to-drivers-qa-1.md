# daemon-passes-cpuprofile-to-drivers-qa-1

evidence · `01M78ZSA2M3KG59PFW928ES6CA` · actor [claude-code-a1181376](../actor/claude-code-a1181376.md)

## result

PASS

## statement

QA at 4f1ad9a09: a real root daemon run with RHIZA_CPUPROFILE=<dir>/daemon.pprof and a dispatch rule sending a ready feature to a scratch command-only workflow started a driver ('run start --foreground ... --workflow backoff'); /proc/<driver>/environ held RHIZA_CPUPROFILE=<dir>/daemon.pprof.01M78ZF6P7JDC1MEXJH2X74VM7.start-1, the driver wrote that 1.7 KB profile on exit and the daemon wrote its own daemon.pprof (311 s) on SIGINT, neither overwriting the other. The suffix is the change set's start log name, not the run ID (01M78ZF88X...) the statement names.

## Links

### Out

- verifies [daemon-passes-cpuprofile-to-drivers](../capability/daemon-passes-cpuprofile-to-drivers.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
