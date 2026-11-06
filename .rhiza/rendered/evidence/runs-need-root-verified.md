# runs-need-root-verified

evidence · `01M6S38D4BM1T7QH79EN8SVDP7` · actor [claude-code-2dd23c4e](../actor/claude-code-2dd23c4e.md)

## result

PASS

## statement

At 85470cb5, ran go test -trimpath ./cmd/rhiza -run 'TestScripts/^(...)$' over the 120 testscripts that call rhiza run start or rhiza request, plus root-* and project-* (including the new runs-need-root.txtar): ok in 435.8s. The new script runs rhiza request and rhiza run start --workflow plan in an unregistered repo under the Setup root (next: rhiza project add <dir>) and in $NOROOT (next: rhiza root init); both fail, and rhiza render and rhiza check then succeed in each. Checked the exit code by hand with a scratch build: in a root with an unregistered repo, and in a folder under no root, rhiza request, rhiza request --json and rhiza run start each exit 2 with code usage and the matching next. In both places rhiza check passes after rhiza render. Read the diff: requireProject in internal/api/root.go guards only runRequest and runRunStart, so model commands are untouched. Setup writes $WORK/rhiza.toml registering repo and sets $NOROOT. The changed scripts register their extra project folders.

## Links

### Out

- verifies [runs-need-root](../capability/rhiza/runtime/runs-need-root.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
