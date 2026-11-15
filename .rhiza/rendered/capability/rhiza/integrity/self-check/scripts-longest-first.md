# scripts-longest-first

capability · `01M7EQE67BJPY5ER3PJFDYRKKM` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./cmd/rhiza -run '^TestScriptOrder$' shows: given timings, scripts come in descending seconds, ties and scripts missing from the file by name, and every *.txtar exactly once; failure case: a malformed timings line (no tab, or seconds not a number) fails with the file's line number. Touches cmd/rhiza/main_test.go and a new cmd/rhiza/order_test.go. Must not change any script's content or environment.

## statement

TestScripts hands testscript its Files longest first by script-timings.tsv, scripts missing from the file after them by name. This orders the list only, not the start: testscript runs each script as a t.Parallel subtest, and Go hands the -parallel slots to the woken subtests in no guaranteed order, so nothing here promises that the longest scripts start first.

## Links

### Out

- depends-on [script-timings-file](script-timings-file.md)
- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)

### In

- depends-on [script-shards-weighted](script-shards-weighted.md)
- blocks [review-01m7eyr9y4fr1bec2wqxx5sqw5-r3-f1](../../../../question/review-01m7eyr9y4fr1bec2wqxx5sqw5-r3-f1.md)
- verifies [scripts-longest-first-verified](../../../../evidence/scripts-longest-first-verified.md)
- verifies [scripts-longest-first-qa-1](../../../../evidence/scripts-longest-first-qa-1.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
