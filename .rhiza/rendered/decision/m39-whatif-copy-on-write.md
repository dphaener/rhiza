# m39-whatif-copy-on-write

decision · `01M7P20XH0X429BY4P79EYNPV5`

## reasoning

The check profile puts 20% of rhiza check in Seed.builder's copyNode/copyFields, run for each of about 64 in-flight change sets over the whole model, while a change set's ops touch a handful of nodes. internal/model is the spec (strict-engine-is-spec), so the change must not alter any view, only how it is built

## statement

model.WhatIfFrom stops deep-copying every node and link of the as-is seed for each in-flight change set: Seed.builder shares the seed's records read-only and copies a node or link (its fields, history and aspects) only when an op of the view writes it, so the seed and other views built from it never change. The strict engine's results stay byte-for-byte the same, checked by the parity and model tests

## Links

### Out

- derived-from [m39-test-time-research](../source/m39-test-time-research.md)

### In

- decided-by [whatif-copies-only-written](../capability/rhiza/integrity/check/whatif-copies-only-written.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
