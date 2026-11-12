# agent-docs-consistent-verified

evidence · `01M783Q7CPFM6PSFXZ9CGCMHTG` · actor [claude-code-b7ecde6f](../actor/claude-code-b7ecde6f.md)

## result

PASS

## statement

At cb349dd15 ran go test -trimpath ./cmd/rhiza -run TestAgentDocsConsistent -v -count=1: TestAgentDocsConsistent passes on the repo's AGENTS.md and agents/, and TestAgentDocsConsistentFixtures passes its table: ok (no problems), unindexed doc (reports 'does not index b.md'), no AGENTS link (reports 'does not link agents/README.md'), broken link (reports 'links missing gone.md'). Read cmd/rhiza/agentdocs_test.go: it checks every agents/*.md is linked from agents/README.md, that AGENTS.md links agents/README.md, and that relative links in AGENTS.md and agents/*.md (outside code spans, skipping URLs and anchors) resolve. gofmt and go vet clean. make self-check fails only on stale .rhiza/rendered actor files from in-run actor records, not on this leaf.

## Links

### Out

- verifies [agent-docs-consistent](../capability/rhiza/integrity/own-agent-docs/agent-docs-consistent.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
