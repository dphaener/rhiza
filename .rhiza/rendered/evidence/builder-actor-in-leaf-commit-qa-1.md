# builder-actor-in-leaf-commit-qa-1

evidence · `01M6533NTRNRHRYCF9SCEY7T0P` · actor [claude-code-2c09082f](../actor/claude-code-2c09082f.md)

## result

PASS

## statement

QA cycle 1 with a binary built from 7cd456e. Real run: the build run's leaf commits 3d9beb0 and 8f19664 each carry the record of the actor their Actor: trailer names (claude-code-32a2e794, claude-code-27db46d5). The build worktree's only untracked .rhiza/log entries are review sessions' actor records, which the Decision leaves out of scope. Hands-on: on a scratch clone at 6127282 holding five real untracked actor-record change sets copied from the build worktree, with RHIZA_ACTOR=claude-code-d2a10324, I ran the `git add` line taken from `rhiza guide build` of the built binary through eval. It staged only that actor's three-entry record change set. After the commit with Leaf:/Actor: trailers the other four stayed untracked. On a clean clone of that commit, `rhiza show claude-code-d2a10324 --json` reads kind actor, and `rhiza check` folds with no problems; it fails only on the expected stale render. build-evidence-commit and internal/guide tests pass. `make check` exits 0 (an earlier run failed once on a read-budget timing flake, and a rerun of go test ./... and make check passed).

## Links

### Out

- verifies [builder-actor-in-leaf-commit](../capability/rhiza/builder-actor-in-leaf-commit.md)

## History

- `01M5ZNV3K796MKWEZM3TBQXGZ8` Builder sessions' actor records are committed with the work they sign: after a build run no builder's actor record is left untracked in the build worktree, so every Actor: trailer names a committed actor.
