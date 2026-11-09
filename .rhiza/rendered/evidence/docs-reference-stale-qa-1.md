# docs-reference-stale-qa-1

evidence · `01M6ZF4Q56262J8EZMDT9VY8J6` · actor [claude-code-d9ac25eb](../actor/claude-code-d9ac25eb.md)

## result

PASS

## statement

QA at 48d2315c with a scratchpad build (go build -trimpath ./cmd/rhiza) in a fresh git init + rhiza init + render repo: after rhiza docs, rhiza check passes. I appended a line to config.md, deleted workflow-file.md and added my own notes.md: rhiza check exits 1 with 'fail render docs/reference/config.md: differs from what rhiza docs writes' and 'fail render docs/reference/workflow-file.md: missing', both with next: rhiza docs, and notes.md is ignored. After rhiza docs it rewrote the two files and check exits 0 again. TestScripts/docs-reference-stale passes.

## Links

### Out

- verifies [docs-reference-stale](../capability/rhiza/documentation/docs-reference-stale.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
