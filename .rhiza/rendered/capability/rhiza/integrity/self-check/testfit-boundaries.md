# testfit-boundaries

capability · `01M7EQF5ZRFYQWZJCW8CKY1S8V` · parent [self-check](../self-check.md)

## acceptance

/usr/bin/time -f '%e s' go test -trimpath -count=1 ./internal/slice -run '^TestFit$' passes within 5 s (2026-11-15: 42 s) and go test -trimpath -count=1 ./internal/slice passes; coverage kept: go test -trimpath -count=1 -coverprofile $TMPDIR/new.out ./internal/slice -run '^TestFit$', and the same into $TMPDIR/old.out in a detached worktree of 88516d956: every internal/slice/budget.go block with a count above 0 in old.out has a count above 0 in new.out; failure case: a test feeds the factored order checker a hand-made sequence that skips a state and one that goes back a state, and each fails with the budget and both states named; TestFitWithin's doc comment no longer says TestFit sweeps every budget. Touches internal/slice/budget_test.go only; must not change budget.go or loosen any assertion TestFit makes at a budget it fits.

## statement

internal/slice's TestFit checks big's cuts at the fit's state boundaries instead of fitting at every token budget: the 'cuts in order' subtests (markdown and JSON) and the 'footer' subtest find each cut state's first budget by bisection, and the cut order is checked just below, at and just above each boundary plus at a fixed seeded sample of budgets between, with every per-budget assertion the sweep made.

## Links

### Out

- decided-by [fast-suite-proposal-approved](../../../../decision/fast-suite-proposal-approved.md)
- decided-by [testfit-bisect-boundaries](../../../../decision/testfit-bisect-boundaries.md)

### In

- verifies [testfit-boundaries-verified](../../../../evidence/testfit-boundaries-verified.md)
- verifies [testfit-boundaries-qa-1](../../../../evidence/testfit-boundaries-qa-1.md)

## History

- `01M7EQ42Y1M8ENZCQ202D761JM` Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).
