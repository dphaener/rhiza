# builder-actor-interview

source · `01M65085GNHRY9JC623VBJEVAF`

## quotes

```json
[]
```

## raw_sha256

fa7d0b1c95dec7aa18d9c4df7ce55c4a2056b598d5200935552522343f3ec9ca

## session

3140003e-5643-4487-b77d-01a18a0128af

## summary

The plan run for change set 01M5ZNV3 raised two questions: which actors must have their records committed in the build worktree, and how their entries get committed. The arbiter answered with the Decision builder-actor-records-committed-with-leaf: each builder's actor record is committed in the builder's own leaf commit, and only builder sessions are in scope. No owner words were available, so no quotes are recorded.

## title

Builder sessions' actor records, planned under autonomous arbitration

## Links

### In

- derived-from [builder-actor-in-leaf-commit](../capability/rhiza/builder-actor-in-leaf-commit.md)

## History

- `01M5ZNV3K796MKWEZM3TBQXGZ8` Builder sessions' actor records are committed with the work they sign: after a build run no builder's actor record is left untracked in the build worktree, so every Actor: trailer names a committed actor.
