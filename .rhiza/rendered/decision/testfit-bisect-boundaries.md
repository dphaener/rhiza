# testfit-bisect-boundaries

decision · `01M7EXVP2KEB9D20D32T0NYMKJ`

## reasoning

The brief says test the state boundaries instead of sweeping every budget. Bisection assumes the state is monotone in the budget, which is what the test asserts, so the boundary neighbours and a seeded sample between keep a check on monotonicity off the boundaries; a fixed seed keeps the run deterministic. The footer subtest sweeps the same way for the first two-cut budget, so it uses the same search. The 'empty sections are never cut' sweep is over tiny and stays as it is.

## statement

TestFit's sweeps over big (cuts in order, for markdown and JSON, and footer) find each cut state's first budget by bisection over the budget and check the budgets just below, at and just above each boundary, plus 64 budgets drawn from a fixed-seed math/rand source between uncut and 1, running every per-budget assertion the sweep ran (state order, every state reached, never-cut sections unchanged, tokens within the budget unless over, fit's byte rules); the monotone-state check is factored into a function a test feeds a hand-made out-of-order sequence

## wording

(assumed, would have escalated)

## Links

### Out

- derived-from [cheap-start-plan-review](../source/cheap-start-plan-review.md)

### In

- decided-by [testfit-boundaries](../capability/rhiza/integrity/self-check/testfit-boundaries.md)

## History

- `01M7EQ42Y1M8ENZCQ202D761JM` Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).
