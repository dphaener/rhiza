# m39-api-index-per-test-copy

decision · `01M7P20X15CAXSAPT33THYGRC2`

## reasoning

The intent offers shared prebuilt index fixtures or a cheaper index setup. Sharing one open index between parallel tests is ruled out by api-tests-on-index's must_not; a per-test copy of a prebuilt file gives the prebuilt fixture without sharing. The index's freshness check (folder stat keys, then entry names) already handles a copied tree whose inodes and change times differ, so the copy needs no production change

## statement

internal/api's tests keep one index file per test, as api-tests-on-index requires, but stop building it from nothing: a snapshot template is indexed once when it is built and each copy gets a copy of that index file beside its tree, so a test catches up only what it writes after the copy. Tests on a fresh root keep a fresh index; their setup cost is cut only where the profile shows it

## Links

### Out

- derived-from [m39-test-time-research](../source/m39-test-time-research.md)

### In

- decided-by [api-tests-copy-prebuilt-index](../capability/rhiza/integrity/self-check/api-tests-copy-prebuilt-index.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
