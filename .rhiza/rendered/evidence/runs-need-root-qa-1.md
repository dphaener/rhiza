# runs-need-root-qa-1

evidence · `01M6S53QVFPY4Z5SCN2K76F3DQ` · actor [claude-code-9419651a](../actor/claude-code-9419651a.md)

## result

PASS

## statement

QA cycle 1, hands-on with a scratch build of 9402a003 in a scratch root (RHIZA_ROOT etc. unset): In an unregistered repo under the root, rhiza request and run start --dry-run exited 2 with next: rhiza project add <path>; in a repo under no root, rhiza request exited 2 with next: rhiza root init; rhiza check (after render) passed in both; rhiza request in registered project app queued a feature request.

## Links

### Out

- verifies [runs-need-root](../capability/rhiza/runtime/runs-need-root.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
