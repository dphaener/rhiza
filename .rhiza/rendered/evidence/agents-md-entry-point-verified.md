# agents-md-entry-point-verified

evidence · `01M783J5YJRGE6FG94YTR84ZD4` · actor [claude-code-c7ea9f1e](../actor/claude-code-c7ea9f1e.md)

## result

PASS

## statement

At ef3679587: CLAUDE.md is exactly the line @@AGENTS.md (cat -A). AGENTS.md matches the old CLAUDE.md word for word apart from its title and intro line (diff against HEAD~1:CLAUDE.md), and holds the sections Ignore the ~/code/.claude/rules files, Stack and shape, Bug requests, Commands, Tests and Binaries. New test TestClaudeMDImportsAgentsMD in cmd/rhiza/claudemd_test.go checks both. The existing TestClaudeMDScopedTests and TestClaudeMDCommandsNameMakeUpdate now read AGENTS.md. go test -trimpath -count=1 ./cmd/rhiza -run TestClaudeMD -v: all 3 PASS. No other Go, Makefile, script or txtar file reads CLAUDE.md. make self-check fails only on unrendered actor files and index.md, which the ship step renders; this is not part of the leaf's acceptance.

## Links

### Out

- verifies [agents-md-entry-point](../capability/rhiza/integrity/own-agent-docs/agents-md-entry-point.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
