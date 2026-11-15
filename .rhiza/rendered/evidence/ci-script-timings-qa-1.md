# ci-script-timings-qa-1

evidence · `01M7FPJNNK7CB0WDV4AV4WN73H` · actor [claude-code-a30f328a](../actor/claude-code-a30f328a.md)

## result

PASS

## statement

QA at 6e4fed7ee with a scratch build of tools/scripttimings: replay on the testdata fixtures prints the failing script's text and --- FAIL:/FAIL lines plus script-time lines, exit 0 for replay-pass and 1 for replay-fail, replay-build, replay-killed and none; hand-made streams: a skipped subtest's text dropped and no script-time for it, a stream cut mid-line, empty input and non-JSON input all exit 1 with the no-package-end message; a real go test -json of TestScripts/version and describe-golden replays to ok plus two script-time lines, exit 0. make test-scripts SHARD=2/3 with a fake GO_TEST printing fixtures exits 0 for pass/0 and fails with go test's status kept (Error 1, 3, 137) for fail/1, pass/3, killed/0, build/1, killed/137; make -n shows RHIZA_TEST_SHARD='2/3', -json, the replay and the lock.

## Links

### Out

- verifies [ci-script-timings](../capability/rhiza/integrity/self-check/ci-script-timings.md)

## History

- `01M7EQ3K7PW2XYR684CXSX74TG` Parallel scripts, longest first, weighted CI shards: cmd/rhiza's scripts run at a -parallel bounded by cores and memory, start longest first, and CI's script shards split by a checked-in timing file instead of round-robin by name. Evidence: scripts use ~1.2 cores at -parallel 4; at -parallel 12 they take 381 s instead of 639 s with 0 failures; CI shards took 61-239 s.
