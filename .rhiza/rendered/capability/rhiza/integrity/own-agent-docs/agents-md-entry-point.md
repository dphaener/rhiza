# agents-md-entry-point

capability · `01M76NZQ9E4VVYH1CBT1ZDACE7` · parent [own-agent-docs](../own-agent-docs.md)

## acceptance

A test in cmd/rhiza reads CLAUDE.md and finds only the import line @@AGENTS.md, and reads AGENTS.md and finds the sections Ignore the ~/code/.claude/rules files, Stack and shape, Bug requests, Commands, Tests and Binaries; go test -trimpath ./cmd/rhiza -run that test passes

## statement

AGENTS.md at the repo root holds the repo guidance CLAUDE.md held (ignore the ~/code rules, stack and shape, bug requests, commands, tests, binaries), and CLAUDE.md is only the line @@AGENTS.md, so Claude Code and other agents read the same one file

## Links

### In

- verifies [agents-md-entry-point-verified](../../../../evidence/agents-md-entry-point-verified.md)
- verifies [agents-md-entry-point-qa-1](../../../../evidence/agents-md-entry-point-qa-1.md)
- verifies [pr-494-agents-md-entry-point](../../../../evidence/pr-494-agents-md-entry-point.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
