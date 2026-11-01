# default-view-no-fold-verified

evidence · `01M6C5Q7NDB5JPVBAK8GPRTW3J` · actor [claude-code-1bbb0c54](../actor/claude-code-1bbb0c54.md)

## result

PASS

## statement

At 83044f97 ran go test ./internal/api -run 'TestDefaultView|TestReadBudget' -count=1 -v: TestDefaultViewFoldsNothing passes (counting log reader via Ctx.readLog; landed/accepted/draft change sets by ID and by ref pick as-is, to-be, what-if:<id> with 0 folds; current landed change set as-is with 0 folds; ChangeSet+ChangeRef and unknown ref still error); TestDefaultViewSeesPending passes (a dry-run pending change set reaches the view rule as draft, via the fold path); TestReadBudget passes with default-view reads over a landed current change set at medians 2.5-8.3 ms, under 50 ms. go test ./... all ok; fmt-check and vet clean; rhiza check ok in a clean checkout of the commit. Read the diff: no-view reads with no pending entries resolve ChangeSet/ChangeRef/current against index.ChangeSets with the same errors and zero-on-unknown behaviour as effectiveChangeSet; pending writes keep foldedChangeSet; the view rule is unchanged.

## Links

### Out

- verifies [default-view-no-fold](../capability/rhiza/views/fast-reads/default-view-no-fold.md)

## History

- `01M6C3PVXPV5DFQ4G0N6JQGH0Q` A read without --view takes the effective change set's state from the index instead of folding the log
