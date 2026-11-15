# testfit-boundaries-qa-1

evidence · `01M7FN89CQFJQWBMZNCGKRQAGD` · actor [claude-code-9160dfa9](../actor/claude-code-9160dfa9.md)

## result

PASS

## statement

QA at c9e38ec85: /usr/bin/time go test -trimpath -count=1 ./internal/slice -run '^TestFit$' ok in 2.67 s wall (was 42 s); go test -trimpath -count=1 -p 2 ./internal/slice ok (19.0 s); TestCheckOrder passes, and its skip and go-back cases assert exact errors naming the budget and both states. Coverage: -coverprofile for TestFit here and in a detached worktree of 88516d956: 39 budget.go blocks covered in old.out, 0 missing in new.out. git diff from the merge-base leaves budget.go untouched; the brief-two-sections subtest now also asserts not over and exactly two cut, and TestFitWithin's comment no longer says TestFit sweeps every budget.

## Links

### Out

- verifies [testfit-boundaries](../capability/rhiza/integrity/self-check/testfit-boundaries.md)

## History

- `01M7EQ42Y1M8ENZCQ202D761JM` Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).
