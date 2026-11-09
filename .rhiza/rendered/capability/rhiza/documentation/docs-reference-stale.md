# docs-reference-stale

capability · `01M6W5H0026NYNDG8X1JYS6R7C` · parent [documentation](../documentation.md)

## acceptance

a testscript edits a generated file and sees the staleness check fail naming it, then pass after rhiza docs

## statement

make check fails when docs/reference/ differs from what rhiza docs writes, naming the stale files, as it does for .rhiza/rendered/

## Links

### Out

- decided-by [m30-generated-reference](../../../decision/m30-generated-reference.md)

### In

- verifies [docs-reference-stale-verified](../../../evidence/docs-reference-stale-verified.md)
- verifies [docs-reference-stale-qa-1](../../../evidence/docs-reference-stale-qa-1.md)
- verifies [docs-reference-stale-verified-2](../../../evidence/docs-reference-stale-verified-2.md)
- verifies [pr-412-generated-reference](../../../evidence/pr-412-generated-reference.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
