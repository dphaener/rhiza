# root-init-qa-1

evidence · `01M6S53PVJ4HF32THHA3M6YJ0Y` · actor [claude-code-9419651a](../actor/claude-code-9419651a.md)

## result

PASS

## statement

QA cycle 1, hands-on with a scratch build of 9402a003 in a scratch root (RHIZA_ROOT etc. unset): rhiza root init root wrote root/rhiza.toml (exit 0); a second root init on it exited 2 with next: rhiza root show; root show --json gave the same root from the root, from a registered project's subfolder app/sub and from a git worktree sibling app.feature-x.

## Links

### Out

- verifies [root-init](../capability/rhiza/runtime/root-init.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
