# docs-commands-exist-qa-1

evidence · `01M6ZF4Q6WEGK1JQ58ADDA2GBG` · actor [claude-code-d9ac25eb](../actor/claude-code-d9ac25eb.md)

## result

PASS

## statement

QA at 48d2315c with a scratchpad build (go build -trimpath ./cmd/rhiza) in a fresh git init + rhiza init + render repo: a README naming `rhiza no-such-command` in backticks, plus a ```bash fence with 'rhiza change nosuch' and '$ rhiza bogus', and a docs/guide/a.md naming `rhiza help`: rhiza check exits 1 with one docs-command fail for each, naming the file (rhiza help really does exit 2 with unknown command). Real commands (`rhiza change new "x" --json`, version, tui app overview, api call check, describe run.findings, render), paths (.rhiza/log, rhiza.toml, github.com/x/rhiza) and unquoted prose are accepted. In this repo, rhiza check shows no docs-command finding. Gap, reported as minor: a ~~~ fence is not treated as code. TestScripts/docs-commands-exist and go test ./internal/check/... pass.

## Links

### Out

- verifies [docs-commands-exist](../capability/rhiza/documentation/docs-commands-exist.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
