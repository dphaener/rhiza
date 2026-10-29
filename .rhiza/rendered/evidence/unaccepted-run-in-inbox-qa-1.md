# unaccepted-run-in-inbox-qa-1

evidence · `01M62ZXSPP8S749RD0M93MJ47E` · actor [claude-code-ea5d2a63](../actor/claude-code-ea5d2a63.md)

## result

PASS

## statement

At 2da4c0b with the built binary, after daemon --once skipped them, rhiza inbox lists one skipped item for each of the landed, rejected, withdrawn, proposed and draft runs ('interrupted, and not resumed: its change set is landed', 'next: rhiza run abandon RUNLND', ...), none for the resumed accepted run; the daemon emitted a blocked OSC 7501 notice per skipped run; rhiza run abandon RUNREJ removes its item. Autonomous runs get no inbox item (pre-existing autonomous exclusion).

## Links

### Out

- verifies [unaccepted-run-in-inbox](../capability/rhiza/orchestration/inbox/unaccepted-run-in-inbox.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
