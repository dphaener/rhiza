# api-tests-free-disk-hermetic-qa-1

evidence · `01M6XSMGP7059RJ5E9YTYK1KRP` · actor [claude-code-5a427573](../actor/claude-code-5a427573.md)

## result

PASS

## statement

QA claude-code-5a427573 at 8aff029d. Built the binary (go build -trimpath) and ran it in a scratch repo: rhiza health reads the real disk through the new diskFreeBytes var ('39.9 GiB free, minimum 10.0 GiB', matching df's 40G); with RHIZA_TEST_NOW and RHIZA_TEST_FREE_BYTES=5368709120 the disk part is 'bad 5.0 GiB free, under the minimum'; RHIZA_TEST_FREE_BYTES without RHIZA_TEST_NOW is ignored (real 39.9 GiB). go test -trimpath ./internal/api -run for TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow, TestHealthSchedulesPart, TestFreeBytesOverride, TestDaemonWeb*, TestWebRootHandlerServesEveryProject and TestLearnedLimitRatchetsAtPoll: every test PASS (the package then reports the known rootguard failure on ~/code/.rhiza-root/usage.live.json, written by the live root service). go test -trimpath ./cmd/rhiza -run 'TestScripts/^(health|daemon-notify-disk-low)$' -> ok. Mutation: with the RHIZA_TEST_FREE_BYTES case renamed in pinFreeBytes, each of the three schedule/health tests run alone fails with the panic 'test read the host's real free disk space'; reverted.

## Links

### Out

- verifies [api-tests-free-disk-hermetic](../capability/rhiza/integrity/self-check/api-tests-free-disk-hermetic.md)

## History

- `01M6XE6M6DXRA2DTC19CTHYYK9` Fix Main is never tested, so two PRs that are each green can break it together: after scheduled triggers (#374), daemon watchdog (#375) and autonomy per workflow (#376) merged on 2026-11-08, PR #379, which changes only .rhiza/, failed internal/api tests TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow and TestHealthSchedulesPart (CI run 42473839029), while main's last CI run is from 2026-11-06 (CI no longer runs on push since the 'stop re-testing on main' change). Every open PR then goes red until someone notices. Either test the merge result before merging (require the PR branch to be up to date with main and re-run CI after main moves, or use a merge queue), or run CI on main after each merge and have the daemon pause new merges and file a bug when main goes red.
