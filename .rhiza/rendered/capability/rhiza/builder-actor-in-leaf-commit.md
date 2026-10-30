# builder-actor-in-leaf-commit

capability · `01M6508E9HB879QTYQ1JJFJ8WF` · parent [rhiza](../rhiza.md)

## acceptance

A build run script leaves git status clean of builder actor entries under .rhiza/log, and the builder's leaf commit carries that actor's record so fold-ok reads the actor as committed

## statement

Each builder session's actor record is committed in its own leaf commit, so after a build run every builder's Actor: trailer names an actor whose record is tracked

## Links

### Out

- decided-by [builder-actor-records-committed-with-leaf](../../decision/builder-actor-records-committed-with-leaf.md)
- derived-from [builder-actor-interview](../../source/builder-actor-interview.md)

### In

- verifies [builder-actor-in-leaf-commit-verified](../../evidence/builder-actor-in-leaf-commit-verified.md)
- verifies [builder-actor-in-leaf-commit-qa-1](../../evidence/builder-actor-in-leaf-commit-qa-1.md)
- verifies [pr-98-builder-sessions-actor](../../evidence/pr-98-builder-sessions-actor.md)

## History

- `01M5ZNV3K796MKWEZM3TBQXGZ8` Builder sessions' actor records are committed with the work they sign: after a build run no builder's actor record is left untracked in the build worktree, so every Actor: trailer names a committed actor.
