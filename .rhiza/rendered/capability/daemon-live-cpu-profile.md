# daemon-live-cpu-profile

capability · `01M78V4WAQ555KCQSZVXT3CPDQ`

## acceptance

go test -trimpath ./internal/api/... -run DaemonProfile: a Go test runs the daemon loop in-process with a fake clock, issues the profile request with --seconds 1, and asserts the daemon writes a non-empty gzip-compressed pprof file at --out (or the default path in the daemon's state folder), keeps polling afterwards, and that a second request while one is in progress is refused with an error naming the capture in progress; no test binds a socket

## statement

A hidden 'rhiza daemon profile [--seconds N] [--out <file>]' (N defaults to 30) asks the running daemon to capture a CPU profile of itself for N seconds and write it to the file, without restarting it, and prints the file's path

## Links

### Out

- decided-by [decision-idle-cpu-scope](../decision/decision-idle-cpu-scope.md)

### In

- verifies [daemon-live-cpu-profile-verified](../evidence/daemon-live-cpu-profile-verified.md)
- verifies [daemon-live-cpu-profile-qa-1](../evidence/daemon-live-cpu-profile-qa-1.md)
- verifies [pr-511-fix-the-root-daemon-still](../evidence/pr-511-fix-the-root-daemon-still.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
