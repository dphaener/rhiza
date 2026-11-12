# agent-docs-consistent

capability · `01M76NZQB0E9EP3387YTWEW6MC` · parent [own-agent-docs](../own-agent-docs.md)

## acceptance

go test -trimpath ./cmd/rhiza -run <the test> passes on the repo, and its table cases fail on a fixture with an unindexed agents/ doc, one without the AGENTS.md link and one with a broken relative link

## statement

A Go test fails when an agents/*.md file is missing from agents/README.md's index, when AGENTS.md does not link agents/README.md, or when a relative link in AGENTS.md or under agents/ points at a file that does not exist

## Links

### In

- verifies [agent-docs-consistent-verified](../../../../evidence/agent-docs-consistent-verified.md)
- verifies [agent-docs-consistent-qa-1](../../../../evidence/agent-docs-consistent-qa-1.md)
- verifies [pr-494-agents-md-entry-point](../../../../evidence/pr-494-agents-md-entry-point.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
