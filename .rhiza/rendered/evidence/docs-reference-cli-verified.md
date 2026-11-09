# docs-reference-cli-verified

evidence · `01M6ZBYD7QDFRG87HHHHKRVNFS` · actor [claude-code-0d6fde16](../actor/claude-code-0d6fde16.md)

## result

PASS

## statement

At ae8ad5d3: go test -trimpath ./internal/docs/... passes; go test -trimpath ./cmd/rhiza -run 'TestScripts/(docs-reference-cli|describe-golden)$' -count=1 passes. The docs-reference-cli testscript runs rhiza docs --out into $WORK, checks the '## rhiza request' section lists --kind and --milestone with a pattern whose lines may not start with '#', so it cannot run into a later section, checks the exit codes, checks a fixture whose request section lacks --milestone fails, and checks a second run is byte-identical. I built the binary, ran rhiza docs --out into a scratch dir, and read the request section: usage line, arguments, flags with descriptions, --json output fields; the doc has 114 command sections and a shared Exit codes section. Mutation check for the attempt-2 finding: I deleted the request section's --milestone line from the generated cli.md and ran the testscript's --milestone regexp (Go regexp, with (?m)) on it: no match (it matches the unedited file), and the --kind regexp still matches. So the finding is fixed.

## Links

### Out

- verifies [docs-reference-cli](../capability/rhiza/documentation/docs-reference-cli.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
