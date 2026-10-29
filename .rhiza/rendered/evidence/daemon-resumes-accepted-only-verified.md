# daemon-resumes-accepted-only-verified

evidence · `01M62Z2YX65ZZS0E7E6S2Q4HS0` · actor [claude-code-41f48c39](../actor/claude-code-41f48c39.md)

## result

PASS

## statement

At 1b01b58, go test ./internal/daemon -run Resumable -v passes TestResumableOnlyAcceptedChangeSets: of interrupted runs on accepted, landed, rejected, withdrawn, proposed and draft change sets, Resumable returns only run-accepted. go test ./cmd/rhiza -run TestScripts/daemon-resume passes: with real change sets, the interrupted run on the accepted one gets a Rex tab and a daemon-resumed entry, the landed and draft ones get neither. The diff wires Resumable's accepted predicate in daemonPoll1 to ws.changeSet(cs).State == StateAccepted. make fmt-check vet test exits 0; self-check fails only on stale rendered actor files from this run's sessions.

## Links

### Out

- verifies [daemon-resumes-accepted-only](../capability/rhiza/orchestration/daemon/daemon-resumes-accepted-only.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
