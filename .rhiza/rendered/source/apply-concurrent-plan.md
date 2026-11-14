# apply-concurrent-plan

source · `01M7AN69WJ7YGN0EXYKXP4BKNG`

## quotes

```json
[]
```

## raw_sha256

f3e403d44e521bbab2b59b21d3cd1f0aa167a978da0c8cfb1bad140685ba5303

## session

339ae36a-8982-47d2-b716-01a1d5511141

## summary

Autonomous plan run; no owner interview and no questions were needed. Research: scratchOps in internal/api/apply.go fails apply on the first fold problem of any kind, while check (internal/check/check.go fold) downgrades changeset.ErrConcurrent to a warning; Fold still applies both concurrent entries in fold order, so the folded ops are usable. scratchOps runs both on the start ops and after each appended item, so both calls fail today. Plan: one leaf that skips ErrConcurrent problems in scratchOps (errors.Is) and still fails on any other fold problem, tested by a testscript that builds a concurrent pair the way check.txtar does (hide one entry while writing the other) and then runs rhiza apply on it. The change set serves m4-runs-hold-up, as the earlier bug fix to rhiza session message (01M7AA6QGX335J8XPJX3D6T4DJ) does: plan and review runs write through apply, and this failure needs hand repair.

## title

Plan of the apply ErrConcurrent bug fix (autonomous, no owner)

## Links

### In

- derived-from [apply-skips-concurrent-only](../decision/apply-skips-concurrent-only.md)
- derived-from [apply-concurrent-serves-runs-hold-up](../decision/apply-concurrent-serves-runs-hold-up.md)

## History

- `01M7AN1DTS8RCJ60YD93BADX4C` Fix rhiza apply refuses to write to a change set whose fold has a concurrent-entries warning (ErrConcurrent), though rhiza check only warns on it and single writes (rhiza set/add/link) succeed. scratchOps in internal/api/apply.go returns problems[0] for any fold problem. Repro: change set 01M79KH59D0HEME6B85PXXB056, where entries 01M7AHXPZYKTZVKB02T0NAN2YZ and 01M7AK05HZNK14DCH9GW0ARKYR both write op 01M79KPE0KMZHN6W97EDHJ5H21 (a later entry, 01M7AKMF541ZKACQ21FDF8XVEP, follows both). Any rhiza apply --change on it fails with 'changeset: concurrent entries: ...'. Expected: apply ignores ErrConcurrent as check does.
