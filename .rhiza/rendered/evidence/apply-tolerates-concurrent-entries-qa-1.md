# apply-tolerates-concurrent-entries-qa-1

evidence · `01M7AR32RGPMNDP5MSWWPXMBW0` · actor [claude-code-864834fc](../actor/claude-code-864834fc.md)

## result

PASS

## statement

QA at 8aad96ca3, binary built fresh with make build to a scratch path, run in a hermetic scratch repo (env -i, own HOME). Built a concurrent pair (hid an entry while a second rhiza set wrote op a, restored it): check printed ok with only the concurrent-entries warn. apply --dry-run (add b, set a, link b depends-on a) exit 0 and wrote nothing; apply exit 0, wrote one entry whose after lists both concurrent heads; change show listed both node.create and the link; render+check ok with only the warn. Repeated applies on it work (a no-op one says the writes cancel out); a second concurrent pair in the same change set still lets apply write. A canonical hand-written ops entry dropping an unwritten op (ErrRevision) in another change set: apply --dry-run and apply both exit 1 with 'changeset: illegal revision', log entry count unchanged; the same with concurrent pairs plus that entry in one change set also fails and writes nothing. go test -trimpath ./cmd/rhiza -run 'TestScripts/(apply-concurrent|apply|check)$' and ./internal/api -run TestApply pass.

## Links

### Out

- verifies [apply-tolerates-concurrent-entries](../capability/apply-tolerates-concurrent-entries.md)

## History

- `01M7AN1DTS8RCJ60YD93BADX4C` Fix rhiza apply refuses to write to a change set whose fold has a concurrent-entries warning (ErrConcurrent), though rhiza check only warns on it and single writes (rhiza set/add/link) succeed. scratchOps in internal/api/apply.go returns problems[0] for any fold problem. Repro: change set 01M79KH59D0HEME6B85PXXB056, where entries 01M7AHXPZYKTZVKB02T0NAN2YZ and 01M7AK05HZNK14DCH9GW0ARKYR both write op 01M79KPE0KMZHN6W97EDHJ5H21 (a later entry, 01M7AKMF541ZKACQ21FDF8XVEP, follows both). Any rhiza apply --change on it fails with 'changeset: concurrent entries: ...'. Expected: apply ignores ErrConcurrent as check does.
