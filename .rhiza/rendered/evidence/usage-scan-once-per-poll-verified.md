# usage-scan-once-per-poll-verified

evidence · `01M78Y7JXYV7YDQ0GRECM77G38` · actor [claude-code-4fdc4cef](../actor/claude-code-4fdc4cef.md)

## result

PASS

## statement

Verifier at 5fe205c8d: go test -trimpath -count=1 ./internal/api/... ./internal/adapter/claudecode/... -run UsageScan passes (TestUsageScanOncePerPoll); also -race -count=2, ./internal/adapter/... in full and ./internal/api -run Usage pass, go vet clean. The test uses a fake ctx.Now clock and 120 transcripts and counts claudecode.UsageScanStats walks and reads across polls that each call usageRecordsFor twice (two projects): first poll 1 walk, 120 reads; 10s later 0 walks, 0 reads; one file appended 1 read, 0 walks; a new file within the minute is not walked, and the poll a minute later does 1 walk and 1 read and picks it up; final records reflect.DeepEqual a full UsageRecords read. Diff read: usageRecordsFor passes ctx.Now() to the new adapter.UsageReaderAt; the claude-code adapter walks a root at most once per usageWalkEvery (1m), otherwise Lstats only its known files and rereads only size/mtime changes, dropping vanished ones.

## Links

### Out

- verifies [usage-scan-once-per-poll](../capability/usage-scan-once-per-poll.md)

## History

- `01M78TV0ZFXRBX1M08V6G7H19C` Fix The root daemon still burns CPU while idle: on 2026-11-12 the 30s sampler measured rhiza daemon at a steady ~0.45 core and the per-run drivers (rhiza run start --foreground) at ~0.25 core, together 44% of all measured CPU on the 4-core VM, all Go code, including idle stretches with no run doing work. Earlier bug 01M6G0KF (daemon pins a full core: each 5 s poll re-reads and re-folds the whole log; cache on stamps) is still planned, unbuilt; this is the same problem at a lower level after other fixes, and the drivers' polling adds to it. Wanted: profile the daemon and a driver at idle (pprof), make a poll with nothing changed cheap (skip fold/pick when log, refs, run records and settings are unchanged), slow driver polling while a step waits on CI, and add a CPU budget test for an idle poll. Evidence: sampler data /home/rhiza/tmp/claude-1000/-home-rhiza-code-rhiza/823cd4a5-db47-4770-9e31-1f8134228fec/scratchpad/perf/samples.csv; report ~/orchestrator/perf-2026-11-12/report-prospective.md.
