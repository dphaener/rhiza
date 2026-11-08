# main-never-tested-plan-interview

source · `01M6XEMX4KANVKMW49X4XJMNSW`

## quotes

```json
[]
```

## raw_sha256

bfb844d72dc14eb9ce87982b0589d732d24555017bafde4f067b1f3a3cc7b162

## session

db26d06b-42cb-4546-92e4-01a1bae405c7

## summary

Autonomous plan run of change set 01M6XE6 (no owner present). The intent said two green PRs (#374/#375/#376) broke main together and that #379 then failed the internal/api tests TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow and TestHealthSchedulesPart (CI run 42473839029). It proposed a merge gate, a merge queue, or CI on main with the daemon pausing merges. The planner checked the evidence. Attempt 1's api (1/2) job log shows health's disk part bad ('3.7 GiB free, under the minimum of 10.0 GiB'), so the daemon held schedule starts. Attempt 2 at the same commit passed, and so did #379's next head 0f9215bc. Every merged PR head already contained origin/main (ship merge refuses otherwise, and ship check --main dispatches CI when the merge commit's tree differs), so main's tree was the tested tree. Root cause: repoFreeBytes (internal/api/daemon.go) honours RHIZA_TEST_FREE_BYTES only through the Ctx's Getenv. cmd/rhiza scripts pin 1 TiB, but internal/api tests build Ctx values without it and read the host's real free disk. The planner raised the scope as a question. The arbiter ruled (Decision main-red-was-host-disk-not-merge): build no merge gating; make internal/api tests hermetic about free disk with a pinned 1 TiB figure, so that a test which reads the real figure fails; narrow the change set to this flaky-test fix. The planner wrote one leaf under self-check, next to scripts-free-disk-hermetic, and linked the change set to milestone-25, because PRs going red for no defect stall the daemon's autonomous shipping. The change set's intent text cannot be edited, so the narrowed scope lives in the Decision and the leaf.

## title

Fix Main is never tested: plan interview (autonomous)

## Links

### In

- derived-from [main-red-was-host-disk-not-merge](../decision/main-red-was-host-disk-not-merge.md)

## History

- `01M6XE6M6DXRA2DTC19CTHYYK9` Fix Main is never tested, so two PRs that are each green can break it together: after scheduled triggers (#374), daemon watchdog (#375) and autonomy per workflow (#376) merged on 2026-11-08, PR #379, which changes only .rhiza/, failed internal/api tests TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow and TestHealthSchedulesPart (CI run 42473839029), while main's last CI run is from 2026-11-06 (CI no longer runs on push since the 'stop re-testing on main' change). Every open PR then goes red until someone notices. Either test the merge result before merging (require the PR branch to be up to date with main and re-run CI after main moves, or use a merge queue), or run CI on main after each merge and have the daemon pause new merges and file a bug when main goes red.
