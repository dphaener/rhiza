# scripts-run-on-scaled-clock-qa-1

evidence · `01M7EXVMS6B9403CTNWB7W1S4Z` · actor [claude-code-e5658f1e](../actor/claude-code-e5658f1e.md)

## result

PASS

## statement

QA at 27088b322: go test -count=1 -p 2 -v ./cmd/rhiza -run 'TestScripts/(run-ship-conflict-fix|run-ship-conflict-fix-more|run-ship-ci-fix|run-ship-main-dispatch-fail|arbiter-agent-tabs|plan-workflow-ci-infra-wait|session-message-busy|wait-timeout)$' all PASS together: conflict-fix 25.3 s, conflict-fix-more 29.3 s, ci-fix 25.8, main-dispatch-fail 7.0, arbiter-agent-tabs 3.5, plan-workflow-ci-infra-wait 7.0, session-message-busy 0.6, wait-timeout 4.4. The split's two scripts hold the same 94 distinct exec/stdout/stderr/exists lines as the old run-ship-conflict-fix. By hand with the built binary, the wait-timeout workflow gives 3 checks on RealClock and 40 at scale 1000.

## Links

### Out

- verifies [scripts-run-on-scaled-clock](../capability/rhiza/integrity/self-check/scripts-run-on-scaled-clock.md)

## History

- `01M7EQ2DHBZE64PPPZ36GQ599V` Scripts wait on a virtual clock: cmd/rhiza's testscripts run the engine and the api's run polls on a virtual clock, so the builtin workflows' 30 s polls, 2 min backoff and 2 s spawn polls cost no wall time in tests, while production keeps RealClock and its timing unchanged. Evidence: run-ship-ci-fix takes 165 s wall for 18 s CPU; scripts sleep about 90% of the time (~450 s of a cold make test).
