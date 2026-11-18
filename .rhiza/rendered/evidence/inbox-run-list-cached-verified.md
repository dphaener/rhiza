# inbox-run-list-cached-verified

evidence · `01M7PM55G675AKGHWR02RNA9N7` · actor [claude-code-35b2a147](../actor/claude-code-35b2a147.md)

## result

PASS

## statement

Verified at c292adb2f. go test -trimpath -count=1 -p 2 ./internal/run: ok (TestStoreListDiskCacheParsesOnlyChanged, TestStoreListDiskCacheIgnoresBadFile corrupt/garbage/other version, TestStoreListNoDiskCacheWritesNothing pass; with diskStore's DiskCache flipped off in a scratch worktree the first fails 'parsed 3 journals, want only the appended one', so it needs the on-disk cache). Tests set XDG_CACHE_HOME and HOME to TempDirs; ~/.cache/rhiza/run-list untouched by the test runs. go test -trimpath -count=1 -p 2 ./internal/api -run Inbox: ok 12.5s. cmd/rhiza TestScripts/(inbox-.*|tui-inbox-app)$ with RHIZA_/SHIP_/REX_ unset: ok. Profile pairs in ~/code/rhiza, base a7e6ab696 vs leaf, warm call first, back to back at load ~4.6: run.(*Store).List 27.48% (360ms) -> 13.04% (90ms) and 37.84% (560ms) -> 10.29% (70ms); /usr/bin/time warm call 1.14 s -> 0.72 s and 1.35 s -> 0.62 s. rhiza inbox --json base vs leaf byte-identical. Diff touches only store.go, store_test.go and inbox.go's two rs.List call sites (plus the leaf's .rhiza records); no TODOs, journal format and List output unchanged.

## Links

### Out

- verifies [inbox-run-list-cached](../capability/rhiza/orchestration/inbox/inbox-run-list-cached.md)

## History

- `01M7P1Q9GKHXA37YBG2M8JP661` Milestone 39's test-time targets missed after all five change sets: cold make test 400 s wall, 963 CPU-s (630 user + 332 sys) at 5d2596d on a quiet box (target 240 s). internal/api went from 129 s to 318 s once its tests run on the index (api-tests-on-index, #599): each test builds a SQLite index. cmd/rhiza is 254 s (was 222 s) despite in-process rhiza (#605); internal/index 113 s (was 48 s). Also rhiza inbox 645 ms (target 200 ms, was 531 ms before #606) and rhiza check 8.4 s (target 5 s, was 7.1 s). Profile each before changing anything (go test -cpuprofile per package; one built binary for inbox and check), then fix: shared prebuilt index fixtures or a cheaper index setup for api tests, what still costs in the scripts, and inbox's and check's hot paths. Never lower a target. Raw: /home/rhiza/tmp/claude-1000/perf-data/m39/cold.log
