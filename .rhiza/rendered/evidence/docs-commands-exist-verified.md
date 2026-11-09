# docs-commands-exist-verified

evidence · `01M6ZETPA50RNDHTMP38B263VD` · actor [claude-code-5a3e087f](../actor/claude-code-5a3e087f.md)

## result

PASS

## statement

Verified at 33da0496. go test -trimpath ./cmd/rhiza -run 'TestScripts/(docs-commands-exist|describe-golden)$' passes; docs-commands-exist.txtar puts a README naming `rhiza no-such-command` in backticks and a docs/guide.md naming `rhiza change nosuch` in a code block, and sees rhiza check fail naming each file, and pass for real commands, groups, paths and prose. The attempt-1 finding is fixed: a scratchpad build in a fresh init+render repo passes rhiza check with a README naming `rhiza version`, `rhiza api call check`, `rhiza tui app overview` (exit 0), and fails with `rhiza no-such-command` (exit 1, fail docs-command README.md). api.FrameCommands now feeds both check and the help text, so the two cannot drift (help still lists the three). go test ./internal/check/... and ./internal/api -run 'Check|Docs' pass. The rule runs inside rhiza check, which make check calls through self-check; make self-check on this repo shows no docs-command finding (only expected mid-run stale-render fails). The wider internal/api suite failed here only on the known usage-state leak into ~/code/.rhiza-root, unrelated to this diff.

## Links

### Out

- verifies [docs-commands-exist](../capability/rhiza/documentation/docs-commands-exist.md)

## History

- `01M6W5E3JQXD62WN4RHG546KRX` Generated reference: rhiza docs writes docs/reference/ (CLI commands, rhiza.toml settings, workflow file format) from the code, and make check fails when it is stale or the docs name a command that does not exist
