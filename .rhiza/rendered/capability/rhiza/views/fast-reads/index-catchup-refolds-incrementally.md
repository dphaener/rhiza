# index-catchup-refolds-incrementally

capability · `01M6NEZN44VV5TKPM24EY59X4R` · parent [fast-reads](../fast-reads.md)

## acceptance

go test -trimpath ./internal/index/... : a property test over generated logs appends entries one at a time and asserts after each catch-up that the stored change sets (with history) and fold problems equal changeset.FoldLog over the whole log; a Go test with decode and fold counters appends one entry to an unrelated change set of a 100+ change set log and asserts catch-up decodes only that change set's entries and folds only that change set; a budget test times catch-up after one appended entry on the 50,000-node budget model under 50 ms median

## statement

The index's catch-up decodes and folds only the entries it newly reads: the change sets they touch are refolded from their own stored entries (index.EntriesOf) and the rest keep their stored folded state, so appending one entry costs work in proportion to that entry's change set, not a whole-log decode (entry.Decode with canon.Check) and changeset.FoldLog; the index stores each change set's folded history and entries, so change show's notes and apply's validation can take whole change sets from it; stored change sets and fold problems still equal FoldLog over the whole log

## Links

### Out

- decided-by [model-reads-go-through-the-index](../../../../decision/model-reads-go-through-the-index.md)

### In

- depends-on [run-driver-idle-poll-cheap](../../orchestration/daemon/run-driver-idle-poll-cheap.md)
- depends-on [change-show-ship-check-from-index](change-show-ship-check-from-index.md)
- depends-on [apply-validates-from-index](apply-validates-from-index.md)
- verifies [index-catchup-refolds-incrementally-verified](../../../../evidence/index-catchup-refolds-incrementally-verified.md)
- verifies [index-catchup-refolds-incrementally-qa-1](../../../../evidence/index-catchup-refolds-incrementally-qa-1.md)
- verifies [pr-304-fix-run-drivers-burn-cpu-while](../../../../evidence/pr-304-fix-run-drivers-burn-cpu-while.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
