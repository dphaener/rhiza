# api-tests-free-disk-hermetic-verified

evidence · `01M6XSAHRT14Q0AMPA422KPBCP` · actor [claude-code-3295c0e4](../actor/claude-code-3295c0e4.md)

## result

PASS

## statement

Verifier claude-code-3295c0e4 at 6b7d571f (attempt 2). Read the diff: freebytes_test.go adds pinFreeBytes (RHIZA_TEST_FREE_BYTES=1099511627776 beside RHIZA_TEST_NOW) and setFreeBytes; TestMain swaps diskFreeBytes (daemon.go, was daemon.FreeBytes) for a panic; the schedule harness drops MinFree "0" and schedHealthCtx drops the minFree=0 stub, both pinning instead; attempt 2 pins daemonRepo, web_root and usage_learn ctxs, and the two RHIZA_TEST_NOW-free web listener tests turn minFree/trimGoCache off. Ran: go test -trimpath ./internal/api -run 'TestDaemonPollStartsEachProjectsSchedule|TestDaemonScheduleUnresolvedWorkflow|TestHealthSchedulesPart|TestFreeBytesOverride' -count=1 -> ok, all four PASS. TestFreeBytesOverride sets 5 GiB and sees daemonDiskLow low and health's disk part HealthBad. go test -trimpath ./cmd/rhiza -run 'TestScripts/(health|daemon-notify-disk-low)' -count=1 -> ok. Mutation check in a scratch worktree with the RHIZA_TEST_FREE_BYTES case removed from pinFreeBytes: each of the three schedule/health tests, run alone, fails with the panic 'test read the host's real free disk space'. Attempt 1's finding f1 is fixed: go test -trimpath ./internal/api -count=1 -timeout 30m -v gives 531 PASS lines, no --- FAIL and no free-disk panic, ending in PASS. The package then reports FAIL only from rootguard: ~/code/.rhiza-root/usage.live.json changed, a file the live rhiza-root.service (active on this VM) writes, unrelated to this diff. An earlier run under the default 600s timeout timed out under load 8 on 4 cores; TestLinkEndpointKinds alone took 119s.

## Links

### Out

- verifies [api-tests-free-disk-hermetic](../capability/rhiza/integrity/self-check/api-tests-free-disk-hermetic.md)

## History

- `01M6XE6M6DXRA2DTC19CTHYYK9` Fix Main is never tested, so two PRs that are each green can break it together: after scheduled triggers (#374), daemon watchdog (#375) and autonomy per workflow (#376) merged on 2026-11-08, PR #379, which changes only .rhiza/, failed internal/api tests TestDaemonPollStartsEachProjectsSchedule, TestDaemonScheduleUnresolvedWorkflow and TestHealthSchedulesPart (CI run 42473839029), while main's last CI run is from 2026-11-06 (CI no longer runs on push since the 'stop re-testing on main' change). Every open PR then goes red until someone notices. Either test the merge result before merging (require the PR branch to be up to date with main and re-run CI after main moves, or use a merge queue), or run CI on main after each merge and have the daemon pause new merges and file a bug when main goes red.
