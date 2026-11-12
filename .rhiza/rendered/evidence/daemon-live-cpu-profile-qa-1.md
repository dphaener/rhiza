# daemon-live-cpu-profile-qa-1

evidence · `01M78ZSA158C2A7TC8TGMXNS2P` · actor [claude-code-a1181376](../actor/claude-code-a1181376.md)

## result

PASS

## statement

QA at 4f1ad9a09 with a scratch build, on a scratch root (two projects) outside ~/code running a real 'rhiza daemon': 'rhiza daemon profile --seconds 5 --out live1.pprof' printed the path and exit 0; the file is gzip pprof (go tool pprof: Duration 5s, 310ms samples), the daemon kept polling after it; a second request issued 2s later exited 3 with 'a profile capture is already in progress'. With no daemon it exits 2 'no daemon is running at <root>'. Noted: on a daemon started with RHIZA_CPUPROFILE the capture fails with 'cpu profiling already in use', and on a cold start the request waits behind the first poll (94 s) so the CLI timed out at 23 s.

## Links

### Out

- verifies [daemon-live-cpu-profile](../capability/daemon-live-cpu-profile.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
