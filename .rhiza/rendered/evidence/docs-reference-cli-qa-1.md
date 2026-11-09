# docs-reference-cli-qa-1

evidence · `01M6ZF4PR9C9EKTV8K83XWXWSY` · actor [claude-code-d9ac25eb](../actor/claude-code-d9ac25eb.md)

## result

PASS

## statement

QA at 48d2315c with a scratchpad build (go build -trimpath ./cmd/rhiza) in a fresh git init + rhiza init + render repo: rhiza docs wrote docs/reference/cli.md (2667 lines, plus config.md and workflow-file.md), exit 0. The '## rhiza request' section has the usage line, <text>, --kind, --milestone, --workflow and --ref with descriptions, and the --json output fields; a shared Exit codes section lists 0-4. Two rhiza docs --out runs gave byte-identical trees, and a re-run reports '3 unchanged'. go test ./internal/docs/... and TestScripts/docs-reference-cli pass.

## Links

### Out

- verifies [docs-reference-cli](../capability/rhiza/documentation/docs-reference-cli.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
