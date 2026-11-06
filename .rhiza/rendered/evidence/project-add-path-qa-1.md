# project-add-path-qa-1

evidence · `01M6S53Q8AKS6BRQ7JKWWY7FMA` · actor [claude-code-9419651a](../actor/claude-code-9419651a.md)

## result

PASS

## statement

QA cycle 1, hands-on with a scratch build of 9402a003 in a scratch root (RHIZA_ROOT etc. unset): rhiza project add app registered app with its repo ID (app had .rhiza/); project add bare --name plain registered plain without a repo; rhiza.toml listed both; an outside path, nested/deep, a non-git folder and a duplicate name each exited 2 with a next.

## Links

### Out

- verifies [project-add-path](../capability/rhiza/runtime/project-add-path.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
