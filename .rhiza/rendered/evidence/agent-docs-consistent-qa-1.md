# agent-docs-consistent-qa-1

evidence · `01M7840775E19HG53A3WRQJXTX` · actor [claude-code-20ebe8d3](../actor/claude-code-20ebe8d3.md)

## result

PASS

## statement

At head a5f095ccc ran go test -trimpath -count=1 ./cmd/rhiza -run TestAgentDocsConsistent -v: the repo case and the fixture table (ok, unindexed doc, no AGENTS link, broken link) pass. Mutated the real docs by hand: an unindexed agents/qa.md, AGENTS.md's index link changed to agents/READ.md, and [x](missing.md) appended to agents/README.md each made TestAgentDocsConsistent fail with the matching problem; a missing link inside a code span did not; restoring the files made it pass.

## Links

### Out

- verifies [agent-docs-consistent](../capability/rhiza/integrity/own-agent-docs/agent-docs-consistent.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
