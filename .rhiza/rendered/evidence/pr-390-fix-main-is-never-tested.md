# pr-390-fix-main-is-never-tested

evidence · `01M6XVBKDDT1KP4NCXF8HPC8NT` · actor [dphaener](../actor/dphaener.md)

## head

ccd6c36353d989bdc9b194aff8df15acffa1b156

## pr

https://github.com/dphaener/rhiza/pull/390

## result

10 checks passed: lint SUCCESS, scripts (1/4) SUCCESS, scripts (2/4) SUCCESS, scripts (3/4) SUCCESS, scripts (4/4) SUCCESS, api (1/2) SUCCESS, api (2/2) SUCCESS, rest SUCCESS, budgets SUCCESS, check SUCCESS

## statement

Pull request #390 (https://github.com/dphaener/rhiza/pull/390) passed CI at ccd6c36

## status

verified

## Links

### Out

- verifies [api-tests-free-disk-hermetic](../capability/rhiza/integrity/self-check/api-tests-free-disk-hermetic.md)

## History

- `01M6XE6M6DXRA2DTC19CTHYYK9` Fix Main is never tested, so two PRs that are each green can break it together: after scheduled triggers (#374), daemon watchdog (#375) and autonomy per workflow (#376) merged on 2026-11-08, PR #379, which changes only .rhiza/, failed internal/api tests TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow and TestHealthSchedulesPart (CI run 42473839029), while main's last CI run is from 2026-11-06 (CI no longer runs on push since the 'stop re-testing on main' change). Every open PR then goes red until someone notices. Either test the merge result before merging (require the PR branch to be up to date with main and re-run CI after main moves, or use a merge queue), or run CI on main after each merge and have the daemon pause new merges and file a bug when main goes red.
