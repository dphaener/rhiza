# apply-concurrent-serves-runs-hold-up

decision · `01M7AN6YGQZ03KYB3R430Z5N5T`

## statement

The apply ErrConcurrent bug fix serves m4-runs-hold-up, as the earlier rhiza session message bug fix does: plan, revise and review runs write through rhiza apply, and a change set left with a concurrent pair after a merge otherwise needs hand repair

## Links

### Out

- derived-from [apply-concurrent-plan](../source/apply-concurrent-plan.md)

## History

- `01M7AN1DTS8RCJ60YD93BADX4C` Fix rhiza apply refuses to write to a change set whose fold has a concurrent-entries warning (ErrConcurrent), though rhiza check only warns on it and single writes (rhiza set/add/link) succeed. scratchOps in internal/api/apply.go returns problems[0] for any fold problem. Repro: change set 01M79KH59D0HEME6B85PXXB056, where entries 01M7AHXPZYKTZVKB02T0NAN2YZ and 01M7AK05HZNK14DCH9GW0ARKYR both write op 01M79KPE0KMZHN6W97EDHJ5H21 (a later entry, 01M7AKMF541ZKACQ21FDF8XVEP, follows both). Any rhiza apply --change on it fails with 'changeset: concurrent entries: ...'. Expected: apply ignores ErrConcurrent as check does.
