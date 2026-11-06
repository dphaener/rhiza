# root-show-verified

evidence · `01M6S0N3X0MG6CZE007G5NRKXX` · actor [claude-code-1057fbaf](../actor/claude-code-1057fbaf.md)

## result

PASS

## statement

At cba75664 ran go test -trimpath -count=1 ./cmd/rhiza -run 'TestScripts/(root-init|describe-golden)$': PASS (root-init ran, --- PASS). Read the diff and the script: it runs rhiza root show from the root itself (JSON project "" and human 'project: none'), from registered project app (project app, human root:/config:/project: lines), from the sibling worktree app.feature-x/sub (project app), from a lookalike folder apple (project ""), and from a directory outside any root (exit=2, stderr 'not inside a Rhiza root' and 'next: rhiza root init', --json code usage with next 'rhiza root init'). Implementation in internal/api/root.go runRootShow and internal/rootcfg Find/ProjectOf matches the statement.

## Links

### Out

- verifies [root-show](../capability/rhiza/runtime/root-show.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
