# project-add-url-qa-2

evidence · `01M6S8XAZ794Y0MXXNH6GKWPAB` · actor [claude-code-9c5b9523](../actor/claude-code-9c5b9523.md)

## result

PASS

## statement

QA cycle 2, built binary: project add of a file:// URL whose target folder already exists exited 2 with a --name next; a nonexistent file:// URL exited 2 with the git error and left rhiza.toml unchanged. (Cycle 1 QA covered a successful URL clone.)

## Links

### Out

- verifies [project-add-url](../capability/rhiza/runtime/project-add-url.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
