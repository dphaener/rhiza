# conduct-plans-leaves-verified-2

evidence · `01M7BV18KRZ3Y3F6MT7RY501VV` · actor [claude-code-d39bd10b](../actor/claude-code-d39bd10b.md)

## result

PASS

## statement

At bdc764ab3 (attempt 2): go test -trimpath ./internal/api -run 'TestPickNoLeaves|TestDaemonPicks|TestConductorWorkflow' ok; -run 'TestDaemon|TestPick|TestConductor|TestReadyWorkflow|TestStartWorkflow|TestUmbrella' ok; go vet ./internal/api clean; go test -trimpath ./cmd/rhiza -run TestScripts/run-ship-needs-leaves ok; -run 'TestScripts/(daemon-picks|run-ship)' ok. TestPickNoLeaves dispatches a leafless milestone change set through the real [ready] rule and changeConducts: picked for an agent-with-checks workflow, no-leaves skip for a leaves-first one, backlog still no-leaves, umbrella not-ready. The script refuses a leafless change set under a command-step ship-change with 'has no leaves' and next 'rhiza run start --workflow plan', then runs it to success under a conductor ship-change with fake bin/claude and bin/rex. Diff: conductorWorkflow/changeConducts in daemon.go, Conducts wired in daemon_fold.go, run.go guard; no TODOs.

## Links

### Out

- verifies [conduct-plans-leaves](../capability/rhiza/orchestration/daemon/conduct-plans-leaves.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
