# import-records-built

capability · `01M6J3C5Q7TJ7HR8GP6PD4JW1B` · parent [guide](../guide.md)

## acceptance

The import topic states the one-record-per-area rule, the statement and acceptance rules and the follow-up anchoring step; a testscript follows those steps on a small fixture repo with tests and gets one landed record change set per area, every capability with a statement and an acceptance naming a test file, then a follow-up change set whose anchors are linked implements to those capabilities, and rhiza check passes.

## statement

The import procedure records what a repo already does as record change sets, one per area, whose capabilities each state one user-visible behaviour and give as acceptance the repo's existing tests that prove it, then anchors each area's capabilities to the code that implements them in a follow-up change set with rhiza add anchor.

## Links

### Out

- decided-by [import-anchors-later](../../../decision/import-anchors-later.md)

### In

- verifies [import-records-built-verified](../../../evidence/import-records-built-verified.md)
- verifies [import-records-built-qa-1](../../../evidence/import-records-built-qa-1.md)
- verifies [pr-328-import-procedure](../../../evidence/pr-328-import-procedure.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
