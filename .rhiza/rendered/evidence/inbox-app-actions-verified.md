# inbox-app-actions-verified

evidence · `01M674P12Z97EM68P4EMYCV9SK` · actor [claude-code-a28ad899](../actor/claude-code-a28ad899.md)

## result

PASS

## statement

At 4e20b3c ran go test ./internal/tui/inboxapp -run 'Test(Open|Answer|Take|Approve|Retry|Abandon|Actions|Esc)' -v -count=1: all 8 pass. internal/tui/inboxapp/actions_test.go drives the Model with a fake Runner recording ops and a fake Exec recording argv (no real Rex or agent): o and enter call inbox.open q1; a + typed 'use pg' + enter calls session.message SESS1 'use pg'; t calls session.message SESS1 with the recommendation 'sqlite'; y on the sign-off calls inbox.approve s2; c + 'add tests' + enter calls inbox.request-changes s2 <file> whose findings[0].reason is 'add tests' (file removed after); R on the failed run execs 'rex window new RUNF00000 -c /repo --json -- /bin/rhiza run resume RUNF00000003 --retry' (gave-up: no --retry); x then n calls nothing, x then y calls run.abandon RUNF00000003; actions are ignored on other kinds. Read the diff: each act reloads the list; positional args match InboxOpen/Approve/RequestChanges, SessionMessage and RunAbandon request arg tags; run.resume has --retry; tui.go wires Exec/Self/Root/RexSession from the Ctx. make fmt-check vet and go test ./... pass (one make test run hit the known read-budget timing flake; reruns passed). make self-check fails only on stale .rhiza/rendered actor files recorded by this run, not on the leaf.

## Links

### Out

- verifies inbox-app-actions (retired)

## History

- `01M6732DJSQ14B2WTZYZ9MPRA5` An inbox app lists what waits on the owner, questions, sign-offs and failed runs, and acts on each in place: answer, approve, request changes, retry, abandon or open its tab
