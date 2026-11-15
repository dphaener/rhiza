# api-waits-on-run-clock-qa-1

evidence · `01M7EXVMQKAZM3S0NDF2NN79YT` · actor [claude-code-e5658f1e](../actor/claude-code-e5658f1e.md)

## result

PASS

## statement

QA at 27088b322 with the scratch build: rhiza ship deploy against a refused deploy URL (rhiza.ship.deployTimeout 3s, RHIZA_DEPLOY_INTERVAL=1s) reported 'after 3 polls over 3s' in 2.12 s wall with the scale unset and 2.11 s with abc (timing unchanged), and the same 3 polls in 0.21 s at scale 1000 (deadline on the clock it sleeps on). go test -count=1 -p 2 ./internal/api -run 'TestRunClock|TestTypePointer|TestShipDeploy|TestSpawnSleep|TestRetry|TestSessionMessage|TestRun' ok in 10.9 s.

## Links

### Out

- verifies [api-waits-on-run-clock](../capability/rhiza/integrity/self-check/api-waits-on-run-clock.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).
