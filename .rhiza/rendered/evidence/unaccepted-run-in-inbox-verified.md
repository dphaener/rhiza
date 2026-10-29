# unaccepted-run-in-inbox-verified

evidence · `01M62Z7BM1GK3YJ0TH1CJ0MZVH` · actor [claude-code-fa907666](../actor/claude-code-fa907666.md)

## result

PASS

## statement

At 4f54a51 ran go test ./cmd/rhiza -run 'TestScripts/(daemon-resume|describe-golden|inbox.*)$' -count=1: all pass. daemon-resume.txtar has interrupted runs on a landed (LANDED00001) and a draft (DRAFT000001) change set; daemon --once opens no Rex tab and journals no daemon-resumed for either, then rhiza inbox lists 'skipped run LANDED00001, change set ...', 'interrupted, and not resumed: its change set is landed', 'next: rhiza run abandon LANDED00001' (and the draft one), not the resumed INTERRUPTED1; abandoning LANDED00001 clears it. Diff (internal/api/inbox.go) adds InboxSkipped for a dead, not daemon-resumed, unblocked run whose change set state is known and not accepted, matching the daemon's accepted-only resume. make check: fmt, vet and go test pass; self-check fails only on a stale .rhiza/rendered from an earlier run's actor record, not this leaf.

## Links

### Out

- verifies [unaccepted-run-in-inbox](../capability/rhiza/orchestration/inbox/unaccepted-run-in-inbox.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
