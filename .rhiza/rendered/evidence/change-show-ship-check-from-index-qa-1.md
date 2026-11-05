# change-show-ship-check-from-index-qa-1

evidence · `01M6NYQH1B6E2AGGJKN8W806P0` · actor [claude-code-edef540b](../actor/claude-code-edef540b.md)

## result

PASS

## statement

QA at 35496535, hands-on with binaries built from this branch (rhiza-new) and from the merge-base (rhiza-old), each run in its own scratch git repo that holds a copy of this repo's real .rhiza (2,617 entries, 51 MB), with RHIZA_ROOT/CHANGE/RUN unset. `rhiza change show 01M6NE0XK2H70R52C5JAPWK71M --json`: with the index warm, new took 0.16-0.34 s CPU per call against 9.1-10.4 s on old. After this branch's own writes it still took 0.19-0.31 s, against 8.2-10.0 s on old. The first call built the index (8.3 s). Normalised change show output of a change set written with 8 writes is byte-identical across both binaries. TestReadBudget on a 4-core VM at load ~3-8: change show draft/accepted fresh, medians 36.4/29.2 ms; after an unrelated entry, 37.6/31.8 ms (budget 50 ms). PASS. TestScripts/(change-show-serves|ship-check-settles-conflict|ship-check-ci-test-findings|ship-check-follows-moved-head): PASS. go test -trimpath -timeout 20m with make's BUDGET_TESTS skip over ./internal/api/... ./internal/index/... ./internal/changeset/...: all ok.

## Links

### Out

- verifies [change-show-ship-check-from-index](../capability/rhiza/views/fast-reads/change-show-ship-check-from-index.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
