# builder-actor-records-committed-with-leaf

decision · `01M6507263M8AWCJA5KRQEZYTM`

## reasoning

The change set intent says builder sessions' actor records are committed with the work they sign, and that no builder's record is left untracked after a build run. That fixes the scope to builder actors' records. The rhiza-build guide has a builder make one commit with a Leaf: trailer, so a separate record commit would break that shape. Including the record in the builder's one commit keeps the Actor: trailer and its committed record together. Nothing in the slice (no decisions, constraints or evidence) contradicts this.

## statement

Every builder actor record in the build worktree is committed in the builder's own leaf commit, so the Actor: trailer names an actor whose record is tracked after the run. Only builder sessions' actor records are in scope.

## wording

(assumed, would have escalated)

## Links

### Out

- answers [question-3140003e-1](../question/question-3140003e-1.md)

### In

- decided-by [builder-actor-in-leaf-commit](../capability/rhiza/builder-actor-in-leaf-commit.md)

## History

- `01M5ZNV3K796MKWEZM3TBQXGZ8` Builder sessions' actor records are committed with the work they sign: after a build run no builder's actor record is left untracked in the build worktree, so every Actor: trailer names a committed actor.
