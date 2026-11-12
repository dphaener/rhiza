# agents-docs-index

capability · `01M76NZQA1943A46ZY93ZKTZJX` · parent [own-agent-docs](../own-agent-docs.md)

## acceptance

The consistency test (agent-docs-consistent) passes on the repo: agents/README.md exists, AGENTS.md links it, and every agents/*.md other than README.md is listed in it

## statement

agents/README.md indexes every doc under agents/, one line each saying which agent reads it and when, and AGENTS.md links to it near its top as where role instructions live

## Links

### In

- verifies [agents-docs-index-verified](../../../../evidence/agents-docs-index-verified.md)
- verifies [agents-docs-index-qa-1](../../../../evidence/agents-docs-index-qa-1.md)
- verifies [pr-494-agents-md-entry-point](../../../../evidence/pr-494-agents-md-entry-point.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
