# root-daemon-round-robin-verified-2

evidence · `01M6SK7NWBZPPNQY5X6VZY2BFW` · actor [claude-code-954e6fe9](../actor/claude-code-954e6fe9.md)

## result

PASS

## statement

Verifier, attempt 2, at ccb69422. Ran go test -trimpath ./internal/daemon/... -run RoundRobin -v: TestRoundRobinAcrossProjects passes (maxRuns 3 over A(3), B(2), C(1) picks A1, B1, C1) and TestRoundRobinPerProjectMax passes (A max 1, maxRuns 4 picks A1, B1, C1, B2). Read the diff: runDaemon (internal/api/daemon.go) now reads every project's demand through prepare/daemonDemand, calls daemon.RoundRobin with the machine-wide running count and the home project's maxRuns, and passes each project only its share into daemonPoll1; rhiza.project.maxRuns is now ProjectPick.Max. This fixes review finding f1 (RoundRobin had no non-test caller). go test -trimpath ./cmd/rhiza -run TestScripts/daemon passes, including the new daemon-round-robin.txtar (max-runs 1 starts one run in total; max-runs 3 starts 2+1 alternating; project maxRuns 1 gives 1+2) and the updated daemon-milestone-request. go test ./internal/api -run 'Daemon|Project|Settings|Usage|Governor' passes. The other scripts that mention maxRuns or projects pass, except tui-autonomy-app, which fails the same way at the parent commit 11590484 (a settings source check that has nothing to do with this leaf). The full internal/api package went past go's 10-minute default timeout in TestLinkEndpointKinds, which is unrelated; CI runs it.

## Links

### Out

- verifies [root-daemon-round-robin](../capability/rhiza/orchestration/daemon/root-daemon-round-robin.md)

## History

- `01M6RN3RCGAZ9RNTWTZPKFM61M` Root daemon: one daemon on the root builds every registered project round-robin under one machine-wide run cap and usage governor, with machine settings in the root config and project settings in each repo, under one service
