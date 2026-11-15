# scripts-longest-first-verified

evidence · `01M7FMP15H3P7QJTW8QM5CWHGJ` · actor [claude-code-63d24d17](../actor/claude-code-63d24d17.md)

## result

PASS

## statement

At 29dafac45: re-ran the Checks: block under the machine lock with RHIZA_/SHIP_/REX_ unset: go test -trimpath -p 2 ./cmd/rhiza -run '^(TestScriptOrder|TestScriptShard|TestScriptShardEnvError)$' exit 0 (all 9 TestScriptOrder subtests pass, also with -count=1); go test -trimpath -p 2 ./cmd/rhiza -run 'TestScripts/^(add|apply)$' exit 0; rhiza ship generate exit 0 (0 files changed); rhiza render exit 0 (no tracked change). Diff read: only cmd/rhiza/main_test.go (TestScripts now passes Files from scriptFiles, Dir dropped) and new cmd/rhiza/order_test.go plus the builder's model log entries and actor render. TestScriptOrder covers descending seconds, ties and untimed scripts by name, a shard ordered too, missing timings file by name, and every testdata/script/*.txtar exactly once against the real script-timings.tsv; malformed cases (no tab, seconds not a number, negative, empty name, duplicate) each assert the error contains '<file>:<line>:'. No .txt scripts exist, so the *.txtar glob matches testscript's Dir set. No script content, Setup or environment changed. Per Decision longest-first-is-list-order the code only sorts the list (sort.SliceStable), with no start-order mechanism; comments say so. No stubs or TODOs.

## Links

### Out

- verifies [scripts-longest-first](../capability/rhiza/integrity/self-check/scripts-longest-first.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
