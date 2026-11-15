# ci-script-timings-verified

evidence · `01M7FM5FC2GRCW2D5SDE7S71A0` · actor [claude-code-3d35ede6](../actor/claude-code-3d35ede6.md)

## result

PASS

## statement

Verified at ad0133144. Re-ran the Checks block: go test -p 2 -trimpath ./tools/scripttimings ok (TestReplayPass/Fail/BuildFailure/Killed/NoPackageEnd pass under -run '^TestReplay'); go test -trimpath ./cmd/rhiza -run '^TestMakeTestShardTargets$' ok; make -n test-scripts SHARD=2/3 shows go test -json with RHIZA_TEST_SHARD='2/3' piped to the replay, the tool built before lock "the cmd/rhiza scripts", unlock after, and go test's status kept through a file with the replay's as fallback. make test-scripts with GO_TEST printing fixtures exits 0 for replay-pass/0 and 2 for replay-fail/1, replay-pass/3 (Error 3, go test's status kept), replay-killed/0 and replay-build/1. go test -json -run '^TestScripts$/^version$' ./cmd/rhiza piped to replay prints script-time version 0.04, exit 0. rhiza ship generate and rhiza render in a clean detached worktree change nothing. The diff has no .github change, CI still runs make test-scripts SHARD=<shard>, it touches only tools/scripttimings, Makefile, cmd/rhiza/makeupdate_test.go and its own .rhiza records, and has no TODOs.

## Links

### Out

- verifies [ci-script-timings](../capability/rhiza/integrity/self-check/ci-script-timings.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
