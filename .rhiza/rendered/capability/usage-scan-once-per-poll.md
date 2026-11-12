# usage-scan-once-per-poll

capability · `01M78V4WPDAC0P0JMCA17H4TC9`

## acceptance

go test -trimpath ./internal/api/... ./internal/adapter/claudecode/... -run 'UsageScan': a Go test with a fake clock and a transcript tree of 100+ files counts directory walks and file reads across root polls of two projects: one walk per poll at most, none in a poll less than a minute after the last full walk, zero file reads when nothing changed, one read when one known file grows, and a new file found by the next full walk; the usage windows computed match those of a full read

## statement

The usage governor reads Claude Code transcripts once per root daemon poll for all projects, not once per project, re-reads only transcript files whose size or mtime changed, and walks the transcript tree in full at most once a minute, re-checking only the files it already knows between full walks

## Links

### Out

- decided-by [decision-idle-cpu-scope](../decision/decision-idle-cpu-scope.md)

### In

- depends-on [root-daemon-idle-poll-no-work](root-daemon-idle-poll-no-work.md)
- verifies [usage-scan-once-per-poll-verified](../evidence/usage-scan-once-per-poll-verified.md)
- verifies [usage-scan-once-per-poll-qa-1](../evidence/usage-scan-once-per-poll-qa-1.md)
- verifies [pr-511-fix-the-root-daemon-still](../evidence/pr-511-fix-the-root-daemon-still.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
