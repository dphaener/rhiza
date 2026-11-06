# project-add-path-verified

evidence · `01M6S4B1BNR2CZ9QJGDF4M57BQ` · actor [claude-code-c3d6e71e](../actor/claude-code-c3d6e71e.md)

## result

PASS

## statement

At 3b87a097 (attempt 2): the fix commit drops ReadOnly from project.ls and root.show, and `go test -trimpath -short -count=1 -run TestReadOnlyMarker ./internal/api/` passes, which clears attempt 1's finding. `go test -trimpath ./cmd/rhiza -run 'TestScripts/(project-add-path|describe-golden|root-init)$'` and `go test -trimpath ./internal/rootcfg/...` pass. The project-add-path testscript adds a checkout and finds it in rhiza.toml and in project ls --json, records the repo ID with --name when .rhiza/ exists, and refuses an outside path, a nested path, a non-git folder and a duplicate name. By hand with a scratch build, each of the four refusals exits 2 with a next: line, and rhiza.toml lists [projects.app] path = "app". Separate from this leaf: TestRunStartQueuesForTheDaemon in internal/api fails ('repo is not inside a Rhiza root'). It passes at 85470cb5^ and fails at 85470cb5, so the runs-need-root leaf introduced it, and it will fail CI's make check.

## Links

### Out

- verifies [project-add-path](../capability/rhiza/runtime/project-add-path.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
