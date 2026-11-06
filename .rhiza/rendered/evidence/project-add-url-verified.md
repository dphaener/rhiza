# project-add-url-verified

evidence · `01M6S2ZAZM8WJQSFVE9GGZEQCT` · actor [claude-code-078ababe](../actor/claude-code-078ababe.md)

## result

PASS

## statement

Verifier claude-code-078ababe at c4efeed7: go test -trimpath ./cmd/rhiza -run 'TestScripts/(project-add-url|project-add-path|describe-golden)$' -count=1 passes. project-add-url.txtar creates a local bare repo, runs rhiza project add file://$WORK/remote.git, sees 'Registered project remote at remote.', remote/.git under the root, projects.remote in rhiza.toml and name remote in project ls --json; --name again clones to again/; a taken name is refused; file://$WORK/missing.git exits non-zero with a next: hint, leaves no missing/ folder, and cmp shows rhiza.toml unchanged. Re-ran by hand with a scratch build (no network): clone + register worked, the bad URL exited 2 with git's error, rhiza.toml still listed only r, and project ls showed it. Read the diff: isGitURL/urlName/runProjectClone in internal/api/project.go check name and destination before cloning and remove the partial folder on any failure after clone. go vet and gofmt clean. Unrelated to this leaf: go test -short ./internal/api fails TestReadOnlyMarker because project.ls (6571762e) and root.show (eac8b7d4) are marked ReadOnly but not in the test's web-UI list. This commit touches no ReadOnly mark.

## Links

### Out

- verifies [project-add-url](../capability/rhiza/runtime/project-add-url.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
