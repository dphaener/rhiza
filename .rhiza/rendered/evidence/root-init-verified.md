# root-init-verified

evidence · `01M6S0MZ41VPJ4MSB19NCDPXTM` · actor [claude-code-c860a90e](../actor/claude-code-c860a90e.md)

## result

PASS

## statement

At eac8b7d4, go test -trimpath ./cmd/rhiza -run 'TestScripts/(root-init|describe-golden)$' -count=1 passes: root-init.txtar runs rhiza root init top, then rhiza root show --json from top, from the registered project top/app and from the sibling worktree top/app.feature-x/sub, matching root .*/top each time (project "" then "app" twice), and a second root init exits non-zero with 'already a Rhiza root' and 'next: rhiza root show' (code usage under --json). By hand with a scratch build in a fresh dir: the second init exited 2 with that next, and root show --json printed the identical absolute root and config from top, top/app and top/app.wt/sub. go test ./internal/rootcfg passes (Find, Init ErrExists, ProjectOf rejects apple and app.), and the internal/api registry, describe, schema and humaner tests pass with root.init/root.show registered. Diff read: rootcfg walks up to the nearest rhiza.toml; Init uses O_EXCL so it writes nothing on a second call.

## Links

### Out

- verifies [root-init](../capability/rhiza/runtime/root-init.md)

## History

- `01M6RN3M8EC8FEDEJ8X6YVCPZ4` Root and registry: a root directory's rhiza.toml registers its projects, rhiza root init and rhiza project add/ls/rm adopt or clone them, and every command inside a project resolves its root and project
