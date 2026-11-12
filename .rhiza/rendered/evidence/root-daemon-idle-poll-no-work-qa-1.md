# root-daemon-idle-poll-no-work-qa-1

evidence · `01M78ZSA928XM8FDZRDDHN178C` · actor [claude-code-a1181376](../actor/claude-code-a1181376.md)

## result

PASS

## statement

QA at 4f1ad9a09: a real root daemon over a scratch root with two registered projects, nothing changing, against the real ~/.claude/projects (2,309 project folders, 4,093 transcripts) used 0.032 core over 30 s and 0.040 core over 75 s (/proc stat ticks), under the 0.05 core target; it still started a newly accepted change set's run when one appeared. go test -run RootIdlePoll on internal/api passes. Noted: the first poll after start reads every transcript, 94 s and 41 CPU-s cold, 16 s warm, tripping the 15 s watchdog note.

## Links

### Out

- verifies [root-daemon-idle-poll-no-work](../capability/root-daemon-idle-poll-no-work.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
