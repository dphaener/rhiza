# project-ls-rm-qa-2

evidence · `01M6S8XB23WVRKGMW09AJM6CHV` · actor [claude-code-9c5b9523](../actor/claude-code-9c5b9523.md)

## result

PASS

## statement

QA cycle 2, built binary: project ls showed app initialised and plain not initialised (--json with repo on app only); project rm plain exited 0, the bare folder stayed on disk and ls listed only app.

## Links

### Out

- verifies [project-ls-rm](../capability/rhiza/runtime/project-ls-rm.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
