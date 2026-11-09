# docs-reference-workflow-qa-1

evidence · `01M6ZF4Q0Z5N86BCRHZSQFGDFS` · actor [claude-code-d9ac25eb](../actor/claude-code-d9ac25eb.md)

## result

PASS

## statement

QA at 48d2315c with a scratchpad build (go build -trimpath ./cmd/rhiza) in a fresh git init + rhiza init + render repo: docs/reference/workflow-file.md has a section for each of the 8 step kinds internal/workflow defines (agent, command, wait, leaves, checks, signoff, workflow, arbiter), plus file keys, common step keys, role keys and prompt placeholders. go test ./internal/workflow/... and TestScripts/docs-reference-workflow pass.

## Links

### Out

- verifies [docs-reference-workflow](../capability/rhiza/documentation/docs-reference-workflow.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
