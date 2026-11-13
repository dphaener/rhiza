# decision-idle-cpu-scope

decision · `01M78V2RH7Z677JJBPANNDM53X`

## reasoning

Intent: the two earlier CPU fixes served milestones that are now off the roadmap, so they can't be used. m4-orchestrated-by-rhiza (the daemon orchestrates everything with no orchestrator session) is open, and it is the intent that the daemon load-hold fix 01M75T74 served. The daemon's own ~0.45 core counts toward the load average that this hold reads. So idle daemon CPU directly reduces how much work the daemon-orchestrated loop can start. Profiling: cmd/rhiza/main.go only wraps a whole process in pprof (RHIZA_CPUPROFILE). Restarting the root service to profile it would lose the warm caches and the long-uptime state that the sampler measured, so the daemon needs an on-demand window. Drivers are per-run and short-lived, so start-to-exit profiling is enough for them. AGENTS.md forbids tests that touch the network or bind sockets, so the CI-safe path is a benchmark over a fixture. Budgets: the change set's intent asks for exactly this: skip fold and pick when inputs are unchanged, plus a CPU budget test for an idle poll. Counting folds is deterministic under load. A CPU ceiling follows the existing *_budget_test.go pattern (watchdog, inbox gate, write). The report puts daemon plus drivers at 44% of CPU and estimates about 0.7 core recoverable, which supports ~0.05 and ~0.01 core targets. Walking ~4,000 transcripts (2.6 GB) every 5 s once per project is the most obvious idle cost. Sharing one walk and keying it on size and mtime keeps the governor correct, and a once-a-minute full walk still catches new files. CI on this repo takes minutes, so checking every 30 s for the whole wait mostly costs process spawns. Backing off to 2 minutes, and returning to 30 s on any change, caps the added landing latency at about 2 minutes.

## statement

(1) Intent: change set 01M78TV0 serves m4-orchestrated-by-rhiza. (2) Profiling: add a hidden live capture for the running daemon, `rhiza daemon profile [--seconds N, default 30] [--out <file>]`, which asks the live daemon to run pprof.StartCPUProfile for N seconds and write the file, with no restart. Drivers are profiled by the existing RHIZA_CPUPROFILE start-to-exit hook, which the daemon passes down to the drivers it starts, plus an in-process Go benchmark of one idle daemon poll and one idle driver poll over a production-sized log fixture (go test -cpuprofile), so builders can profile deterministically without the VM. (3) Budgets and cadence: a daemon poll with nothing changed (log, refs, run records, settings, usage inputs all unchanged by stamp) does no fold and no pick; a budget test asserts zero folds on an unchanged poll and a CPU-time ceiling for it. Target steady idle cost is at most 0.05 core for the daemon and 0.01 core per driver. The usage governor scans Claude transcripts once per poll for all projects, not once per project, re-reads only files whose size or mtime changed, and does the full walk at most once a minute. While a ship-change step waits on CI, the ship check interval backs off from 30 s to at most 2 minutes while CI is pending, and goes back to 30 s on any state change.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-73c618f6-1](../question/question-73c618f6-1.md)
- supersedes [daemon-poll-caches-fold-not-picks](daemon-poll-caches-fold-not-picks.md)
- derived-from [idle-cpu-plan-interview](../source/idle-cpu-plan-interview.md)

### In

- decided-by [daemon-live-cpu-profile](../capability/daemon-live-cpu-profile.md)
- decided-by [daemon-passes-cpuprofile-to-drivers](../capability/daemon-passes-cpuprofile-to-drivers.md)
- decided-by [idle-poll-benchmarks](../capability/idle-poll-benchmarks.md)
- decided-by [root-daemon-idle-poll-no-work](../capability/root-daemon-idle-poll-no-work.md)
- decided-by [usage-scan-once-per-poll](../capability/usage-scan-once-per-poll.md)
- decided-by [ci-wait-backs-off](../capability/ci-wait-backs-off.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
