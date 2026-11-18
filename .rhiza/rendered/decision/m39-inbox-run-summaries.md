# m39-inbox-run-summaries

decision · `01M7P20X9RTKTAMZ7RGZ39SMDK`

## reasoning

The inbox profile puts 37% in run.(*Store).List: readFile and parseFile of every journal, the in-process journalCache gone with each CLI process. Finished runs' journals never change, so a cache keyed on the file's stat lets one call parse only live ones, without changing the journal format older runs were written in

## statement

rhiza inbox stops parsing every run journal on each call: run.Store.List (or a listing the inbox uses in its place) reuses what it parsed of a journal whose file has not changed since, kept in the user cache keyed by path, size and modification time, and parses only new or changed journals. A changed journal is always read fresh

## Links

### Out

- derived-from [m39-test-time-research](../source/m39-test-time-research.md)

### In

- decided-by [inbox-run-list-cached](../capability/rhiza/orchestration/inbox/inbox-run-list-cached.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
