# unaccepted-run-in-inbox-verified-2

evidence · `01M630AMWZPZYFVDFYTKR3JBHZ` · actor [claude-code-a1465882](../actor/claude-code-a1465882.md)

## result

PASS

## statement

At 6929af3 ran go test ./cmd/rhiza -run 'TestScripts/daemon-resume$' -count=1 (ok), plus go test ./internal/daemon ./internal/api (ok), fmt-check and vet (clean). In daemon-resume.txtar, interrupted runs LANDED00001 and LANDEDAUTO1 (autonomous) on a landed change set get no Rex tab and no daemon-resumed entry. rhiza inbox lists both as 'skipped run ..., interrupted, and not resumed: its change set is landed' with 'next: rhiza run abandon <run>', and abandoning one clears it. Read the diff: internal/api/inbox.go now runs the skipped check before the autonomous skip, using changeset.State.Terminal() (the same test as run resume), which fixes the review's autonomous-run finding. make check's only failure is self-check's stale render of another session's actor record, which is unrelated to this commit.

## Links

### Out

- verifies [unaccepted-run-in-inbox](../capability/rhiza/orchestration/inbox/unaccepted-run-in-inbox.md)

## History

- `01M62MKZ6SA46DYC1K2R0W254S` The daemon resumes only interrupted runs whose change set is still accepted: a run left interrupted on a change set that has since landed, been withdrawn or rejected is never resumed.
