# default-view-no-fold

capability · `01M6C3QBZ6N7WNS21SYQPCYWFM` · parent [fast-reads](../fast-reads.md)

## acceptance

A Go test with a counting log reader: a default-view read with a landed, an accepted and a draft current change set picks as-is, to-be and what-if as today and folds no change sets; with a pending write it still sees the pending state; TestReadBudget times a default-view read with a landed current change set under the 50 ms budget

## statement

A read without --view takes the effective change set's state from the index, so with a current change set it costs no fold of the log and picks the same view as before

## Links

### Out

- decided-by [default-view-state-from-index](../../../../decision/default-view-state-from-index.md)

### In

- verifies [default-view-no-fold-verified](../../../../evidence/default-view-no-fold-verified.md)
- verifies [default-view-no-fold-qa-1](../../../../evidence/default-view-no-fold-qa-1.md)
- verifies [pr-178-default-reads-skip-the-fold](../../../../evidence/pr-178-default-reads-skip-the-fold.md)

## History

- `01M6C3PVXPV5DFQ4G0N6JQGH0Q` A read without --view takes the effective change set's state from the index instead of folding the log
