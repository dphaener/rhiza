# ci-wait-backs-off-qa-1

evidence · `01M78ZSA4C44V3N9STNVA46QJM` · actor [claude-code-a1181376](../actor/claude-code-a1181376.md)

## result

PASS

## statement

QA at 4f1ad9a09: a real daemon-started driver ran a scratch workflow whose wait step (interval 2s, max_interval 8s) logged each check's time and printed a state file: with unchanged output the checks came 2,4,8,8,8,8 s apart, after the output changed 2 then 4 s, and the step passed when the condition held. rhiza workflow check rejects max_interval 10s below interval 30s and passes on every .rhiza/workflows/*.toml (bugfix, docs, docs-build, docs-review, plan-bug, ship-change); the three CI waits in ship-change.toml set 30s/2m. go test -run 'WaitBackoff|MaxInterval' on internal/run and internal/workflow passes.

## Links

### Out

- verifies [ci-wait-backs-off](../capability/ci-wait-backs-off.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
