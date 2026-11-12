# daemon-live-cpu-profile-verified

evidence · `01M78W18Z0G87QBBQV7G2CHFWF` · actor [claude-code-0ff2bdb5](../actor/claude-code-0ff2bdb5.md)

## result

PASS

## statement

At 65b2e01fe, go test -trimpath -count=1 ./internal/api/... -run 'DaemonProfile|Describe' -v passes (also under -race): TestDaemonProfile runs daemonLoop in-process with test-driven ticks and a fake capture clock, issues daemon.profile with seconds 1 and --out, asserts the 1s capture timer, a non-empty gzip pprof at --out, polls during and after the capture, and a second request refused with 'a profile capture is already in progress'; TestDaemonProfileDefaultPath checks the default path lands in the root's state folder; TestDaemonProfileNoDaemon checks the no-daemon error. The request and result travel as files in .rhiza-root, so no socket is bound. go test -trimpath ./cmd/rhiza -run 'TestScripts/(describe|help|docs|daemon|usage|reference)' passes. A scratch build hides the op from 'rhiza daemon --help' while 'rhiza daemon profile --help' shows --seconds (default 30) and --out. Diff read: Hidden flag filters help, subcommands, docs, check and unfiltered describe.

## Links

### Out

- verifies [daemon-live-cpu-profile](../capability/daemon-live-cpu-profile.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
