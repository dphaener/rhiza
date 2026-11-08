# api-tests-free-disk-hermetic

capability · `01M6XENB2AST9TAX4MTWCDY0VK` · parent [self-check](../self-check.md)

## acceptance

internal/api tests reach repoFreeBytes only through a Ctx (or one shared test helper building it) whose Getenv pins RHIZA_TEST_FREE_BYTES=1099511627776 beside RHIZA_TEST_NOW, and the package's test setup makes consulting the real daemon.FreeBytes fail the test; go test -trimpath ./internal/api -run 'TestDaemonPollStartsEachProjectsSchedule|TestDaemonScheduleUnresolvedWorkflow|TestHealthSchedulesPart|TestFreeBytesOverride' -count=1 passes, where TestFreeBytesOverride is a new internal/api test that overrides the helper's figure with one under 10 GiB and sees daemonDiskLow report low (so a due schedule is held) and health's disk part go bad, proving a test can still set its own figure; go test -trimpath ./cmd/rhiza -run 'TestScripts/(health|daemon-notify-disk-low)' passes unchanged; with the pinned figure removed from the helper, the three schedule and health tests fail on the real-disk guard rather than passing on the host's disk

## statement

Every internal/api Go test sees a fixed 1 TiB of free disk unless it sets its own figure, so the daemon's minFree guard and health's disk part never read the CI host's real free space, and a runner short on disk no longer turns PRs red

## Links

### Out

- decided-by [main-red-was-host-disk-not-merge](../../../../decision/main-red-was-host-disk-not-merge.md)
- decided-by [test-setup-default-free-bytes](../../../../decision/test-setup-default-free-bytes.md)

### In

- verifies [api-tests-free-disk-hermetic-verified](../../../../evidence/api-tests-free-disk-hermetic-verified.md)
- verifies [api-tests-free-disk-hermetic-qa-1](../../../../evidence/api-tests-free-disk-hermetic-qa-1.md)
- verifies [pr-390-fix-main-is-never-tested](../../../../evidence/pr-390-fix-main-is-never-tested.md)

## History

- `01M6XE6M6DXRA2DTC19CTHYYK9` Fix Main is never tested, so two PRs that are each green can break it together: after scheduled triggers (#374), daemon watchdog (#375) and autonomy per workflow (#376) merged on 2026-11-08, PR #379, which changes only .rhiza/, failed internal/api tests TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow and TestHealthSchedulesPart (CI run 42473839029), while main's last CI run is from 2026-11-06 (CI no longer runs on push since the 'stop re-testing on main' change). Every open PR then goes red until someone notices. Either test the merge result before merging (require the PR branch to be up to date with main and re-run CI after main moves, or use a merge queue), or run CI on main after each merge and have the daemon pause new merges and file a bug when main goes red.
