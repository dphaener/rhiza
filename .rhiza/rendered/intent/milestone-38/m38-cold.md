# m38-cold

intent · `01M7EQ74ZYGCCHY4QW22JZJ5SH` · parent [milestone-38](../milestone-38.md)

## statement

A cold make test on the 4-core VM, with nothing else running, finishes within 4 minutes of wall time (2026-11-15: 854 s). Measure: go clean -testcache && /usr/bin/time -f '%e s' make test.

## Links

### Out

- decided-by [fast-suite-proposal-approved](../../decision/fast-suite-proposal-approved.md)
- decided-by [fast-suite-measuring-commands](../../decision/fast-suite-measuring-commands.md)
- decided-by [m38-closed-cold-interim](../../decision/m38-closed-cold-interim.md)

### In

- serves change set `01M7EQ2DHBZE64PPPZ36GQ599V` (Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).)
- serves change set `01M7EQ3K7PW2XYR684CXSX74TG` (Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.)
- serves change set `01M7EQ3V3JMPYXS0WKE0DM8CSM` (internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.)
- serves change set `01M7EQ42Y1M8ENZCQ202D761JM` (Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).)
- relates-to [m38-cold-at-a98fd40](../../evidence/m38-cold-at-a98fd40.md) — Evidence that m38-cold is not met at a98fd4051 (interim result)

## History

- `01M7EQ19FNS1VDXT4750YX89RD` The test suite is fast. Owner, 2026-11-14/15, after the orchestrator measured the suite (cold make test 854 s wall for 972 CPU-s on 4 cores, warm 254 s): a cold make test within 4 minutes and a warm one within 60 s on the 4-core VM, the median PR CI run within 3 minutes, with no test skipped and no coverage lost; five change sets, the gate bug first, projection out of scope.
