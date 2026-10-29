# daemon-resumes-accepted-only-qa-1

evidence · `01M62ZXQ654TV6H233XC2VB1ND` · actor [claude-code-ea5d2a63](../actor/claude-code-ea5d2a63.md)

## result

PASS

## statement

At 2da4c0b, built rhiza to $TMPDIR and in a scratch repo made six real change sets (accepted, landed, rejected, withdrawn, proposed, draft) via change new/propose/accept/land/reject/withdraw, wrote an interrupted-run journal (dead pid) for each, and ran rhiza daemon --permission-mode manual --once --json with Rex faked as a logger. It reported resumed: [RUNACC] only; the fake rex log has one window new ... run resume RUNACC; only RUNACC.jsonl gained a daemon-resumed entry. Same result with autonomous runs.

## Links

### Out

- verifies [daemon-resumes-accepted-only](../capability/rhiza/orchestration/daemon/daemon-resumes-accepted-only.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
