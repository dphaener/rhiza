# inbox-stale-rebase-resolution

decision · `01M6358252CFGR7B85YBE8WRQE`

## reasoning

The two change sets each added a change-set lookup to inboxItems: main a state lookup, this branch a gone predicate. gone is the state lookup plus Terminal(), so one parameter serves both and no call site changes on main.

## statement

The rebase onto main is resolved by keeping main's changeState (from resume-only-in-flight runs) as inboxItems' parameter and deriving gone from it inside inboxItems (a known change set in a terminal state). goneChanges is dropped. Both changes keep their behaviour, and make check passes on the rebased branch.

## wording

decided by the orchestrator acting as arbiter, under the owner's standing word to run autonomously

## Links

### Out

- answers [rebase-conflict-inbox-clears-stale-failures](../question/rebase-conflict-inbox-clears-stale-failures.md)

## History

- `01M62MKZFRPZ2C7QTYAF815XAG` The inbox leaves out failed runs whose change set has since landed, been withdrawn or rejected, and rhiza run abandon also takes a failed run, so every failed item can be cleared; no item's next command is one that can only fail again.
