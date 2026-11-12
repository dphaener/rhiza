# root-daemon-idle-poll-no-work-verified

evidence · `01M78YYKTFBJWSHD2T1M8EX7CJ` · actor [claude-code-b94e70ae](../actor/claude-code-b94e70ae.md)

## result

PASS

## statement

Verifier at f86804689: go test -trimpath ./internal/api/... -run 'RootIdlePoll' -v -count=1 passed. TestRootIdlePoll (two projects, idle-poll fixture) asserts the second and third unchanged polls add zero log reads, plan reads, folds, picks (new daemonPickCount) and transcript reads, then changes a log entry, a ref, a setting, a run record and a transcript file in turn and asserts the next poll redoes the work and the poll after is idle again. TestRootIdlePollBudget (quiettest median, budgetLock, like TestReadBudget) measured median 69.7 ms, own CPU 52.7 ms, under the 100 ms budget. Diff read: daemonDemand memoizes daemon.Pick on fold identity plus a signature of runs (id, entries, liveness), clock-based holds, running count, unjournaled and workflow/dispatch file stamps, so holds and liveness are still evaluated each poll. Also go test -trimpath ./internal/api/ -run 'Daemon|Pick|Watchdog|Usage' passed, go vet and gofmt clean.

## Links

### Out

- verifies [root-daemon-idle-poll-no-work](../capability/root-daemon-idle-poll-no-work.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
