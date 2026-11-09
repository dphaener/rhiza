# docs-reference-config-qa-1

evidence · `01M6ZF4PV4F56Y9C3DTANEZBSQ` · actor [claude-code-d9ac25eb](../actor/claude-code-d9ac25eb.md)

## result

PASS

## statement

QA at 48d2315c with a scratchpad build (go build -trimpath ./cmd/rhiza) in a fresh git init + rhiza init + render repo: docs/reference/config.md has '### rhiza.autonomy.plan' with 'Type: auto or owner' and 'Default: owner', which matches the settings definition in internal/api/daemon_settings.go. TestScripts/docs-reference-config passes.

## Links

### Out

- verifies [docs-reference-config](../capability/rhiza/documentation/docs-reference-config.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
