# project-ls-rm-verified

evidence · `01M6S22H2JHBCXY9QC8XKESAHB` · actor [claude-code-9c43fd95](../actor/claude-code-9c43fd95.md)

## result

PASS

## statement

At commit 1b7da491 I ran go test -trimpath ./cmd/rhiza -run 'TestScripts/(project-ls-rm|describe-golden|project-)' (ok), go test ./internal/rootcfg/... (ok), go test ./internal/api -run 'Project|Root' (ok), go vet and gofmt (clean). The project-ls-rm testscript registers a bare git checkout and an initialised one, sees 'not initialised' and 'initialised' in project ls and initialised false/true in project ls --json, runs project rm bare, finds bare/.git still on disk and bare gone from rhiza.toml and ls --json, and gets 'no project named bare' with a next: line for a second rm. I also ran a scratch build by hand in a hermetic root: ls --json lists name, path, initialised and repo (repo ID present for the initialised project); rm --dry-run left rhiza.toml unchanged; rm removed only [projects.bare] and kept the folder; a second rm exited 2 with next: rhiza project ls. I read the diff: rootcfg.Remove deletes only the map entry and rewrites the TOML file, and never touches project files.

## Links

### Out

- verifies [project-ls-rm](../capability/rhiza/runtime/project-ls-rm.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
