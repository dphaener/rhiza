# agents-docs-index-verified

evidence · `01M783MVMZ134JWHYE7ZGZB2T7` · actor [claude-code-a2ff416e](../actor/claude-code-a2ff416e.md)

## result

PASS

## statement

At ad06447d3 (detached checkout) I checked the acceptance's conditions directly, since the agent-docs-consistent test is a later leaf not yet built: agents/README.md exists; AGENTS.md links [agents/README.md](agents/README.md) on line 6, right under its intro, as where role instructions live; agents/ holds no *.md other than README.md, so the index (which states the one-line form '- [name.md](name.md): <which agent> reads it <when>' for future docs) is complete; README.md's ../AGENTS.md link resolves. go test -trimpath ./cmd/rhiza -run 'TestClaudeMD|TestAgents|TestMakeUpdate' passes (TestClaudeMDScopedTests, TestClaudeMDImportsAgentsMD, TestClaudeMDCommandsNameMakeUpdate, TestMakeUpdateDryRun). make self-check fails only on stale .rhiza/rendered actor files and index.md, which run close-out's rhiza render writes; no fold or model errors.

## Links

### Out

- verifies [agents-docs-index](../capability/rhiza/integrity/own-agent-docs/agents-docs-index.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
