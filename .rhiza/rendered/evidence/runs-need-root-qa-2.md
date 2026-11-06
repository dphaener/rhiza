# runs-need-root-qa-2

evidence · `01M6S8XB577K0N41SS4X0Y1KC1` · actor [claude-code-9c5b9523](../actor/claude-code-9c5b9523.md)

## result

PASS

## statement

QA cycle 2, built binary: rhiza request and rhiza run start in an unregistered repo under a root exited 2 with next rhiza project add <path>, and in a repo outside any root exited 2 with next rhiza root init; rhiza check passed in both after render; request in registered app queued. go test -trimpath -short ./internal/api -run 'TestRunStart|TestDaemonPlanRequests|Amend|Process|TestRequest|Project|Root' now passes (cycle 1 blocker fixed by 56e58aa7).

## Links

### Out

- verifies [runs-need-root](../capability/rhiza/runtime/runs-need-root.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
