# api-tests-parallel-runtime-verified

evidence · `01M7F61VZT1PWQYP443XG7S4DW` · actor [claude-code-4ec56230](../actor/claude-code-4ec56230.md)

## result

PASS

## statement

Verified at a8b6ed9ca. The -run pattern covers 208 Test funcs from the daemon_/run_/ship_/watchdog_/health_/usage_/feed_ files without newSessionFixture (none of the touched files call it). go test -p 2 -trimpath -count=1 -shuffle=on (77.8s) and -shuffle=424242 (76.9s): both ok. -race split into four sequential runs by prefix, each file's serial and parallel tests together: daemon_ 55s, run_ 26s, ship_ 15s, watchdog_|health_|usage_|feed_ 212s. All exit 0 with no DATA RACE. The single-command -race was not run: the parts total about 308s, over the 4-minute budget. go vet ./internal/api: clean. go test -list '.*': 756 Test + 12 Benchmark = 768, as the builder reported. The diff removes no lines and adds no funcs, so nothing was skipped, deleted or merged. rhiza ship generate: exit 0, 0 files changed. rhiza render: no diff beyond this verifier's own actor. Diff review: the internal/api changes are only t.Parallel() calls and comments in in-scope _test.go files. The other files are the builder's actor record and its rendering. No non-test code changed. 159 tests are parallel and 49 serial, and each serial test has a // Serial: comment. Each reason checks out: it swaps daemonReexec/daemonTicker/webFactory/pollGitTimeout/profileWaitEvery/sweepLeftoversWork/feedTicker/channelUpdate*/prePushCheck/runClock (useClock), sets the fake agents' resumeTranscript/holdFakeRecords/holdFakeLive/usageFakeReset/limitLaunches, or calls t.Setenv/t.Chdir, directly or through exposureRepo/newPrePushRepo/useEnvClock. Every assignment to a package var in these files is in a serial test. No parallel test calls a state-swapping helper (runDaemonTicks, newPrePushRepo, realPrePushCheck, useClock, limitFixture, exposureRepo, watchFixture and others), and none uses run temp folders.

## Links

### Out

- verifies [api-tests-parallel-runtime](../capability/rhiza/integrity/self-check/api-tests-parallel-runtime.md)

## History

- `01M7EQ3V3JMPYXS0WKE0DM8CSM` internal/api tests run in parallel: independent tests call t.Parallel and share snapshot fixtures instead of rebuilding the model per test; projection replacing FoldLog stays out of scope. Evidence: 749 tests, no t.Parallel in the repo, 58% of CPU in store.ReadLog because each write re-reads the log; ~180 s.
