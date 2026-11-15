# script-timings-file-verified

evidence · `01M7FKFJGFC86WDPF6ZY9BMTX2` · actor [claude-code-024cf479](../actor/claude-code-024cf479.md)

## result

PASS

## statement

At 960376888 (env -u RHIZA_*/SHIP_*/REX_*, capped scope behind the machine lock): go test -trimpath -p 2 ./tools/scripttimings exit 0; with -count=1 -v all 7 tests pass (TestMergeFromJSON, TestMergeFromRunLog over a four-job gh run view --log fixture, TestMergeNoPreviousFile, TestMergeFailuresLeaveFileUnchanged covering no passing subtest, no script-time line, three malformed previous files and bad flags with the file byte-identical and no temp left, TestParseScriptTimeLines, TestFormatScriptTimeRoundTrips, TestCheckedInFileParses); go vet exit 0; make script-timings JSON=tools/scripttimings/testdata/scripts.json exit 0 with the TSV unchanged; JSON=testdata/none.json exit 2 (tool exit 1) with the TSV unchanged; make script-timings with no source exit 2; rhiza ship generate and rhiza render exit 0. The diff shows failed or unrun scripts keep their previous time, removed scripts are dropped, and output is sorted name<TAB>seconds lines. cmd/rhiza/testdata/script-timings.tsv has 398 lines, the same as the 398 cmd/rhiza/testdata/script/*.txtar, and is byte-identical to the C-sorted script names, each at 1 second. The commit touches only tools/scripttimings/, Makefile, the TSV and its own .rhiza log/rendered files. go.mod and go.sum are unchanged (std library only), and there are no TODOs or stubs. The CI refresh (make script-timings RUN=...) is the orchestrator's step before merge, so it was not run here.

## Links

### Out

- verifies [script-timings-file](../capability/rhiza/integrity/self-check/script-timings-file.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
