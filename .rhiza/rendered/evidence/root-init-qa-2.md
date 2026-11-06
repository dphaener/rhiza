# root-init-qa-2

evidence · `01M6S8XAMXYN54E504CRWVWAXN` · actor [claude-code-9c5b9523](../actor/claude-code-9c5b9523.md)

## result

PASS

## statement

QA cycle 2, scratch build of HEAD 0099b258 in a scratch dir (RHIZA_* unset): rhiza root init wrote rhiza.toml (exit 0), a second root init exited 2 with next rhiza root show; root show --json gave the same root from the root, from registered project app and from its app.feat sibling worktree.

## Links

### Out

- verifies [root-init](../capability/rhiza/runtime/root-init.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
