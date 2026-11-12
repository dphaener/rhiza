# agents-md-entry-point-qa-1

evidence · `01M784073ZCNJPY5KRTX6XDN7X` · actor [claude-code-20ebe8d3](../actor/claude-code-20ebe8d3.md)

## result

PASS

## statement

At head a5f095ccc: CLAUDE.md is the single line @@AGENTS.md and this Claude Code session loaded AGENTS.md's full guidance through that import; diff of the old CLAUDE.md body against AGENTS.md shows the six sections moved verbatim. go test -trimpath -count=1 ./cmd/rhiza -run TestClaudeMDImportsAgentsMD passes. make self-check from a clean detached checkout of the head: ok.

## Links

### Out

- verifies [agents-md-entry-point](../capability/rhiza/integrity/own-agent-docs/agents-md-entry-point.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
