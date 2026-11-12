# unstable-ci-guidance-agents-md-qa-1

evidence · `01M78407AS1FNSMYG88B0QQRFF` · actor [claude-code-20ebe8d3](../actor/claude-code-20ebe8d3.md)

## result

PASS

## statement

At head a5f095ccc TestClaudeMDScopedTests (reading AGENTS.md) and go test -trimpath -count=1 ./internal/guide -run 'TestBuildText|TestReviewText' pass; AGENTS.md keeps the unstable-CI reproduce-locally rule unchanged.

## Links

### Out

- verifies [unstable-ci-guidance](../capability/rhiza/integrity/self-check/unstable-ci-guidance.md)

## History

- `01M76NWXV16NF94MBV6WYESRM0` AGENTS.md is this repo's agent entry point: CLAUDE.md's guidance moves into AGENTS.md and CLAUDE.md becomes an @AGENTS.md import, AGENTS.md links agents/README.md (the index of agents/ docs), a Go test keeps the index and relative links consistent, and the tests and capabilities that read CLAUDE.md read AGENTS.md
