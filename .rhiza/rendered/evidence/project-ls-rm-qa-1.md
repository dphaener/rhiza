# project-ls-rm-qa-1

evidence · `01M6S53QHFNHK3H3GYVGRM8CPB` · actor [claude-code-9419651a](../actor/claude-code-9419651a.md)

## result

PASS

## statement

QA cycle 1, hands-on with a scratch build of 9402a003 in a scratch root (RHIZA_ROOT etc. unset): rhiza project ls (text and --json) listed app initialised with its repo ID and plain not initialised; project rm src2 unregistered it and left src2/ on disk; project rm ghost exited 2 with next: rhiza project ls.

## Links

### Out

- verifies [project-ls-rm](../capability/rhiza/runtime/project-ls-rm.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
