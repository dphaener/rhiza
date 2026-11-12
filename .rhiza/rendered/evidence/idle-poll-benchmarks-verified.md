# idle-poll-benchmarks-verified

evidence · `01M78XRZXFKSCYV1BF9BBKVX6S` · actor [claude-code-e8a86d83](../actor/claude-code-e8a86d83.md)

## result

PASS

## statement

At 770814896, ran go test -trimpath ./internal/api/ -run '^$' -bench 'IdleRootPoll|IdleDriverPoll' -benchtime 3x -cpuprofile $TMPDIR/idle.pprof: PASS in 75s (IdleRootPoll ~201ms/op, IdleDriverPoll ~56us/op) and wrote a 145KB profile that go tool pprof reads. Read the diff: the fixture is built in t.TempDir dirs, git and Detach are faked, Getenv is stubbed, no network/Rex/agent/gh; 4 runs in flight, 4,000 transcripts under a temp CLAUDE_CONFIG_DIR. In a throwaway scratch worktree, folded the grown fixture log: 1,099 change sets, 0 log problems, 0 fold problems. go test -run 'DaemonIdlePoll|DaemonPollCache|DaemonFold' ./internal/api passes with the refactored newDaemonFoldFix.

## Links

### Out

- verifies [idle-poll-benchmarks](../capability/idle-poll-benchmarks.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
