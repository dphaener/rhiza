# docs-reference-stale-verified

evidence · `01M6ZC5KNRMDNGNDMS017WN717` · actor [claude-code-f6574bbc](../actor/claude-code-f6574bbc.md)

## result

PASS

## statement

At 08b59b62 ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/docs-reference-stale$' -v: PASS. The script runs rhiza docs and check (pass), appends to docs/reference/config.md and sees rhiza check fail with stdout naming docs/reference/config.md (not cli.md), 'differs from what rhiza docs writes' and 'next: rhiza docs', then passes after rhiza docs; a deleted cli.md fails as 'missing' and passes after rhiza docs; an extra notes.md is ignored. Read the diff: check.docs() compares api.docsFiles() (the same map runDocs now writes) with docs/reference under the render rule, skipped when the folder is absent. go test ./internal/check (ok), go test -short ./internal/api -run '^Test(Check|Docs)' (ok), go vet and gofmt -l on the touched packages are clean. The full internal/api run failed only on things this change doesn't touch: a budget test timed out waiting on machine load, and the rootguard flagged ~/code/.rhiza-root/usage.live.json changing during the run. make self-check fails only on .rhiza/rendered being stale for the run's untracked actor records, which is expected before landing; there is no docs/reference finding.

## Links

### Out

- verifies [docs-reference-stale](../capability/rhiza/documentation/docs-reference-stale.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
