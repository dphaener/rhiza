# testfit-boundaries-verified

evidence · `01M7FMFF3XATHGH8PQ6TYQR599` · actor [claude-code-3456e2fb](../actor/claude-code-3456e2fb.md)

## result

PASS

## statement

Verified at c9e38ec85. I re-ran the Checks: block and the acceptance under the machine lock. /usr/bin/time go test -trimpath -count=1 -p 2 ./internal/slice -run '^TestFit$' exited 0 in 1.48 s wall (the limit is 5 s). go test -trimpath -count=1 -p 2 ./internal/slice exited 0 (16.6 s). TestCheckOrder passes. Its 'skips a state' and 'goes back a state' cases each check the exact error text, which names the budget and both states. Coverage: new.out came from HEAD and old.out from a detached worktree of 88516d956 under $TMPDIR, which I removed afterwards. Of 48 internal/slice/budget.go blocks, 39 are covered in old.out and none of them is missing from new.out. budget.go is unchanged between 88516d956 and HEAD. The only code file in the diff is internal/slice/budget_test.go; the rest is the builder's .rhiza actor record. No assertion was loosened: the state closure keeps the over/tokens-within-budget and never-cut-sections checks on every probed budget, including the bisection probes. Order and every-state-reached still hold, and the order check now also catches a skipped state between neighbouring budgets. The footer subtest finds the first briefing budget by bisection, keeps its three assertions and adds a not-over, exactly-two-cut check. TestFitWithin's comment no longer says TestFit sweeps every budget. go vet passes, gofmt is clean, and there are no TODOs.

## Links

### Out

- verifies [testfit-boundaries](../capability/rhiza/integrity/self-check/testfit-boundaries.md)

## History

- `01M7EQ42Y1M8ENZCQ202D761JM` Cheap process start and boundary TestFit: a rhiza process no longer pays go-runewidth's eager init at start (22 ms of 39 ms, across 6,600+ starts in scripts), and internal/slice TestFit tests the fit's state boundaries instead of sweeping every token budget (42 s).
