# apply-tolerates-concurrent-entries-verified

evidence · `01M7AQPQKZJ5BZT1R6F4Z8QJ22` · actor [claude-code-0aac2ac0](../actor/claude-code-0aac2ac0.md)

## result

PASS

## statement

At e61a2366e re-ran the Checks block: go test -trimpath ./cmd/rhiza -run 'TestScripts/(apply-concurrent|apply|check)$' exit 0; go test -trimpath ./internal/api -run TestApply exit 0; rhiza ship generate exit 0; rhiza render exit 0. Diff: scratchOps in internal/api/apply.go skips problems wrapping changeset.ErrConcurrent and returns any other fold problem, covering both the start ops and each appended item. New apply-concurrent.txtar builds concurrent entries as check.txtar does, shows the check warn with no fail, runs apply --dry-run and apply (add b, set a), change show lists node.create, render+check prints ok and only the warn; then a hand-written ops entry dropping an unwritten op makes apply --dry-run and apply fail with 'changeset: illegal revision' and the log entry count stays 2. In a scratch worktree with the old apply.go the new script fails at apply --dry-run, so it exercises the fix. No stubs or TODOs; only the leaf's files plus actor records changed.

## Links

### Out

- verifies [apply-tolerates-concurrent-entries](../capability/apply-tolerates-concurrent-entries.md)

## History

- `01M7AN1DTS8RCJ60YD93BADX4C` Fix rhiza apply refuses to write to a change set whose fold has a concurrent-entries warning (ErrConcurrent), though rhiza check only warns on it and single writes (rhiza set/add/link) succeed. scratchOps in internal/api/apply.go returns problems[0] for any fold problem. Repro: change set 01M79KH59D0HEME6B85PXXB056, where entries 01M7AHXPZYKTZVKB02T0NAN2YZ and 01M7AK05HZNK14DCH9GW0ARKYR both write op 01M79KPE0KMZHN6W97EDHJ5H21 (a later entry, 01M7AKMF541ZKACQ21FDF8XVEP, follows both). Any rhiza apply --change on it fails with 'changeset: concurrent entries: ...'. Expected: apply ignores ErrConcurrent as check does.
