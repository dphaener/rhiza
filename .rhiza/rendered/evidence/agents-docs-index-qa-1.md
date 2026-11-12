# agents-docs-index-qa-1

evidence · `01M78407594PA6R1EBAMNHZ5HE` · actor [claude-code-20ebe8d3](../actor/claude-code-20ebe8d3.md)

## result

PASS

## statement

At head a5f095ccc: agents/README.md exists, links back to ../AGENTS.md and holds the index (no role docs yet); AGENTS.md links [agents/README.md](agents/README.md). TestAgentDocsConsistent passes on the repo; adding an unindexed agents/qa.md by hand makes it fail with 'agents/README.md does not index qa.md', and removing the file makes it pass again.

## Links

### Out

- verifies [agents-docs-index](../capability/rhiza/integrity/own-agent-docs/agents-docs-index.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
