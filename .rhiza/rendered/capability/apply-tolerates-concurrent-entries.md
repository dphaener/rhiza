# apply-tolerates-concurrent-entries

capability · `01M7AN6YYW3GR94680716HY50V`

## acceptance

A new testscript cmd/rhiza/testdata/script/apply-concurrent.txtar builds a change set with two concurrent entries writing the same op the way check.txtar's concurrent section does (hide one entry's file while a second rhiza set writes, then restore it), shows rhiza check prints the concurrent-entries warn and no fail, then runs rhiza apply on a file that adds a capability and sets a field on the existing one: apply --dry-run and apply both succeed, rhiza change show lists the new op, and rhiza check after rhiza render still prints ok with only the warn. The same script then adds a hand-written ops entry to another change set's log that drops an op ID never written (an ErrRevision fold problem) and runs rhiza apply --dry-run and rhiza apply against that change set: both fail with 'changeset: illegal revision' in stderr and the change set's log gains no entry. go test -trimpath ./cmd/rhiza -run 'TestScripts/(apply-concurrent|apply|check)$' and go test -trimpath ./internal/api -run TestApply pass.

## statement

rhiza apply (and apply --dry-run) writes to a change set whose fold reports only concurrent-entries problems (changeset.ErrConcurrent), both from its start ops and after each appended item, through scratchOps in internal/api/apply.go, while any other fold problem still fails the apply with that problem's error

## Links

### Out

- decided-by [apply-skips-concurrent-only](../decision/apply-skips-concurrent-only.md)

### In

- verifies [apply-tolerates-concurrent-entries-verified](../evidence/apply-tolerates-concurrent-entries-verified.md)
- verifies [apply-tolerates-concurrent-entries-qa-1](../evidence/apply-tolerates-concurrent-entries-qa-1.md)
- verifies [pr-556-fix-rhiza-apply-refuses-to](../evidence/pr-556-fix-rhiza-apply-refuses-to.md)

## History

- `01M7AN1DTS8RCJ60YD93BADX4C` Fix rhiza apply refuses to write to a change set whose fold has a concurrent-entries warning (ErrConcurrent), though rhiza check only warns on it and single writes (rhiza set/add/link) succeed. scratchOps in internal/api/apply.go returns problems[0] for any fold problem. Repro: change set 01M79KH59D0HEME6B85PXXB056, where entries 01M7AHXPZYKTZVKB02T0NAN2YZ and 01M7AK05HZNK14DCH9GW0ARKYR both write op 01M79KPE0KMZHN6W97EDHJ5H21 (a later entry, 01M7AKMF541ZKACQ21FDF8XVEP, follows both). Any rhiza apply --change on it fails with 'changeset: concurrent entries: ...'. Expected: apply ignores ErrConcurrent as check does.
