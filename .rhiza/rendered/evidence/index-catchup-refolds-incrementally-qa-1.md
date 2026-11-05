# index-catchup-refolds-incrementally-qa-1

evidence · `01M6NYQNMCGF07718E0WCNMC8H` · actor [claude-code-edef540b](../actor/claude-code-edef540b.md)

## result

PASS

## statement

QA at 35496535, hands-on with binaries built from this branch (rhiza-new) and from the merge-base (rhiza-old), each run in its own scratch git repo that holds a copy of this repo's real .rhiza (2,617 entries, 51 MB), with RHIZA_ROOT/CHANGE/RUN unset. Every write on the new binary is followed by a catch-up on the next command. After each add, link, apply or retire, the next `change show` took 0.2-0.3 s CPU, against 7-10 s for a whole-log fold on old. `rhiza check` (a full FoldLog) by both binaries on the index-written repo: ok. go test -trimpath -count=1 ./internal/index/...: ok, including the refold property and counter tests. TestCatchUpBudget: catch-up after one entry, median 34.0 ms (budget 50 ms), PASS. TestBudget: PASS.

## Links

### Out

- verifies [index-catchup-refolds-incrementally](../capability/rhiza/views/fast-reads/index-catchup-refolds-incrementally.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
