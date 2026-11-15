# script-timings-measured-verified

evidence · `01M7FYHD7FBMS6ZQ397HWZD7Q3` · actor [claude-code-177bd8d1](../actor/claude-code-177bd8d1.md)

## result

PASS

## statement

At 0d73a0435: cmd/rhiza/testdata/script-timings.tsv has 398 lines, the same as the 398 cmd/rhiza/testdata/script/*.txtar files, with the names matching one to one. The four scripts are at local times: bugfix-workflow-autonomous 2.55, plan-review-workflow 2.13, plan-workflow-autonomous 2.68, plan-workflow-autonomous-left 4.19 (all under 15). run-ship-conflict-fix-more 15.75 is 4th of the top ten. I re-ran make script-timings RUN=43721082089 (exit 0, 398 fresh times); it changed only the four scripts' lines, back to CI's 48-88 s, so every other script holds run 43721082089's time. The scoped go test of the four scripts with -json (exit 0, locally 1.77-4.0 s) then make script-timings JSON=four.json (exit 0, 4 fresh times) reproduced the committed values within timing noise. go test -trimpath -p 2 ./tools/scripttimings -run '^TestCheckedInFileParses$' passed. go test -trimpath -p 2 -count=1 ./cmd/rhiza -run '^(TestScriptShard|TestScriptOrder)$' passed on the committed file, including 'checked-in timings, four shards balanced' (each shard total within 0.15 of the mean, shard_test.go:152). The commit's only code-tree change is the tsv; the rest is its actor record and log entries. Not checked here: the CI jobs-within-30 s measure, which is the orchestrator's on the next PR. Note: the commit message repeats its Checks block and trailers twice (cosmetic).

## Links

### Out

- verifies [script-timings-measured](../capability/script-timings-measured.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
