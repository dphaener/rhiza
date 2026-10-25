# ordered-landing-tests

evidence · `01M5T2DGXDS88SEX05WYNJSGR3`

## result

make check exit 0 at commit ea4c7f7

## statement

Go tests and testscripts show land refusing a change set until every dependency has landed, and refusing one with a conflict, with next.

## tests

go test: TestLandDependencies, TestLandConflict, TestLandRevisedLanded, TestLandProblem (internal/api); testscript: stacked.txtar, stacked-revised.txtar, stacked-revised-disjoint.txtar, stacked-revised-rewrite.txtar, stacked-revised-drop.txtar

## Links

### Out

- verifies [ordered-landing](../capability/genesis/change-sets/ordered-landing.md)

## History

- `01M5SWTT2A55K0REHQ6JNJQZ6W` Writes see the accepted, unlanded change sets they depend on, so a sequence of change sets can build on each other. This is the gap behind the built/intended split and milestone 2's most important item.
