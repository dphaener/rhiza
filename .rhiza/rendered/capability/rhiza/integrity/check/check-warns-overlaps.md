# check-warns-overlaps

capability · `01M69DHG2Y91KWW8R5BJEFQDH7` · parent [check](../check.md)

## acceptance

A case in internal/check/check_test.go and a testscript where two proposed change sets set the same field of one node: rhiza check exits 0 and its warning names both change sets; existing scripts whose output changes are updated

## statement

rhiza check warns under the in-flight rule on each pair of proposed or accepted change sets that write overlapping aspects of the same target, naming both change sets, with next rhiza slice --change <A> --view what-if

## Links

### Out

- decided-by [overlaps-reuse-model](../../../../decision/overlaps-reuse-model.md)

### In

- verifies [check-warns-overlaps-verified](../../../../evidence/check-warns-overlaps-verified.md)
- verifies [check-warns-overlaps-qa-1](../../../../evidence/check-warns-overlaps-qa-1.md)
- verifies [pr-149-genesis-detects-conflicting](../../../../evidence/pr-149-genesis-detects-conflicting.md)

## History

- `01M5SWTNWB9GEJ5AA8R7K867KB` Genesis detects conflicting decisions across change sets in flight and surfaces them to the owner.
