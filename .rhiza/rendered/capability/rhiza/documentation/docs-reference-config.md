# docs-reference-config

capability · `01M6W5GZXB8YMB2G6142DGQTS7` · parent [documentation](../documentation.md)

## acceptance

a testscript checks rhiza.autonomy.plan appears with its default; a unit test fails if a settings definition has no description

## statement

rhiza docs writes docs/reference/config.md: every rhiza.toml setting with its type, default and description, from the settings definitions

## Links

### Out

- decided-by [m30-generated-reference](../../../decision/m30-generated-reference.md)

### In

- verifies [docs-reference-config-verified](../../../evidence/docs-reference-config-verified.md)
- verifies [docs-reference-config-qa-1](../../../evidence/docs-reference-config-qa-1.md)
- verifies [pr-412-generated-reference](../../../evidence/pr-412-generated-reference.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
