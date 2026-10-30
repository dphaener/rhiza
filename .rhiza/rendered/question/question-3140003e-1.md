# question-3140003e-1

question · `01M65062MMCP5FWMAFPZENP75F`

## session

3140003e-5643-4487-b77d-01a18a0128af

## statement

Which actors must have their records committed in the build worktree? Decides the check's scope. / How do an actor's record entries get committed? Decides the commit shape.

## Links

### In

- answers [builder-actor-records-committed-with-leaf](../decision/builder-actor-records-committed-with-leaf.md)

## History

- `01M5ZNV3K796MKWEZM3TBQXGZ8` Builder sessions' actor records are committed with the work they sign: after a build run no builder's actor record is left untracked in the build worktree, so every Actor: trailer names a committed actor.
