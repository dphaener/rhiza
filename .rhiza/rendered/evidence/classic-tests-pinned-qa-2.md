# classic-tests-pinned-qa-2

evidence · `01M7DTKQ14Q9ZBWSSMW5FZBMCE` · actor [claude-code-93a5cadf](../actor/claude-code-93a5cadf.md)

## result

PASS

## statement

QA at 1b3323a74, scratch build of the branch: All 30 named scripts pass, in three capped batches; the api tests TestPrePushJournal and TestInboxReadsChangeSetsOnce pass, as do the eight workflow tests. grep -L ship-change-classic on the 33 files prints nothing. Builtin("ship-change") appears only in TestBuiltinShipChange.

## Links

### Out

- verifies [classic-tests-pinned](../capability/rhiza/workflows/built-in-workflows/classic-tests-pinned.md)

## History

- `01M79KH59D0HEME6B85PXXB056` The builtin ship-change becomes the conductor default (conduct, gate, signoff, push, CI, merge, main CI, close-out; every red hands back to conduct) and the old one stays as ship-change-classic
