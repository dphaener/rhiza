# ci-script-timings

capability · `01M7EZ13K1QGZHVMCYCGB79TTD` · parent [self-check](../self-check.md)

## acceptance

go test -trimpath ./tools/scripttimings -run '^TestReplay' turns a go test -json fixture of passing, failing and skipped TestScripts subtests into the failing ones' text (with their --- FAIL: lines) and a script-time line per passing one, and a fixture whose package passed exits 0; failure cases: a fixture with a failed subtest, a build failure, or no package end event (a killed go test) exits non-zero. go test -trimpath ./cmd/rhiza -run '^TestMakeTestShardTargets$' shows make -n test-scripts SHARD=2/3 passing -json and the replay, with RHIZA_TEST_SHARD, the lock and go test's failing status kept. CI's scripts step stays make test-scripts SHARD=<shard>. Touches tools/scripttimings/, Makefile, cmd/rhiza/makeupdate_test.go. Must not change the CI matrix, add a workflow step or action, or drop the machine-wide test lock.

## statement

make test-scripts runs go test -json through tools/scripttimings's replay, which prints go test's text for every test that did not pass and for the package (its --- FAIL:, FAIL and panic: lines as go test prints them), then one script-time<TAB><name><TAB><seconds> line per passing TestScripts subtest, and exits non-zero whenever go test failed, so each CI scripts job's log carries its scripts' times.

## Links

### Out

- depends-on [script-timings-file](script-timings-file.md)
- decided-by [script-times-from-ci-log](../../../../decision/script-times-from-ci-log.md)

### In

- verifies [ci-script-timings-verified](../../../../evidence/ci-script-timings-verified.md)
- verifies [ci-script-timings-qa-1](../../../../evidence/ci-script-timings-qa-1.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
