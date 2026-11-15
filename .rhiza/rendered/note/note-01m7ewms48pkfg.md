# note-01m7ewms48pkfg

note · `01M7EWMS48PKFG084HD6HHJ474`

## statement

verifier finding (attempt 1, session 37a728f8): acceptance 'each of the first five scripts' elapsed under 30 s' FAILS for run-ship-conflict-fix: 33.40 s alone (go test -count=1 -p 2 -trimpath -v ./cmd/rhiza -run 'TestScripts/run-ship-conflict-fix$'; 38.5 s CPU in 34.9 s wall, load ~5), 51.84 s in the seven-script run. CPU-bound, so the clock cannot fix it: split its seven ship cases into two scripts each under 30 s, or amend the leaf's bound for this script. Everything else verified: all seven named scripts pass uncached (ci-fix 14.7 s alone, 31.4 s in the loaded group; main-dispatch-fail 11.3, arbiter-agent-tabs 8.6, plan-workflow-ci-infra-wait 8.9, session-message-busy 1.1, wait-timeout 4.2-4.8 with >=10 scaled checks and exactly 3 on RealClock); diff touches only main_test.go, arbiter-agent-tabs.txtar and wait-timeout.txtar; arbiter-agent-tabs' 'block inspect' -> exit 1 matches the old outcome (runProcess, internal/api/run.go:1981, calls it at run start before any tab, so the old fallthrough waited 40 s for $WORK/ready then exited 1) and its assertions are unchanged. rhiza run findings --leaf was refused: no leaves step in this conduct run.

## Links

### Out

- discusses [scripts-run-on-scaled-clock](../capability/rhiza/integrity/self-check/scripts-run-on-scaled-clock.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).
