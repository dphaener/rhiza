# root-show-qa-1

evidence · `01M6S53Q3T3G1182ZWSSQYQ4ZK` · actor [claude-code-9419651a](../actor/claude-code-9419651a.md)

## result

PASS

## statement

QA cycle 1, hands-on with a scratch build of 9402a003 in a scratch root (RHIZA_ROOT etc. unset): rhiza root show --json gave project "" at the root itself, project app in app/sub and in the app.feature-x sibling worktree, project "" in an unregistered folder under the root, and outside any root exited 2 with next: rhiza root init.

## Links

### Out

- verifies [root-show](../capability/rhiza/runtime/root-show.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
