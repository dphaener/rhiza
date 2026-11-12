# root-daemon-idle-poll-no-work-verified-2

evidence · `01M7909VYBSDC5SHDRD979S0H8` · actor [claude-code-5ecd501e](../actor/claude-code-5ecd501e.md)

## result

PASS

## statement

Verifier, attempt 2, at 61d125ffa (fix commit: supersedes links root-daemon-idle-poll-no-work -> daemon-idle-poll-no-fold and decision-idle-cpu-scope -> daemon-poll-caches-fold-not-picks, plus render). Read the diff: two link.create supersedes ops, the builder's actor record and re-rendered files, no code change. rhiza show confirms both older nodes now show 'superseded by'. Re-ran the acceptance: go test -trimpath ./internal/api/... -run 'RootIdlePoll' -v -count=1 -> TestRootIdlePoll PASS (second unchanged poll folds nothing and computes no pick; each stamp change redoes the work), TestRootIdlePollBudget PASS (median 71 ms, own 54 ms CPU < 100 ms budget). make self-check in a clean detached worktree at the commit under env -i: exit 0. (Under the run's env it fails the same way at the parent commit, on the run actor's unrendered file and docs, so this commit does not cause it.)

## Links

### Out

- verifies [root-daemon-idle-poll-no-work](../capability/root-daemon-idle-poll-no-work.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
