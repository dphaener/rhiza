# project-add-url-qa-1

evidence · `01M6S53QBJ8VJJW0GN1TN1A8SC` · actor [claude-code-9419651a](../actor/claude-code-9419651a.md)

## result

PASS

## statement

QA cycle 1, hands-on with a scratch build of 9402a003 in a scratch root (RHIZA_ROOT etc. unset): rhiza project add file://<bare>.git cloned into <root>/src and registered it (project ls shows it); a bad file:// URL exited 2 with the git error and a next and left rhiza.toml byte-identical; cloning into an existing folder (--name notgit) exited 2 with a next.

## Links

### Out

- verifies [project-add-url](../capability/rhiza/runtime/project-add-url.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
