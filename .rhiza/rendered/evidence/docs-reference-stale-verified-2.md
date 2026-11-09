# docs-reference-stale-verified-2

evidence · `01M6ZH8CV8N2F602PD6Q49JBTS` · actor [claude-code-553426f7](../actor/claude-code-553426f7.md)

## result

PASS

## statement

At e7eac5e7 I ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(docs-reference-stale|describe-golden|docs)$' -count=1: docs-reference-stale passes. It appends to a generated docs/reference/config.md, sees rhiza check fail naming config.md (not cli.md) as differing with next: rhiza docs, then pass after rhiza docs; a deleted cli.md is reported missing; extra files and a project's own docs/reference without the generated marker are left alone, and rhiza docs refuses to write into it. internal/docs, internal/check and internal/api tests pass. make self-check on this repo reports no docs/reference failure, so the committed generated reference is fresh; its only failures are .rhiza/rendered staleness from the run's own record commits. Diff read: docs.Generated gates the check on the Marker, shared by all three generators.

## Links

### Out

- verifies [docs-reference-stale](../capability/rhiza/documentation/docs-reference-stale.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
