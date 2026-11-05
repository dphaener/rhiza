# apply-validates-from-index-qa-1

evidence · `01M6NYQCBR3NZ5JRYXEYPD4EHK` · actor [claude-code-edef540b](../actor/claude-code-edef540b.md)

## result

PASS

## statement

QA at 35496535, hands-on with binaries built from this branch (rhiza-new) and from the merge-base (rhiza-old), each run in its own scratch git repo that holds a copy of this repo's real .rhiza (2,617 entries, 51 MB), with RHIZA_ROOT/CHANGE/RUN unset. Timed with /usr/bin/time. A single `rhiza add capability` into a new draft took 0.27-0.34 s CPU (new, index warm) against 9.8-11 s (old). `link ... verifies` took 0.34 s against 9.4 s. `retire --reason` took 0.14 s against 10.2 s. `rhiza apply` of 4 writes (2 adds, mv under parent, set) took 0.34 s against 24.6 s. Refusals matched on both binaries: duplicate handle gave conflict, a missing link target gave not_found, set on a retired node gave not_found, and mv under its own descendant gave the usage error. Normalised `change show --json` of the resulting change set (IDs and times stripped) is byte-identical across both binaries. `rhiza check` by both binaries on both repos: ok, no problems. go test -trimpath -count=1 ./internal/api -run 'TestWritesThroughIndexMatchFold|TestApplyIndexedReadsNoLog': PASS, 18-29 accepted and 51-62 refused per seed. TestReadBudgetApply: median 124.9 ms (budget 200 ms), zero log reads, PASS. TestScripts/(apply|link): PASS.

## Links

### Out

- verifies [apply-validates-from-index](../capability/rhiza/views/fast-reads/apply-validates-from-index.md)

## History

- `01M6NE0XK2H70R52C5JAPWK71M` Fix Run drivers burn CPU while they wait: each 'rhiza run start --foreground' driver uses ~14% of a core continuously (2.5-3 CPU-min over 19 min of wall time, measured 2026-11-05 on drivers 3213906/3213914/3213930 while their runs sat in build/review), so 5-7 live runs spend about a core on supervision alone on the 4-core VM and load reached 41. rhiza apply also takes ~3 s per write in this repo. Profile a driver's poll loop and rhiza apply: likely re-reading and folding the whole model log (now large) on every poll or write; poll cheaply (file mtimes or the run journal), cache the fold, and add a CPU budget test for an idle driver.
