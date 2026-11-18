# check-log-decode-cheaper

capability · `01M7P2X9C68DXSQNHA0QFRKTRM` · parent [check](../check.md)

## acceptance

go test -trimpath -count=1 ./internal/store/... ./internal/entry/... passes, and make self-check passes. Checks: block gives RHIZA_CPUPROFILE of $TMPDIR/rhiza-x check, built at the parent commit and at the commit, in the main checkout, with the pprof -top -cum lines of store.readEntries, entry.Decode and canon.Canonicalize, store.readEntries' share falling, and each run's user+sys CPU from /usr/bin/time. Failure case: a new test in internal/store reads a log holding an entry file that is not canonical and one that fails to parse, and gets the same errors as before the cut; the before and after binaries' rhiza check --json in the main checkout print the same bytes and exit code; the existing store and entry tests keep passing unchanged

## files

internal/store/read.go, internal/store/*_test.go, internal/entry/entry.go, internal/entry/entry_test.go

## must_not

Change the log's file format; accept an entry file the current reader refuses; change which findings rhiza check reports, their severity, rule or text; lower the 5 s target or any budget test threshold

## statement

rhiza check's read of the log costs less: store.readEntries' decode of every entry file (readEntry, entry.Decode and canon.Canonicalize: 17% of a RHIZA_CPUPROFILE of rhiza check in the main checkout on 2026-11-18, canon.Canonicalize 10.5% of the whole) is cut in internal/store and internal/entry, every entry still decoded and checked as strictly as before

## Links

### Out

- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)

### In

- depends-on [check-whatif-finish-cheaper](check-whatif-finish-cheaper.md)
- verifies [check-log-decode-cheaper-verified](../../../../evidence/check-log-decode-cheaper-verified.md)
- verifies [check-log-decode-cheaper-qa-1](../../../../evidence/check-log-decode-cheaper-qa-1.md)
- verifies [check-log-decode-cheaper-qa-2](../../../../evidence/check-log-decode-cheaper-qa-2.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
