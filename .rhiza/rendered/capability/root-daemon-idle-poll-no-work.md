# root-daemon-idle-poll-no-work

capability · `01M78V4WM0TQZYN29GP2GCGGPZ`

## acceptance

go test -trimpath ./internal/api/... -run 'RootIdlePoll': a Go test runs the root daemon poll over two projects twice with nothing changed and asserts through the readLog, planRead and fold counters that the second poll folds nothing and computes no pick; it then changes each stamp in turn (log entry, ref move, run record, setting, transcript file) and asserts the next poll does the work again; a CPU budget test (quiettest, median, like TestReadBudget) on the idle-poll-benchmarks fixture asserts an unchanged root poll costs under 100 ms of process CPU

## statement

A root daemon poll in which, for every registered project, the log, git refs, run records, settings and usage inputs are unchanged by their stamps does no fold and no pick and costs at most 0.05 core of steady CPU at the 5 s poll, while clock-based holds, process liveness and due schedules are still evaluated

## Links

### Out

- depends-on [idle-poll-benchmarks](idle-poll-benchmarks.md)
- depends-on [usage-scan-once-per-poll](usage-scan-once-per-poll.md)
- decided-by [decision-idle-cpu-scope](../decision/decision-idle-cpu-scope.md)
- supersedes [daemon-idle-poll-no-fold](rhiza/orchestration/daemon/daemon-idle-poll-no-fold.md)

### In

- verifies [root-daemon-idle-poll-no-work-verified](../evidence/root-daemon-idle-poll-no-work-verified.md)
- verifies [root-daemon-idle-poll-no-work-qa-1](../evidence/root-daemon-idle-poll-no-work-qa-1.md)
- verifies [root-daemon-idle-poll-no-work-verified-2](../evidence/root-daemon-idle-poll-no-work-verified-2.md)
- verifies [pr-511-fix-the-root-daemon-still](../evidence/pr-511-fix-the-root-daemon-still.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
