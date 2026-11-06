# runs-need-root-verified-2

evidence · `01M6S7SKN007ZRCXQ8MYJCNF4K` · actor [claude-code-a6df410e](../actor/claude-code-a6df410e.md)

## result

PASS

## statement

Attempt 2 at 56e58aa7 (fix for blocker f1: session fixture now writes a rhiza.toml registering its repo). Ran `go test -trimpath -short ./internal/api`: ok (the 21 run/request tests that failed with 'not inside a Rhiza root' now pass). Ran `go test -trimpath ./cmd/rhiza -run TestScripts/(runs-need-root|run-start.*|plan-requests|daemon-milestone-request|daemon-request-slot|tui-request-app|testscript-root|rhiza-root|root-init)$`: ok. runs-need-root.txtar asserts rhiza request and rhiza run start fail with next: rhiza project add in an unregistered repo and next: rhiza root init under $NOROOT, with render and check succeeding in both. A scratch build outside any root confirmed rhiza request and rhiza run start exit exactly 2 with next: rhiza root init, and render then check exit 0. go vet ./internal/api and gofmt clean.

## Links

### Out

- verifies [runs-need-root](../capability/rhiza/runtime/runs-need-root.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
