# faster-reads-budget

evidence · `01M5V5BSX2ZYHR5C0NB8710Y32` · actor [claude-code-9a170ddb](../actor/claude-code-9a170ddb.md)

## result

TestReadBudget on an M3 Pro (11 CPUs), two runs per setting: default procs roadmap 4.93-5.29 ms and ready 4.92-5.11 ms (target 6 ms); GOMAXPROCS=2 roadmap 5.30-6.14 ms and ready 5.25-6.09 ms (target 10 ms); show 1.35-1.73 ms. make check exit 0.

## statement

roadmap and ready read the 50,000-node model well under the 50 ms budget: TestReadBudget medians of about 5 ms on an M3 Pro at default procs and 5.3-6.1 ms at GOMAXPROCS=2, from about 15 ms before the change set.

## tests

go test: TestReadBudget (internal/api), the freshness tests (internal/index), DecodeCanonical tests (internal/entry)

## Links

### Out

- verifies [fast-reads](../capability/rhiza/views/fast-reads.md)

## History

- `01M5SWTJY0AZP9TR3V1QJ54HMY` planGather (behind roadmap and ready) and the index's fresh read get faster, so CI has headroom under 25 ms again.
