# whatif-copies-only-written

capability · `01M7P20XK14P9JE95J3GQVS36H` · parent [check](../check.md)

## acceptance

go test -trimpath -count=1 ./internal/model/... ./internal/check/... ./internal/index/... passes, and go test -trimpath -count=1 -p 2 ./internal/api -run 'Parity|Projection|WhatIf' passes. A new benchmark or test in internal/model shows WhatIfFrom's allocations on a generated model of a few thousand nodes with a one-op change set fall by at least 80% from before (testing.AllocsPerRun or -benchmem, both numbers in the Checks: block). Checks: block gives RHIZA_CPUPROFILE of $TMPDIR/rhiza-x check in the main checkout before and after, showing model.WhatIfFrom's share fall, and /usr/bin/time -f '%e s' both times. Failure case: a test in internal/model builds two what-if views from one seed, each with ops updating the same node's fields, history and aspects, and asserts the seed and the other view are unchanged

## files

internal/model/seed.go (Seed.builder), internal/model/apply.go (builder.node and builder.link, and the History, aspects and Retired writes the copy-on-write must intercept), internal/model/view.go (whatIf, whatIfLayer, layer), internal/model/*_test.go for the new test and benchmark

## must_not

Change any view's nodes, links, problems or versions; change Seed's exported shape or the index's stored format; touch internal/check

## statement

model.WhatIfFrom (internal/model) copies from the seed only the nodes and links the view's ops write, sharing the rest read-only, so building a what-if view costs the change set's ops rather than the whole model

## Links

### Out

- decided-by [m39-whatif-copy-on-write](../../../../decision/m39-whatif-copy-on-write.md)
- decided-by [m39-time-profile-first](../../../../decision/m39-time-profile-first.md)

### In

- depends-on [check-whatif-finish-cheaper](check-whatif-finish-cheaper.md)
- verifies [whatif-copies-only-written-verified](../../../../evidence/whatif-copies-only-written-verified.md)
- verifies [whatif-copies-only-written-qa-1](../../../../evidence/whatif-copies-only-written-qa-1.md)
- verifies [whatif-copies-only-written-qa-2](../../../../evidence/whatif-copies-only-written-qa-2.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
