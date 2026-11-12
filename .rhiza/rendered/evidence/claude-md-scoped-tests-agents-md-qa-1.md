# claude-md-scoped-tests-agents-md-qa-1

evidence · `01M784078HX35FD23ACW0J3XWW` · actor [claude-code-20ebe8d3](../actor/claude-code-20ebe8d3.md)

## result

PASS

## statement

At head a5f095ccc go test -trimpath -count=1 ./cmd/rhiza -run TestClaudeMDScopedTests passes, now reading AGENTS.md; AGENTS.md has the scoped commands and the never make check/make test/bare go test ./... rule, with no unqualified make check line.

## Links

### Out

- verifies [claude-md-scoped-tests](../capability/rhiza/integrity/self-check/claude-md-scoped-tests.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
