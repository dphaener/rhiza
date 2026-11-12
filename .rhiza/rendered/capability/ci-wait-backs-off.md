# ci-wait-backs-off

capability · `01M78V4WSRZERXA7TRPH50HSAV`

## acceptance

go test -trimpath ./internal/run/... ./internal/workflow/... -run 'WaitBackoff|MaxInterval': a Go test drives a wait step with a fake clock and a scripted check whose output stays the same then changes, and asserts the sleeps go 30s, 60s, 120s, 120s and back to 30s after the change, never past the step's timeout; a workflow test parses max_interval, rejects one below interval, and asserts the built-in ship-change CI waits set 2m; rhiza workflow check passes on .rhiza/workflows

## statement

A wait step can back off: with max_interval set, each check whose output is unchanged from the last doubles the interval up to max_interval, and any change in the output brings it back to interval; ship-change's and bugfix's CI waits (rhiza ship check --settled ...) set interval 30s and max_interval 2m, so a driver waiting on pending CI runs a ship check at most every 2 minutes

## Links

### Out

- decided-by [decision-idle-cpu-scope](../decision/decision-idle-cpu-scope.md)

### In

- verifies [ci-wait-backs-off-verified](../evidence/ci-wait-backs-off-verified.md)
- verifies [ci-wait-backs-off-qa-1](../evidence/ci-wait-backs-off-qa-1.md)
- verifies [ci-wait-backs-off-verified-2](../evidence/ci-wait-backs-off-verified-2.md)
- verifies [pr-511-fix-the-root-daemon-still](../evidence/pr-511-fix-the-root-daemon-still.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
