# inbox-landed-plan-run-in-flight-qa-1

evidence · `01M6MQZ0CZJ23Y8EXN7HHDQ6K2` · actor [claude-code-5fde6012](../actor/claude-code-5fde6012.md)

## result

PASS

## statement

QA at dfff1995 with a binary built from the branch, in the same scratch repo (fake rex, RHIZA_*/REX_* unset). An interrupted plan run on a landed record: rhiza inbox printed 'Nothing waits on you' (no skipped item). Interrupted again after the daemon's one resume: rhiza inbox listed 'gave-up run <id> ... interrupted again after the daemon's one resume, next: rhiza run resume <id>', and --json gave kind gave-up with that next. With the journal ending failed ('boom'): inbox listed 'failed run <id> ... next: rhiza run resume --retry <id>', and running that command really retried the step, after which the run succeeded. Negative cases: a ship-change run interrupted on a landed change set whose journal never reached land was listed 'skipped ... its change set is landed, next: rhiza run abandon <id>', and a plan run on a withdrawn change set 'skipped ... its change set is withdrawn, next: rhiza run abandon <id>'. Scoped tests: go test ./internal/api -run Inbox and the inbox*, unaccepted* and run-resume* testscripts pass.

## Links

### Out

- verifies [inbox-landed-plan-run-in-flight](../capability/rhiza/orchestration/inbox/inbox-landed-plan-run-in-flight.md)

## History

- `01M6JZDKFVH14G853J0T0BPAKG` Fix Inbox and daemon treat an interrupted plan run whose record landed as resumable (replaces draft change set 01M6H5VQ231Z8V, which has no leaves: link the new change set relates-to it and withdraw the draft; serve milestone-25)
