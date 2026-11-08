# question-db26d06b-1

question · `01M6XEGGR71R8C0XZETE1S3S60`

## session

db26d06b-42cb-4546-92e4-01a1bae405c7

## statement

Change set 01M6XE6 (Fix Main is never tested) blames a merge interaction: #374/#375/#376 each green, then #379 red on internal/api schedule and health tests (CI run 42473839029). The evidence says otherwise. Attempt 1's log shows health part disk 'bad: 3.7 GiB free, under the minimum of 10.0 GiB', so the daemon held schedule starts. Attempt 2 of the same run passed at the same commit, and so did #379's next head. Every merged head already contained origin/main (ship merge refuses otherwise, and ship check --main dispatches CI when the trees differ), so main's tree was the tested tree. The real defect is that internal/api Go tests read the host's free disk. cmd/rhiza scripts pin RHIZA_TEST_FREE_BYTES; the Go tests use a Ctx with no Getenv and fall through to daemon.FreeBytes. What should this change set build?

## Links

### In

- answers [main-red-was-host-disk-not-merge](../decision/main-red-was-host-disk-not-merge.md)

## History

- `01M6XE6M6DXRA2DTC19CTHYYK9` Fix Main is never tested, so two PRs that are each green can break it together: after scheduled triggers (#374), daemon watchdog (#375) and autonomy per workflow (#376) merged on 2026-11-08, PR #379, which changes only .rhiza/, failed internal/api tests TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow and TestHealthSchedulesPart (CI run 42473839029), while main's last CI run is from 2026-11-06 (CI no longer runs on push since the 'stop re-testing on main' change). Every open PR then goes red until someone notices. Either test the merge result before merging (require the PR branch to be up to date with main and re-run CI after main moves, or use a merge queue), or run CI on main after each merge and have the daemon pause new merges and file a bug when main goes red.
