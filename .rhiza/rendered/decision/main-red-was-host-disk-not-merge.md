# main-red-was-host-disk-not-merge

decision · `01M6XEKESD34G4P39FM1NFECZ3`

## reasoning

The question's evidence checks out in the code: repoFreeBytes (internal/api/daemon.go) honours RHIZA_TEST_FREE_BYTES only through ctxGetenv(ctx) alongside RHIZA_TEST_NOW and otherwise calls daemon.FreeBytes(ctx.Root); cmd/rhiza/main_test.go pins testFreeBytes = 1 TiB for scripts, but the internal/api tests build Ctx values without that Getenv, so they read the runner's real disk. A runner at 3.7 GiB free trips the 10 GiB guard, the daemon holds schedule starts and the health disk part goes bad, which is exactly the three failing tests; attempt 2 at the same commit passed, and #379's next head passed, which a real merge interaction would not do. ci.yml already dispatches CI for a main commit whose tree differs from the tested head, and ship merge requires the head to contain origin/main, so main's tree was the tested tree. Building a merge queue or post-merge red-main pause would add machinery against a failure that did not happen and leave the actual flake in place; fixing the test hermeticity removes the cause and stays inside the change set's purpose (main and PRs going red for no real defect).

## statement

This change set builds no merge gate, merge queue or post-merge CI on main. It makes the internal/api Go tests hermetic about free disk: every test that reaches repoFreeBytes (the daemon's disk guard and health's disk part, via the schedule and health tests) gets a pinned free-disk figure, the same 1 TiB the cmd/rhiza scripts pin with RHIZA_TEST_FREE_BYTES, through the Ctx's Getenv (or one shared test helper that builds such a Ctx), so no internal/api test reads the CI host's real free space. A test that wants low disk sets the figure itself. Acceptance: with the host's free space below 10 GiB (simulated by making daemon.FreeBytes unreachable from the tests, e.g. a test that fails if the real figure is consulted), TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow and TestHealthSchedulesPart still pass. The change set's intent should be narrowed to this flaky-test fix.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-db26d06b-1](../question/question-db26d06b-1.md)
- derived-from [main-never-tested-plan-interview](../source/main-never-tested-plan-interview.md)

### In

- decided-by [api-tests-free-disk-hermetic](../capability/rhiza/integrity/self-check/api-tests-free-disk-hermetic.md)

## History

- `01M6XE6M6DXRA2DTC19CTHYYK9` Fix Main is never tested, so two PRs that are each green can break it together: after scheduled triggers (#374), daemon watchdog (#375) and autonomy per workflow (#376) merged on 2026-11-08, PR #379, which changes only .rhiza/, failed internal/api tests TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow and TestHealthSchedulesPart (CI run 42473839029), while main's last CI run is from 2026-11-06 (CI no longer runs on push since the 'stop re-testing on main' change). Every open PR then goes red until someone notices. Either test the merge result before merging (require the PR branch to be up to date with main and re-run CI after main moves, or use a merge queue), or run CI on main after each merge and have the daemon pause new merges and file a bug when main goes red.
