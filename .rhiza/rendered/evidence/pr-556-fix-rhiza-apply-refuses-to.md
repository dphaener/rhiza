# pr-556-fix-rhiza-apply-refuses-to

evidence · `01M7AT0TT8N59MZTFZ95G5YVG3` · actor [dphaener](../actor/dphaener.md)

## head

7caf1b69fb546e494c635cd3fd8e7b653dd7240a

## pr

https://github.com/dphaener/rhiza/pull/556

## result

9 checks passed: changes SUCCESS, go-install SUCCESS, model SUCCESS, lint SKIPPED, scripts SKIPPED, api SKIPPED, rest SKIPPED, budgets SKIPPED, check SUCCESS

## statement

Pull request #556 (https://github.com/dphaener/rhiza/pull/556) passed CI at 7caf1b6

## status

verified

## Links

### Out

- verifies [apply-tolerates-concurrent-entries](../capability/apply-tolerates-concurrent-entries.md)

## History

- `01M7AN1DTS8RCJ60YD93BADX4C` Fix rhiza apply refuses to write to a change set whose fold has a concurrent-entries warning (ErrConcurrent), though rhiza check only warns on it and single writes (rhiza set/add/link) succeed. scratchOps in internal/api/apply.go returns problems[0] for any fold problem. Repro: change set 01M79KH59D0HEME6B85PXXB056, where entries 01M7AHXPZYKTZVKB02T0NAN2YZ and 01M7AK05HZNK14DCH9GW0ARKYR both write op 01M79KPE0KMZHN6W97EDHJ5H21 (a later entry, 01M7AKMF541ZKACQ21FDF8XVEP, follows both). Any rhiza apply --change on it fails with 'changeset: concurrent entries: ...'. Expected: apply ignores ErrConcurrent as check does.
