# script-timings-measured

capability · `01M7FT147E37B6A5KCR328M971`

## acceptance

make script-timings RUN=43721082089, then go test -trimpath -p 2 -count=1 -json ./cmd/rhiza -run 'TestScripts/^(bugfix-workflow-autonomous|plan-workflow-autonomous-left|plan-workflow-autonomous|plan-review-workflow)$' > $TMPDIR/four.json and make script-timings JSON=$TMPDIR/four.json, leave a file with one line per cmd/rhiza/testdata/script/*.txtar (wc -l equals the ls count), the four scripts at their local times (each under 15), and the other scripts at run 43721082089's times (run-ship-conflict-fix-more among the top ten). go test -trimpath -p 2 ./tools/scripttimings -run '^TestCheckedInFileParses$' and go test -trimpath -p 2 ./cmd/rhiza -run '^(TestScriptShard|TestScriptOrder)$' pass, TestScriptShard's check that each shard's total is within 15% of the mean holding over the new file. Measure (orchestrator, on the next code PR): gh run view <run> --json jobs shows the four scripts jobs within 30 s of each other (2026-11-15, run 43721082089: 84-156 s). Touches cmd/rhiza/testdata/script-timings.tsv only.

## statement

cmd/rhiza/testdata/script-timings.tsv holds measured seconds for every script instead of the 1-second seed, from CI run 43721082089's script-time lines and, for the four scripts fake-rex-no-dead-wait sped up, a local run of those four, so CI's four scripts shards split by real time.

## Links

### Out

- depends-on [fake-rex-no-dead-wait](fake-rex-no-dead-wait.md)
- decided-by [script-timings-from-ci-then-local](../decision/script-timings-from-ci-then-local.md)

### In

- verifies [script-timings-measured-verified](../evidence/script-timings-measured-verified.md)
- verifies [script-timings-measured-qa-1](../evidence/script-timings-measured-qa-1.md)

## History

- `01M7FRV11EM44XJAJS6VXWC5DK` Finish the fast suite: the slow workflow scripts and internal/api off the critical path, and the parallel-scripts git race, so cold make test and median CI meet milestone 38's targets
