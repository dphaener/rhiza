# scripts-run-on-scaled-clock-verified-2

evidence · `01M7EX3D0ARF958C0FTQBVWJKM` · actor [claude-code-8e0825c2](../actor/claude-code-8e0825c2.md)

## result

PASS

## statement

Re-verified at 27088b322 (attempt 2). Split check: run-ship-conflict-fix-more.txtar's Case 4-7 body is byte-identical to lines 143-236 of run-ship-conflict-fix.txtar at c47bf9635, and the archive sections (setup.sh, review.toml, ship.sh, bin/rex, bin/claude, bin/gh) in both new scripts are byte-identical to the old one; run-ship-conflict-fix keeps cases 1-3 unchanged, only the header comment for cases 4-7 was replaced by a pointer. Ran (RHIZA_*/SHIP_*/REX_* unset, under flock machine-tests.lock and systemd-run MemoryMax=6G TasksMax=300) go test -p 2 -trimpath -v ./cmd/rhiza -run 'TestScripts/(run-ship-conflict-fix|run-ship-conflict-fix-more|run-ship-ci-fix|run-ship-main-dispatch-fail|arbiter-agent-tabs|plan-workflow-ci-infra-wait|session-message-busy|wait-timeout)$' at load avg ~2-3: exit 0, all PASS; run-ship-conflict-fix 27.79 s, run-ship-conflict-fix-more 29.62 s, run-ship-ci-fix 27.95 s, run-ship-main-dispatch-fail 9.42 s, plan-workflow-ci-infra-wait 6.91 s, arbiter-agent-tabs 4.57 s, session-message-busy 1.04 s, wait-timeout 4.59 s (RealClock failure case passes). Every named script is under 30 s, though run-ship-conflict-fix-more and run-ship-ci-fix run close to it when they share the CPU in one go test. Commit Checks rhiza ship generate and rhiza render: exit 0.

## Links

### Out

- verifies [scripts-run-on-scaled-clock](../capability/rhiza/integrity/self-check/scripts-run-on-scaled-clock.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).
