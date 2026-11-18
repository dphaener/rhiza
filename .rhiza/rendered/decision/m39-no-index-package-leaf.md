# m39-no-index-package-leaf

decision · `01M7P20WX9YGJV4ZVJWHMBR4G5`

## reasoning

Run alone on 2026-11-18 (go test -trimpath -count=1 -p 2 -cpuprofile ./internal/index, the VM at load about 3.6), internal/index took 41 s wall and 42 CPU-s (34.9 user + 7.6 sys), below the 48 s it took before milestone 39; its top costs are catchUpFrom and refold under TestBudget's generated model, which the budget needs

## statement

01M7P1Q9 has no leaf for internal/index's tests: their 113 s in the cold make test is contention from the packages beside them (internal/api's 318 s above all), not their own cost, and the api and script leaves cut that contention

## Links

### Out

- derived-from [m39-test-time-research](../source/m39-test-time-research.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
