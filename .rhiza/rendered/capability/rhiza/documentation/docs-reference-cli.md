# docs-reference-cli

capability · `01M6W5GZVA1QCBXVCK83RMHFW8` · parent [documentation](../documentation.md)

## acceptance

a testscript runs rhiza docs into a temp dir and checks the request command's section lists --kind and --milestone; go test of the docs package

## statement

rhiza docs writes docs/reference/cli.md from the op registry: every command with its arguments, flags, descriptions, output and exit codes

## Links

### Out

- decided-by [m30-generated-reference](../../../decision/m30-generated-reference.md)

### In

- verifies [docs-reference-cli-verified](../../../evidence/docs-reference-cli-verified.md)
- verifies [docs-reference-cli-qa-1](../../../evidence/docs-reference-cli-qa-1.md)
- verifies [pr-412-generated-reference](../../../evidence/pr-412-generated-reference.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
