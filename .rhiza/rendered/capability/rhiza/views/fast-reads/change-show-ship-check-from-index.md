# change-show-ship-check-from-index

capability · `01M6NEM5KT6HR07M515NWRDSRD` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath ./internal/api -run TestReadBudget : TestReadBudget times rhiza change show of an in-flight change set on the 50,000-node budget model under the 50 ms median, both with the index fresh and after one entry is appended to an unrelated change set; go test -trimpath ./internal/api/... : a Go test with a counting log reader and a fresh index asserts change show (including its notes from the change set's history) and ship check --settled read the log zero times, also after an unrelated entry is appended; go test -trimpath ./cmd/rhiza -run 'TestScripts/(change-show-serves|ship-check-settles-conflict|ship-check-ci-test-findings|ship-check-follows-moved-head)$' passes unchanged

## statement

rhiza change show and rhiza ship check --settled take the model and change sets from the worktree's index, caught up incrementally, instead of reading and folding the whole log, with the same answers as before

## Links

### Out

- depends-on [index-catchup-refolds-incrementally](index-catchup-refolds-incrementally.md)
- decided-by [model-reads-go-through-the-index](../../../../decision/model-reads-go-through-the-index.md)

### In

- verifies [change-show-ship-check-from-index-verified](../../../../evidence/change-show-ship-check-from-index-verified.md)
- verifies [change-show-ship-check-from-index-qa-1](../../../../evidence/change-show-ship-check-from-index-qa-1.md)
- verifies [change-show-ship-check-from-index-verified-2](../../../../evidence/change-show-ship-check-from-index-verified-2.md)
- verifies [pr-304-fix-run-drivers-burn-cpu-while](../../../../evidence/pr-304-fix-run-drivers-burn-cpu-while.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
