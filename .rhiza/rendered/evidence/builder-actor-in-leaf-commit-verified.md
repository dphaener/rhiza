# builder-actor-in-leaf-commit-verified

evidence · `01M652K7HH9XEZGRJWRFJWSC77` · actor [claude-code-b8ac2d4d](../actor/claude-code-b8ac2d4d.md)

## result

PASS

## statement

At 8f19664, cmd/rhiza/testdata/script/build-evidence-commit.txtar runs a build run whose fake builder stages its actor record with the `rhiza guide build` recipe (added to internal/guide/text/build.md in 3d9beb0) and commits it with its leaf file under Leaf:/Actor: trailers. The script checks that the leaf commit's tree holds an actor record whose handle is the Actor: trailer and RHIZA_ACTOR, that `git status --porcelain --untracked-files=all -- .rhiza` is empty after the run, and that fold-ok passes on a clean clone checked out at the leaf commit, where `rhiza show <actor>` reads kind actor. `go test ./cmd/rhiza -run 'TestScripts/build-' -count=1` and `go test ./internal/guide` pass. Mutation check: in a scratch copy with the eval of the staging recipe removed, build-evidence-commit fails at line 64 (no .rhiza/log entry in the leaf commit), so the script catches a builder that leaves its actor untracked. `make check` exits 0 on a clean clone of 8f19664. The builder's own actor claude-code-27db46d5 is committed in 8f19664.

## Links

### Out

- verifies [builder-actor-in-leaf-commit](../capability/rhiza/builder-actor-in-leaf-commit.md)

## History

- `01M5ZNV3K796MKWEZM3TBQXGZ8` Builder sessions' actor records are committed with the work they sign: after a build run no builder's actor record is left untracked in the build worktree, so every Actor: trailer names a committed actor.
