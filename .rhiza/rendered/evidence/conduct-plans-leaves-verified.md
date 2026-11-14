# conduct-plans-leaves-verified

evidence · `01M7BK8C42MF4B33TVQJ2H77K6` · actor [claude-code-202a10f7](../actor/claude-code-202a10f7.md)

## result

PASS

## statement

At d004463c7 re-ran every Checks: command and the acceptance commands, all exit 0: go test -trimpath ./internal/api -run 'TestPickNoLeaves|TestDaemonPicks|TestConductorWorkflow' (TestPickNoLeaves: conducted leafless milestone change set is a candidate, classic and backlog leafless ones skip no-leaves, the dependent umbrella-style one skips not-ready); go test -trimpath ./cmd/rhiza -run TestScripts/run-ship-needs-leaves (command-step ship-change refuses a leafless change set with 'has no leaves' and next 'rhiza run start --workflow plan'; a ship-change whose first step is an agent step declaring checks reaches 'would start' on the same leafless change set, fake rex/claude only); go test -trimpath ./cmd/rhiza -run 'TestScripts/(daemon-picks|run-ship)'; rhiza ship generate (0 changed); rhiza render (only this verifier's own actor rendered). Read the diff: conductorWorkflow/changeConducts in daemon.go, the Conducts hook in daemonCandidates wired from daemon_fold.go, and run.go gating checkShipLeaves on !conductorWorkflow; backlog short-circuits before Conducts, not-ready is checked first; no stubs or TODOs, nothing outside the leaf beyond the required daemon_fold.go wiring and the commit's log/rendered files.

## Links

### Out

- verifies [conduct-plans-leaves](../capability/rhiza/orchestration/daemon/conduct-plans-leaves.md)

## History

- `01M79KGZ51V2F5YYWMMQTGR6NW` A rhiza-conduct skill conducts a change set in one session: plans its leaves if it has none, spawns a builder and verifier per leaf, three review checks with a critic only on blocker or major, fixes, lands
