# docs-commands-exist

capability · `01M6W5H01JDC5QZZECFJSXQSTK` · parent [documentation](../documentation.md)

## acceptance

a testscript feeds a doc naming rhiza no-such-command and sees the check fail; the check runs in make check

## statement

A check fails when README.md or docs/ names a rhiza command (in backticks or a code block) that the registry does not have

## Links

### Out

- decided-by [m30-generated-reference](../../../decision/m30-generated-reference.md)

### In

- verifies [docs-commands-exist-verified](../../../evidence/docs-commands-exist-verified.md)
- verifies [docs-commands-exist-qa-1](../../../evidence/docs-commands-exist-qa-1.md)
- verifies [pr-412-generated-reference](../../../evidence/pr-412-generated-reference.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
