# apply-validates-from-index

capability · `01M6NEMD8NQT1DE0SEYVGE7X2X` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath ./internal/api/... : a Go test with a counting log reader and a fresh index asserts an apply of several writes reads the log zero times, and a property test over generated logs and writes finds every write accepted or refused, with the same problems, through the index as through the full fold; go test -trimpath ./internal/api -run TestReadBudget : a budget case times one rhiza apply write on the 50,000-node budget model under the 200 ms median; go test -trimpath ./cmd/rhiza -run 'TestScripts/(apply|link)$' passes unchanged

## statement

rhiza apply (and the single-write commands add, set, link, retire, mv) validate a write against the index's folded state plus the write's own new entries, so a write never reads and folds the whole log

## Links

### Out

- depends-on [index-catchup-refolds-incrementally](index-catchup-refolds-incrementally.md)
- decided-by [model-reads-go-through-the-index](../../../../decision/model-reads-go-through-the-index.md)

### In

- verifies [apply-validates-from-index-verified](../../../../evidence/apply-validates-from-index-verified.md)
- verifies [apply-validates-from-index-qa-1](../../../../evidence/apply-validates-from-index-qa-1.md)
- verifies [pr-304-fix-run-drivers-burn-cpu-while](../../../../evidence/pr-304-fix-run-drivers-burn-cpu-while.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
