# own-repo-bug-flow-note-agents-md-qa-1

evidence · `01M78407EH7Z98Q1MEA461DBFC` · actor [claude-code-20ebe8d3](../actor/claude-code-20ebe8d3.md)

## result

PASS

## statement

At head a5f095ccc AGENTS.md's Bug requests section names rhiza request --kind bug, .rhiza/dispatch.toml, plan-bug and bugfix. With a binary built from the head: rhiza request --help shows --kind; rhiza workflow describe plan-bug and bugfix succeed in this checkout and show autonomy plan auto and accept auto; .rhiza/dispatch.toml and both workflow files exist.

## Links

### Out

- verifies [own-repo-bug-flow-note](../capability/rhiza/workflows/workflow-dispatch/own-repo-bug-flow-note.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
