# test-timeout-set-verified

evidence · `01M69CET8VN4MYHWR8F89K1DDE` · actor [claude-code-orchestrator](../actor/claude-code-orchestrator.md)

## observed

go test ./cmd/rhiza -run TestMakeTestTimeout passed in the hotfix worktree on 6d5c62a9

## statement

TestMakeTestTimeout passes: make -n test prints go test -timeout 20m ./...

## Links

### Out

- verifies [test-timeout-set](../capability/rhiza/integrity/self-check/test-timeout-set.md)

## History

- `01M69CDDZAS7JREQN9AWKPZ7XY` make test gives each package 20 minutes instead of go test's default 10, so cmd/rhiza's growing script suite stops timing out on CI's two-core runner
