# apply-skips-concurrent-only

decision · `01M7AN6YB7KKN3XYPSG09XCS3T`

## reasoning

Fold still applies both concurrent entries in fold order, so the ops it returns are the same ones check, set, add and link work from; only ErrConcurrent is an honest-merge warning, the others mean the entry did not apply

## statement

rhiza apply's scratch fold (scratchOps in internal/api/apply.go) skips fold problems that wrap changeset.ErrConcurrent, matched with errors.Is, as rhiza check downgrades them to warnings, and still fails on the first other fold problem (ErrBody, ErrOrder, ErrTransition, ErrRevision), so an item that cannot be applied as built is still refused

## Links

### Out

- derived-from [apply-concurrent-plan](../source/apply-concurrent-plan.md)

### In

- decided-by [apply-tolerates-concurrent-entries](../capability/apply-tolerates-concurrent-entries.md)

## History

- `01M7AN1DTS8RCJ60YD93BADX4C` Fix rhiza apply refuses to write to a change set whose fold has a concurrent-entries warning (ErrConcurrent), though rhiza check only warns on it and single writes (rhiza set/add/link) succeed. scratchOps in internal/api/apply.go returns problems[0] for any fold problem. Repro: change set 01M79KH59D0HEME6B85PXXB056, where entries 01M7AHXPZYKTZVKB02T0NAN2YZ and 01M7AK05HZNK14DCH9GW0ARKYR both write op 01M79KPE0KMZHN6W97EDHJ5H21 (a later entry, 01M7AKMF541ZKACQ21FDF8XVEP, follows both). Any rhiza apply --change on it fails with 'changeset: concurrent entries: ...'. Expected: apply ignores ErrConcurrent as check does.
