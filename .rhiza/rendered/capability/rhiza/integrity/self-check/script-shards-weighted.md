# script-shards-weighted

capability · `01M7EQE69J35XJW25SPBZBH1MS` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./cmd/rhiza -run '^(TestScriptShard|TestScriptShardEnvError)$' shows the shards disjoint and covering every script, greedy by recorded time on a fixture of timings (longest first, each to the least-loaded shard, a missing script weighted at the median, ties to the lower shard), each shard's scripts in scripts-longest-first's order with missing scripts last by name, and with the checked-in timings each shard's total within 15% of the mean; failure case: RHIZA_TEST_SHARD=0/4, 5/4 and 1/0 are still refused. The test-scripts-sharded capability's statement says the split is weighted. Measure (orchestrator, after 5 PRs): gh run view <run> --json jobs shows the four scripts jobs within 60 s of each other (2026-11-15, run 43685514513 after the virtual clock, round-robin: 106-197 s). Touches cmd/rhiza/main_test.go, cmd/rhiza/shard_test.go. Must not change the shard env format or the CI matrix.

## statement

RHIZA_TEST_SHARD=i/n splits cmd/rhiza's scripts by recorded time instead of round-robin by name: each script, longest first by script-timings.tsv (a missing one weighted at the median for placement only), goes to the shard with the least total so far, so CI's four script shards take about the same time. Within a shard, scripts run in scripts-longest-first's order (timed scripts by descending seconds, then scripts missing from the file by name); the median weight never moves a missing script ahead of a timed one.

## Links

### Out

- depends-on [script-timings-file](script-timings-file.md)
- depends-on [scripts-longest-first](scripts-longest-first.md)
- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)

### In

- verifies [script-shards-weighted-verified](../../../../evidence/script-shards-weighted-verified.md)
- verifies [script-shards-weighted-qa-1](../../../../evidence/script-shards-weighted-qa-1.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
