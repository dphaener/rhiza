# import-retires-docs

capability · `01M6J3CQTM8JJ1FG394JB40RHN` · parent [guide](../guide.md)

## acceptance

The import topic names which documents are retired and which are kept; a testscript retires a fixture PLAN.md: rhiza source add with --raw succeeds, the file is deleted in the same commit, and rhiza check passes.

## statement

The import procedure retires a planning document by adding it as a Source (a scrubbed summary, the owner's typed words found in it quoted verbatim, the document hashed with --raw) and deleting it in the same change set, and keeps reference, design and operations documents that capabilities link to.

## Links

### Out

- decided-by [import-retires-planning-docs](../../../decision/import-retires-planning-docs.md)

### In

- verifies [import-retires-docs-verified](../../../evidence/import-retires-docs-verified.md)
- verifies [import-retires-docs-qa-1](../../../evidence/import-retires-docs-qa-1.md)
- verifies [pr-328-import-procedure](../../../evidence/pr-328-import-procedure.md)

## History

- `01M6J36X6EC3E70H0QCW17GREP` Import procedure: rhiza guide import teaches agents to bring an existing repo with built code into Rhiza, the steps Nebari's import follows
