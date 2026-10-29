# inbox-list-verified

evidence · `01M62GT93EFHR56WX2FD2GGA2N` · actor [claude-code-be1936d4](../actor/claude-code-be1936d4.md)

## result

PASS

## statement

At 1275837 I ran go test ./cmd/rhiza -run 'TestScripts/(inbox-list|describe-golden)$' -count=1 (ok) and make check (exit 0). inbox-list.txtar puts fixture journals of a run blocked on a question, a run stopped at sign-off, a failed run, a run interrupted again after a daemon-resumed entry and a done run under .git/rhiza/runs; rhiza inbox lists exactly question, signoff, failed, gave-up in that order (the ordered kinds are matched anchored, in text and --json) with the question, its tab and the failure detail, and never the done run. After an unblocked entry, an approved signoff entry, a real rhiza run abandon of the gave-up run and a new succeeded run on the failed run's change set, rhiza inbox prints 'Nothing waits on you' and --json gives empty items. I read the diff: internal/api/inbox.go skips abandoned failed runs (State.Abandoned) and failed runs with a newer run on their change set (run.Store.List sorts by start time), and journal.go adds only timestamps (ProcessAt, DaemonResumedAt, StoppedAt, EndedAt, SessionBlock.At).

## Links

### Out

- verifies [inbox-list](../capability/rhiza/orchestration/inbox/inbox-list.md)

## History

- `01M620GA5ZK6Z854D07NH0R8WX` The inbox: rhiza inbox lists every run's questions, sign-offs, failed runs and runs the daemon gave up on, with a recommendation for each question, opens each in its tab, takes sign-offs, and the daemon sends one Rex notification per new item.
