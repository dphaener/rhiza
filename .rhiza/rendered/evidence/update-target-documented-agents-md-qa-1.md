# update-target-documented-agents-md-qa-1

evidence · `01M78407J3K9ZA2X3Q3PS6EW3W` · actor [claude-code-20ebe8d3](../actor/claude-code-20ebe8d3.md)

## result

PASS

## statement

At head a5f095ccc make help lists update, and go test -trimpath -count=1 ./cmd/rhiza -run TestClaudeMDCommandsNameMakeUpdate passes, now reading AGENTS.md's Commands block.

## Links

### Out

- verifies [update-target-documented](../capability/rhiza/rex-install-command/update-target-documented.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
